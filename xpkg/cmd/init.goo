package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Config carries the dynamic variables used inside the boilerplate templates.
type Config struct {
	PkgName string
}

// HandleInit manages the entrypoint execution flow for the 'init' subcommand.
func HandleInit(args []string) {
	fs := flag.NewFlagSet("init", flag.ExitOnError)
	pkgNameFlag := fs.String("name", "", "The path/name of the new Go package (lowercase, alphanumeric, e.g., internal/services/my_lib)")

	_ = fs.Parse(args)

	if *pkgNameFlag == "" {
		fmt.Println("[ERR] The '--name' parameter is required.")
		fmt.Println("Usage: xpkg init --name=internal/services/my_awesome_lib")
		os.Exit(1)
	}

	// Clean the path to handle cross-platform separators neatly
	targetPath := filepath.Clean(strings.TrimSpace(*pkgNameFlag))

	// Extract ONLY the last element of the path to act as the valid Go package name
	// E.g., "internal/services/auth" -> "auth"
	pkgName := strings.ToLower(filepath.Base(targetPath))

	if pkgName == "." || pkgName == "/" {
		fmt.Println("[ERR] Invalid package path configuration template resolved.")
		os.Exit(1)
	}

	config := Config{PkgName: pkgName}

	// 1. Safety Check: Abort if the target package directory already exists
	if _, err := os.Stat(targetPath); !os.IsNotExist(err) {
		fmt.Printf("[ERR] target directory '%s' already exists. Aborting execution to protect existing data.\n", targetPath)
		os.Exit(1)
	}

	// 2. Create target library directory structure recursively
	if err := os.MkdirAll(targetPath, 0755); err != nil {
		fmt.Printf("[ERR] Error creating target directory structure: %v\n", err)
		os.Exit(1)
	}

	// Define the strict file mapping between template sources and their final destinations inside targetPath
	filesToGenerate := map[string]string{
		"templates/00_config.tmpl":       filepath.Join(targetPath, "00_config.go"),
		"templates/01_xerrors.tmpl":      filepath.Join(targetPath, "01_xerrors.go"),
		"templates/02_definitions.tmpl":  filepath.Join(targetPath, "02_definitions.go"),
		"templates/03_functions.tmpl":    filepath.Join(targetPath, "03_functions.go"),
		"templates/04_bridge.tmpl":       filepath.Join(targetPath, "04_bridge.go"),
		"templates/05_private_test.tmpl": filepath.Join(targetPath, "05_private_test.go"),
		"templates/MD_STYLE.tmpl":        filepath.Join(targetPath, "MD_STYLE.md"),
		"templates/README.tmpl":          filepath.Join(targetPath, "README.md"),
	}

	// Setup template helper functions (e.g., uppercase conversion)
	funcMap := template.FuncMap{
		"stringsToUpper": strings.ToUpper,
	}

	// 3. Read, parse, and evaluate each template file sequentially
	for tmplPath, finalFileDestination := range filesToGenerate {
		tmplContent, err := TemplateFS.ReadFile(tmplPath)
		if err != nil {
			fmt.Printf("[ERR] Error reading embedded template file '%s': %v\n", tmplPath, err)
			os.Exit(1)
		}

		tmpl, err := template.New(tmplPath).Funcs(funcMap).Parse(string(tmplContent))
		if err != nil {
			fmt.Printf("[ERR] Error parsing template syntax for '%s': %v\n", tmplPath, err)
			os.Exit(1)
		}

		file, err := os.Create(finalFileDestination)
		if err != nil {
			fmt.Printf("[ERR] Error creating target file '%s': %v\n", finalFileDestination, err)
			os.Exit(1)
		}

		err = tmpl.Execute(file, config)
		file.Close()

		if err != nil {
			fmt.Printf("[ERR] Error executing blueprint evaluation for '%s': %v\n", finalFileDestination, err)
			os.Exit(1)
		}
		fmt.Printf("  └─ Generated: %s\n", finalFileDestination)
	}

	fmt.Printf("\n[OKK] Package '%s' successfully bootstrapped inside directory '%s' via xpkg.\n", pkgName, targetPath)
}
