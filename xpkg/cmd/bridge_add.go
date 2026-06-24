package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// BridgeConfig carries the variables to render a custom resource bridge.
type BridgeConfig struct {
	PkgName       string
	ResourceLower string // Usado apenas no cabeçalho do comentário técnico
	InterfaceName string // e.g., "IS3Bridge" ou "IHttpBridge"
	StructName    string // e.g., "sS3Bridge" ou "shttpBridge"
	FactoryName   string // e.g., "NewS3" ou "NewHttp"
	VariableName  string // e.g., "S3" ou "Http"
}

// HandleAddBridge manages the entrypoint execution flow for 'bridge add'.
func HandleAddBridge(args []string) {
	fs := flag.NewFlagSet("bridge-add", flag.ExitOnError)
	flagName := fs.String("name", "", "The specific resource name for the bridge (e.g., http, S3, MyService)")

	_ = fs.Parse(args)

	// 1. Discover current package name from the directory
	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("[ERR] Failed to resolve current working directory: %v\n", err)
		os.Exit(1)
	}
	currentPkgName := strings.ToLower(filepath.Base(workingDir))

	// 2. Resolve the resource name (Inline vs Interactive Assistant Flow)
	rawResource := strings.TrimSpace(*flagName)
	if rawResource == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter Bridge Resource Name (e.g., http, S3, MyService): ")
		rawResource, _ = reader.ReadString('\n')
		rawResource = strings.TrimSpace(rawResource)
	}

	if rawResource == "" {
		fmt.Println("[ERR] Bridge resource name declaration cannot resolve to empty value states.")
		os.Exit(1)
	}

	// 3. Apply your precise casing rules
	filenameLower := strings.ToLower(rawResource)
	targetFilename := fmt.Sprintf("04_bridge_%s.go", filenameLower)

	// Capitalize only the first letter of whatever the user typed for public objects
	resourceCapitalized := strings.ToUpper(rawResource[:1]) + rawResource[1:]

	// Safety check: verify if this specific bridge already exists
	if _, err := os.Stat(targetFilename); !os.IsNotExist(err) {
		fmt.Printf("[ERR] Target bridge file '%s' already exists. Aborting execution to protect data.\n", targetFilename)
		os.Exit(1)
	}

	bConfig := BridgeConfig{
		PkgName:       currentPkgName,
		ResourceLower: filenameLower,
		InterfaceName: fmt.Sprintf("I%sBridge", resourceCapitalized),
		StructName:    fmt.Sprintf("s%sBridge", rawResource),
		FactoryName:   fmt.Sprintf("New%s", resourceCapitalized),
		VariableName:  resourceCapitalized,
	}

	// 4. Read template from embedded filesystem
	tmplPath := "templates/04_bridge_custom.tmpl"
	tmplContent, err := TemplateFS.ReadFile(tmplPath)
	if err != nil {
		fmt.Printf("[ERR] Error reading embedded template file '%s': %v\n", tmplPath, err)
		os.Exit(1)
	}

	// 5. Evaluate and inject the concrete blueprint template
	tmpl, err := template.New("bridge_file").Parse(string(tmplContent))
	if err != nil {
		fmt.Printf("[ERR] Failed parsing template syntax engine: %v\n", err)
		os.Exit(1)
	}

	file, err := os.Create(targetFilename)
	if err != nil {
		fmt.Printf("[ERR] Failed to create target file '%s': %v\n", targetFilename, err)
		os.Exit(1)
	}
	defer file.Close()

	err = tmpl.Execute(file, bConfig)
	if err != nil {
		fmt.Printf("[ERR] Error executing blueprint evaluation for '%s': %v\n", targetFilename, err)
		os.Exit(1)
	}

	fmt.Printf("[OKK] Successfully generated custom structural bridge: %s\n", targetFilename)
}
