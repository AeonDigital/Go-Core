package cmd

import (
	"fmt"
)

// HandleHelp routes the help documentation based on the user's focus.
func HandleHelp(args []string) {
	if len(args) == 0 {
		PrintGeneralHelp()
		return
	}

	target := args[0]
	switch target {
	case "error":
		PrintErrorHelp()
	case "bridge":
		PrintBridgeHelp()
	case "init":
		PrintInitHelp()
	default:
		fmt.Printf("[ERR] No specific help documentation found for '%s'.\n\n", target)
		PrintGeneralHelp()
	}
}

func PrintGeneralHelp() {
	fmt.Println("xpkg - A structural boilerplate generator and architecture management CLI")
	fmt.Println("\nUsage:")
	fmt.Println("  xpkg <resource> <sub-resource/action> [flags]")
	fmt.Println("\nAvailable Resources:")
	fmt.Println("  init          Bootstrap a brand new Go package from strict blueprints.")
	fmt.Println("  error         Manage error tokens, codes, and semantic sub-contexts.")
	fmt.Println("  bridge        Generate custom resource bridges (interfaces, structures, factories).")
	fmt.Println("  help          Show general or resource-specific usage documentation.")
	fmt.Println("\nUse 'xpkg help <resource>' for detailed flag guidelines and action definitions.")
}

func PrintInitHelp() {
	fmt.Println("Usage: xpkg init --name=<package_path_or_name>")
	fmt.Println("\nDescription:")
	fmt.Println("  Bootstraps 8 architecture compliance boilerplate files based on internal templates.")
	fmt.Println("\nFlags:")
	fmt.Println("  --name        The target directory path or name (e.g., 'my_lib' or 'internal/services/auth').")
}

func PrintErrorHelp() {
	fmt.Println("Usage: xpkg error <sub-resource> add [flags]")
	fmt.Println("\nAvailable Sub-Resources:")
	fmt.Println("  code add      Append a new corporate error code identifier inside '01_xerrors.go'.")
	fmt.Println("  subctx add    Append a new specialized error sub-context derived from the core PKGCTX.")
	fmt.Println("\nFlags for 'error code add':")
	fmt.Println("  --family       Family ID index bounds. Use 0 to activate the interactive console assistant.")
	fmt.Println("  --family-title Title for the family (required only when declaring a new sequential index inline).")
	fmt.Println("  --name         Error token key layout schema name (camelCase or snake_case).")
	fmt.Println("  --fields       Comma-separated extra tag metadata keys (e.g., 'field,value').")
	fmt.Println("  --message      Fallback text message for human consumption.")
	fmt.Println("  --tech         Strict technical/architectural scenario documentation text.")
	fmt.Println("\nFlags for 'error subctx add':")
	fmt.Println("  --name         The name of the new specialized sub-context scope (e.g., 'validation').")
}

func PrintBridgeHelp() {
	fmt.Println("Usage: xpkg bridge add [flags]")
	fmt.Println("\nDescription:")
	fmt.Println("  Generates a custom structural bridge source file ('04_bridge_<name>.go') matching structural architecture standards.")
	fmt.Println("\nFlags:")
	fmt.Println("  --name        The specific layout resource definition string identifier (e.g., 'S3', 'http').")
}
