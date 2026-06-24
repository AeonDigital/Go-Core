package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// ErrorInput stores the sanitized metadata required to append a new error token.
type ErrorInput struct {
	FamilyID    int
	FamilyName  string
	IsNewFamily bool
	Constant    string
	Code        string
	Fields      []string
	Message     string
	Technical   string
}

// HandleAddErrorCode manages the entrypoint execution flow for the 'error-code-add' subcommand.
func HandleAddErrorCode(args []string) {
	fs := flag.NewFlagSet("error-code-add", flag.ExitOnError)

	// Inline configuration flags
	flagFamily := fs.Int("family", 0, "Family ID (e.g., 1, 2). Use 0 to interactively select or create one.")
	flagFamilyTitle := fs.String("family-title", "", "Title for the new family (required only if creating a new family inline)")
	flagName := fs.String("name", "", "Error naming key (e.g., fieldRequired, invalid_token)")
	flagFields := fs.String("fields", "", "Comma-separated extra tag fields (e.g., field, value)")
	flagMsg := fs.String("message", "", "Default fallback human-readable error message")
	flagTech := fs.String("tech", "", "Technical description/documentation detailing the error scenario")

	_ = fs.Parse(args)

	targetFile := "01_xerrors.go"
	if _, err := os.Stat(targetFile); os.IsNotExist(err) {
		fmt.Printf("[ERR] Target file '%s' not found. Ensure you are running this command from the package root.\n", targetFile)
		os.Exit(1)
	}

	// 1. Read existing file and parse structural metadata
	contentBytes, err := os.ReadFile(targetFile)
	if err != nil {
		fmt.Printf("[ERR] Failed to read '%s': %v\n", targetFile, err)
		os.Exit(1)
	}
	fileContent := string(contentBytes)

	families, maxCodePerFamily, maxFamilyID := parseXerrorsFile(fileContent)

	// 2. Resolve parameters (Inline vs Interactive Assistant Flow)
	input := &ErrorInput{}
	reader := bufio.NewReader(os.Stdin)

	// Resolve Family Choice
	if *flagFamily > 0 {
		input.FamilyID = *flagFamily
		if input.FamilyID > maxFamilyID+1 {
			fmt.Printf("[ERR] Invalid family ID %d. Next sequential available family is %d.\n", input.FamilyID, maxFamilyID+1)
			os.Exit(1)
		}
		if input.FamilyID == maxFamilyID+1 {
			input.IsNewFamily = true
			input.FamilyName = strings.TrimSpace(*flagFamilyTitle)
			if input.FamilyName == "" {
				fmt.Println("[ERR] Flag '--family-title' is required when declaring a new family identifier inline.")
				os.Exit(1)
			}
		} else {
			input.FamilyName = families[input.FamilyID]
		}
	} else {
		// Run Interactive Prompt for Family Selection
		runFamilyWizard(reader, families, maxFamilyID, input)
	}

	// Resolve Error Name/Constant
	rawName := *flagName
	if rawName == "" {
		fmt.Print("Enter Error Name (camelCase or snake_case): ")
		rawName, _ = reader.ReadString('\n')
	}
	input.Constant = toSnakeUpperCase(rawName)
	if !validateToken(input.Constant) {
		fmt.Printf("[ERR] Error name '%s' normalized to '%s' violates strict constraints (A-Z base, max 32 chars).\n", rawName, input.Constant)
		os.Exit(1)
	}
	// Check for duplicates
	if strings.Contains(fileContent, " "+input.Constant+" ") || strings.Contains(fileContent, "\t"+input.Constant+" ") {
		fmt.Printf("[ERR] Token constant '%s' already exists inside registry. Aborting to protect consistency.\n", input.Constant)
		os.Exit(1)
	}

	// Resolve Next Error Code sequentially
	if input.IsNewFamily {
		input.Code = fmt.Sprintf("E%d001", input.FamilyID)
	} else {
		lastCode := maxCodePerFamily[input.FamilyID]
		if lastCode == 0 {
			input.Code = fmt.Sprintf("E%d001", input.FamilyID)
		} else {
			input.Code = fmt.Sprintf("E%d", lastCode+1)
		}
	}

	// Resolve Extra Fields/Tags
	rawFields := *flagFields
	if rawFields == "" && *flagName == "" { // Only prompt if running interactively
		fmt.Print("Enter Extra Tag Fields (comma-separated, e.g., field, value - optional): ")
		rawFields, _ = reader.ReadString('\n')
	}
	if strings.TrimSpace(rawFields) != "" {
		parts := strings.Split(rawFields, ",")
		for _, p := range parts {
			sanitized := toSnakeUpperCase(p)
			if sanitized != "" {
				if !validateToken(sanitized) {
					fmt.Printf("[ERR] Extra field tag '%s' matches an invalid layout format structure.\n", p)
					os.Exit(1)
				}
				input.Fields = append(input.Fields, sanitized)
			}
		}
	}

	// Resolve Human Message
	input.Message = strings.TrimSpace(*flagMsg)
	if input.Message == "" {
		fmt.Print("Enter Default Fallback Human Message (Required): ")
		input.Message, _ = reader.ReadString('\n')
		input.Message = strings.TrimSpace(input.Message)
	}
	if input.Message == "" {
		fmt.Println("[ERR] Default human message cannot resolve to an empty textual state.")
		os.Exit(1)
	}

	// Resolve Technical Documentation
	input.Technical = strings.TrimSpace(*flagTech)
	if input.Technical == "" {
		fmt.Print("Enter Technical/Architecture Guidance Commentary (Required): ")
		input.Technical, _ = reader.ReadString('\n')
		input.Technical = strings.TrimSpace(input.Technical)
	}
	if input.Technical == "" {
		fmt.Println("[ERR] Technical documentation reference statement is mandatory.")
		os.Exit(1)
	}

	// 3. Inject new generated records back into the physical file
	updatedContent := injectNewErrorIntoString(fileContent, input)

	err = os.WriteFile(targetFile, []byte(updatedContent), 0644)
	if err != nil {
		fmt.Printf("[ERR] Failed writing modifications down to '%s': %v\n", targetFile, err)
		os.Exit(1)
	}

	fmt.Printf("[OKK] Successfully appended '%s' (%s) inside family %d.\n", input.Constant, input.Code, input.FamilyID)
}

// parseXerrorsFile extracts present families and active numeric bounds from the raw file string content.
func parseXerrorsFile(content string) (families map[int]string, maxCode map[int]int, maxFamilyID int) {
	families = make(map[int]string)
	maxCode = make(map[int]int)

	// Regex for tracking constants section family headers
	familyRegex := regexp.MustCompile(`===\s*FAMILY:\s*(\d+)\s*\|\s*TITLE:\s*(.+)`)
	// Regex for catching individual error allocations (e.g., E1001)
	codeRegex := regexp.MustCompile(`"E(\d+)"`)

	scanner := bufio.NewScanner(strings.NewReader(content))
	currentFamilyCtx := 0

	for scanner.Scan() {
		line := scanner.Text()

		if matches := familyRegex.FindStringSubmatch(line); len(matches) > 0 {
			id, _ := strconv.Atoi(matches[1])
			title := strings.TrimSpace(matches[2])
			families[id] = title
			if id > maxFamilyID {
				maxFamilyID = id
			}
			currentFamilyCtx = id
		}

		if currentFamilyCtx > 0 {
			if matches := codeRegex.FindStringSubmatch(line); len(matches) > 0 {
				fullNum, _ := strconv.Atoi(matches[1]) // captures e.g., 1001
				if fullNum > maxCode[currentFamilyCtx] {
					maxCode[currentFamilyCtx] = fullNum
				}
			}
		}
	}
	return
}

func runFamilyWizard(reader *bufio.Reader, families map[int]string, maxFamilyID int, input *ErrorInput) {
	fmt.Println("\nAvailable Structural Error Families:")
	for id, title := range families {
		fmt.Printf("  [%d] %s\n", id, title)
	}
	nextID := maxFamilyID + 1
	fmt.Printf("  [%d] Create a new sequential semantic family\n", nextID)

	fmt.Printf("Select Family Destination Index [%d-%d]: ", 1, nextID)
	choiceStr, _ := reader.ReadString('\n')
	choiceStr = strings.TrimSpace(choiceStr)
	choice, err := strconv.Atoi(choiceStr)

	if err != nil || choice < 1 || choice > nextID {
		fmt.Println("[ERR] Invalid structural index selection.")
		os.Exit(1)
	}

	input.FamilyID = choice
	if choice == nextID {
		input.IsNewFamily = true
		fmt.Print("Enter Corporate Title for the new Error Family: ")
		title, _ := reader.ReadString('\n')
		input.FamilyName = strings.ToUpper(strings.TrimSpace(title))
		if input.FamilyName == "" {
			fmt.Println("[ERR] Family title declaration cannot resolve to empty value states.")
			os.Exit(1)
		}
	} else {
		input.FamilyName = families[choice]
	}
}

func injectNewErrorIntoString(content string, input *ErrorInput) string {
	// 1. Generate Injection Content blocks
	var constBlock strings.Builder
	var regBlock strings.Builder

	expectTokens := "CTX, MSG"
	for _, f := range input.Fields {
		expectTokens += ", " + f
	}
	expectTokens += ", [error]"

	if input.IsNewFamily {
		constBlock.WriteString("\n\t// ============================================================================\n")
		constBlock.WriteString(fmt.Sprintf("\t// === FAMILY: %d | TITLE: %s\n", input.FamilyID, input.FamilyName))
		constBlock.WriteString("\t// ============================================================================\n")

		regBlock.WriteString("\n\t// ============================================================================\n")
		regBlock.WriteString(fmt.Sprintf("\t// === FAMILY: %d\n", input.FamilyID))
		regBlock.WriteString("\t// ============================================================================\n")
	}

	constBlock.WriteString(fmt.Sprintf("\n  // %s belongs to Family %d.\n", input.Constant, input.FamilyID))
	constBlock.WriteString(fmt.Sprintf("  // %s\n", input.Technical))
	constBlock.WriteString(fmt.Sprintf("  // Format expects: %s\n", expectTokens))
	constBlock.WriteString(fmt.Sprintf("  %s xerrors.ErrorCode = \"%s\"\n", input.Constant, input.Code))

	regBlock.WriteString(fmt.Sprintf("\n  %s: {\n", input.Constant))
	regBlock.WriteString(fmt.Sprintf("    Message:   \"%s\",\n", input.Message))
	regBlock.WriteString("    ExtraTags: []string{")
	for i, f := range input.Fields {
		if i > 0 {
			regBlock.WriteString(", ")
		}
		regBlock.WriteString(fmt.Sprintf("\"%s\"", f))
	}
	regBlock.WriteString("},\n  },\n")

	// 2. Perform target slicing using visual layout anchors
	if input.IsNewFamily {
		// New families append at the absolute end of the target placeholder boundary blocks
		content = strings.Replace(content, "  // X_ANCHOR_CONSTANTS_END", constBlock.String()+"  // X_ANCHOR_CONSTANTS_END", 1)
		content = strings.Replace(content, "  // X_ANCHOR_REGISTRY_END", regBlock.String()+"  // X_ANCHOR_REGISTRY_END", 1)
	} else {
		// Existing families look for the next sequential family token boundary or default to boundary tail anchors
		nextFamilyMarker := fmt.Sprintf("=== FAMILY: %d", input.FamilyID+1)

		content = injectBeforeOccurrenceOrAnchor(content, nextFamilyMarker, "  // X_ANCHOR_CONSTANTS_END", constBlock.String())
		content = injectBeforeOccurrenceOrAnchor(content, nextFamilyMarker, "  // X_ANCHOR_REGISTRY_END", regBlock.String())
	}

	return content
}

func injectBeforeOccurrenceOrAnchor(fullText, searchToken, fallbackAnchor, blockToInject string) string {
	idx := strings.Index(fullText, searchToken)
	if idx != -1 {
		// Found next family, rewind slightly to position inside comments context beautifully
		before := fullText[:idx]
		// Backtrack to find the start of the divider comment block
		lastDivider := strings.LastIndex(before, "// ============================================================================")
		if lastDivider != -1 {
			return fullText[:lastDivider] + blockToInject + fullText[lastDivider:]
		}
		return fullText[:idx] + blockToInject + fullText[idx:]
	}
	// Fallback to section tail anchor
	return strings.Replace(fullText, fallbackAnchor, blockToInject+fallbackAnchor, 1)
}
