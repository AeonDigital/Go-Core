package xcli_test

import (
	"strings"
	"testing"

	"github.com/AeonDigital/Go-Core/xcli/pkg/xcli"
	"github.com/AeonDigital/Go-Core/xcli/pkg/xcli/xcliconstt"
	"github.com/AeonDigital/Go-Core/xcli/pkg/xcli/xclistruc"
)

// TestCommandValidateAndHydrateFlagsSuccessAndErrors verifies the entire parsing,
// fallback defaults, and type mapping lifecycle directly from a Command instance.
func TestCommandValidateAndHydrateFlagsSuccessAndErrors(t *testing.T) {
	requiredFlag := true

	// Setup command metadata boundaries structure
	cmd := &xcli.Command{
		Name: "test-cmd",
		Flags: []xclistruc.Flag{
			{
				LongName:  "name",
				ShortName: "n",
				Type:      xcliconstt.TypeString,
				Required:  requiredFlag,
			},
			{
				LongName:     "port",
				Type:         xcliconstt.TypeInt,
				DefaultValue: 8080,
			},
		},
	}

	// 1. Success track execution check
	rawInputsSuccess := map[string]string{
		"name": "aeon",
		"port": "3000",
	}

	ctx, err := cmd.ValidateAndHydrateFlags(rawInputsSuccess)
	if err != nil {
		t.Errorf("expected no validation errors for valid flag inputs, got: %v", err)
	}

	// Utilizing authentic typed getters from the flagValues definition
	if nameVal := ctx.GetString("name"); nameVal != "aeon" {
		t.Errorf("expected string 'aeon', got: %v", nameVal)
	}
	if portVal := ctx.GetInt("port"); portVal != 3000 {
		t.Errorf("expected integer 3000, got: %v", portVal)
	}

	// 2. Success track: Fallback default values allocation logic check
	rawInputsDefault := map[string]string{
		"n": "aeon", // Testing short name lookup path as well
	}

	ctxDef, errDef := cmd.ValidateAndHydrateFlags(rawInputsDefault)
	if errDef != nil {
		t.Errorf("expected success utilizing default allocation paths, got: %v", errDef)
	}
	if portDef := ctxDef.GetInt("port"); portDef != 8080 {
		t.Errorf("expected fallback default value 8080, got: %v", portDef)
	}

	// 3. Error track: Mandatory presence requirement enforcement checker
	rawInputsMissing := map[string]string{
		"port": "3000",
	}

	_, errMissing := cmd.ValidateAndHydrateFlags(rawInputsMissing)
	if errMissing == nil || !strings.Contains(errMissing.Error(), "required") {
		t.Errorf("expected missing required flag error block, got: %v", errMissing)
	}

	// 4. Error track: Type transformation schema mismatch failure trigger
	rawInputsBadType := map[string]string{
		"name": "aeon",
		"port": "not-an-int",
	}

	_, errBadType := cmd.ValidateAndHydrateFlags(rawInputsBadType)
	if errBadType == nil {
		t.Fatalf("expected parsing transformation error block trigger, got nil")
	}

	// 5. Error track: Unregistered or zombie flag tokens security block
	rawInputsZombie := map[string]string{
		"name":   "aeon",
		"zombie": "undead-value",
	}

	_, errZombie := cmd.ValidateAndHydrateFlags(rawInputsZombie)
	if errZombie == nil || !strings.Contains(errZombie.Error(), "unrecognized flag identifier token") {
		t.Errorf("expected zombie variable protection trigger, got: %v", errZombie)
	}

	// 6. Error track: Defensive unmapped or unsupported flag types lookup
	cmdBrokenType := &xcli.Command{
		Name: "broken-cmd",
		Flags: []xclistruc.Flag{
			{
				LongName: "bad-flag",
				Type:     xcliconstt.FlagType("ghost_type"),
			},
		},
	}
	rawInputsBadTypeEnum := map[string]string{
		"bad-flag": "value",
	}

	_, errBadEnum := cmdBrokenType.ValidateAndHydrateFlags(rawInputsBadTypeEnum)
	if errBadEnum == nil || !strings.Contains(errBadEnum.Error(), "unsupported flag type") {
		t.Errorf("expected unmapped type validation blocker trigger, got: %v", errBadEnum)
	}
}

// TestCommandTriggerHelpEnforcesFullVisualRender verifies visual documentation tables
// mapping usage definitions, registered subcommands, and multi-layered flag guidelines.
func TestCommandTriggerHelpEnforcesFullVisualRender(t *testing.T) {
	requiredFlag := true

	// 1. Setup a complex structured command architecture to trigger all internal IF branches
	cmd := &xcli.Command{
		Name:             "main-app",
		ShortDescription: "Core engine initialization",
		LongDescription:  "Detailed technical scope documentation for the terminal operator.",
		Subcommands: map[string]*xcli.Command{
			"start": {
				Name:             "start",
				ShortDescription: "Boot up system daemons",
			},
		},
		Flags: []xclistruc.Flag{
			{
				LongName:         "config",
				ShortName:        "c",
				Type:             xcliconstt.TypeString,
				Required:         requiredFlag,
				ShortDescription: "Target file path layout configuration",
			},
			{
				LongName:         "workers",
				Type:             xcliconstt.TypeInt,
				DefaultValue:     4,
				ShortDescription: "Total dynamic execution threads count",
			},
		},
	}

	// 2. Execute the help render trace. Since it targets standard output bridges,
	// this call verifies that the formatting engine terminates successfully without panicking.
	err := cmd.TriggerHelp()
	if err != nil {
		t.Errorf("expected visual menu documentation layout to render without failures, got: %v", err)
	}

	// 3. Technical Edge Case: Enforce fallback coverage when LongDescription is absent
	cmdFallback := &xcli.Command{
		Name:             "fallback-app",
		ShortDescription: "Simple diagnostic description tool",
	}

	errFallback := cmdFallback.TriggerHelp()
	if errFallback != nil {
		t.Errorf("expected description fallback tracking flow to execute without errors, got: %v", errFallback)
	}
}
