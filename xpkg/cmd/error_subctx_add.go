package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// SubctxInput stores the sanitized metadata required to append a new error sub-context.
type SubctxInput struct {
	BaseConstant string // The parent/base constant found (e.g., XERR_PKGCTX)
	BaseValue    string // The parent/base string value found (e.g., ERR_XJSON)
	NewConstant  string // The final computed new constant (e.g., XERR_PKGCTX_VALIDATION)
	NewValue     string // The final computed new string value (e.g., ERR_XJSON_VALIDATION)
}

// HandleAddErrorSubctx manages the entrypoint execution flow for the 'error subctx add' subcommand.
func HandleAddErrorSubctx(args []string) {
	fs := flag.NewFlagSet("error-subctx-add", flag.ExitOnError)
	flagName := fs.String("name", "", "The name of the new sub-context (e.g., validation, auth)")

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

	// 2. Discover the baseline XERR_PKGCTX configuration
	baseConstant, baseValue := parseMainContext(fileContent)
	if baseConstant == "" || baseValue == "" {
		fmt.Println("[ERR] Could not locate the base XERR_PKGCTX definition or its string value inside the target file.")
		os.Exit(1)
	}

	// 3. Resolve parameters (Inline vs Interactive Assistant Flow)
	rawSubName := *flagName
	if rawSubName == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter Sub-Context Name (e.g., validation, auth): ")
		rawSubName, _ = reader.ReadString('\n')
	}

	// 4. Sanitize and normalize the input using your precise casing rules
	normalizedSub := toSnakeUpperCase(rawSubName)
	if normalizedSub == "" {
		fmt.Println("[ERR] Sub-context name declaration cannot resolve to an empty value state.")
		os.Exit(1)
	}

	// 5. Build dynamic outputs
	input := SubctxInput{
		BaseConstant: baseConstant,
		BaseValue:    baseValue,
		NewConstant:  fmt.Sprintf("%s_%s", baseConstant, normalizedSub),
		NewValue:     fmt.Sprintf("%s_%s", baseValue, normalizedSub),
	}

	// Validate token syntax and length constraints
	if !validateToken(input.NewConstant) {
		fmt.Printf("[ERR] Computed constant '%s' violates strict constraints (A-Z base, max 32 chars).\n", input.NewConstant)
		os.Exit(1)
	}

	// Check for duplication collisions inside the entire registry file
	if strings.Contains(fileContent, " "+input.NewConstant+" ") || strings.Contains(fileContent, "\t"+input.NewConstant+" ") {
		fmt.Printf("[ERR] Sub-context constant '%s' already exists inside registry. Aborting execution.\n", input.NewConstant)
		os.Exit(1)
	}

	// 6. Inject the brand new sub-context block back into the physical file string
	updatedContent := injectNewSubctxIntoString(fileContent, &input)

	err = os.WriteFile(targetFile, []byte(updatedContent), 0644)
	if err != nil {
		fmt.Printf("[ERR] Failed writing modifications down to '%s': %v\n", targetFile, err)
		os.Exit(1)
	}

	fmt.Printf("[OKK] Successfully appended sub-context '%s' (%s).\n", input.NewConstant, input.NewValue)
}

// parseMainContext extracts the original baseline package context constant and its value using regular expressions.
func parseMainContext(content string) (constant string, value string) {
	// Captures pattern: XERR_PKGCTX xerrors.ErrorCode = "ERR_XYZ"
	// Allowing for optional spaces/tabs around assignments
	re := regexp.MustCompile(`(XERR_[A-Z0-9_]*)\s+xerrors\.ErrorCode\s*=\s*"([^"]+)"`)
	matches := re.FindStringSubmatch(content)

	if len(matches) == 3 {
		return matches[1], matches[2]
	}
	return "", ""
}

// injectNewSubctxIntoString prepares the text block and slices it cleanly into the context anchor points.
func injectNewSubctxIntoString(content string, input *SubctxInput) string {
	var subctxBlock strings.Builder

	subctxBlock.WriteString(fmt.Sprintf("\n  // %s defines a specialized sub-context scope derived from %s.\n", input.NewConstant, input.BaseConstant))
	subctxBlock.WriteString(fmt.Sprintf("  %s xerrors.ErrorCode = \"%s\"\n", input.NewConstant, input.NewValue))

	// Injects beautifully right before the visual section boundary marker anchor
	fallbackAnchor := "  // X_ANCHOR_PKGCTX_END"
	content = strings.Replace(content, fallbackAnchor, subctxBlock.String()+fallbackAnchor, 1)

	return content
}
