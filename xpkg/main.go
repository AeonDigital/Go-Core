package main

import (
	"fmt"
	"os"

	"github.com/AeonDigital/Go-Core/xpkg/cmd"
)

func main() {
	cmd.TemplateFS = TemplateFS

	if len(os.Args) < 2 {
		cmd.PrintGeneralHelp()
		os.Exit(1)
	}

	resource := os.Args[1]

	switch resource {
	case "help":
		cmd.HandleHelp(os.Args[2:])

	case "init":
		cmd.HandleInit(os.Args[2:])

	case "bridge":
		if len(os.Args) < 3 {
			fmt.Println("[ERR] Missing action for resource 'bridge'. Expected 'add'.")
			fmt.Println("Run 'xpkg help bridge' for usage documentation.")
			os.Exit(1)
		}
		action := os.Args[2]
		if action == "add" {
			cmd.HandleAddBridge(os.Args[3:])
		} else {
			fmt.Printf("[ERR] Unknown action '%s' for resource 'bridge'.\n", action)
			os.Exit(1)
		}

	case "error":
		if len(os.Args) < 4 {
			fmt.Println("[ERR] Missing parameters for 'error'. Expected 'error code add' or 'error subctx add'.")
			fmt.Println("Run 'xpkg help error' for usage documentation.")
			os.Exit(1)
		}

		subResource := os.Args[2]
		action := os.Args[3]

		switch subResource {
		case "code":
			if action == "add" {
				cmd.HandleAddErrorCode(os.Args[4:])
			} else {
				fmt.Printf("[ERR] Unknown action '%s' for 'error code'. Expected 'add'.\n", action)
				os.Exit(1)
			}

		case "subctx":
			if action == "add" {
				cmd.HandleAddErrorSubctx(os.Args[4:])
			} else {
				fmt.Printf("[ERR] Unknown action '%s' for 'error subctx'. Expected 'add'.\n", action)
				os.Exit(1)
			}

		default:
			fmt.Printf("[ERR] Unknown sub-resource '%s' for 'error'. Use 'code' or 'subctx'.\n", subResource)
			os.Exit(1)
		}

	default:
		fmt.Printf("[ERR] Unknown resource '%s'.\n\n", resource)
		cmd.PrintGeneralHelp()
		os.Exit(1)
	}
}
