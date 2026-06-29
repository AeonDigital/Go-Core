package xcli_test

import (
	"errors"
	"testing"

	"github.com/AeonDigital/Go-Core/sterr/pkg/sterr"
	"github.com/AeonDigital/Go-Core/xcli/pkg/xcli"
	"github.com/AeonDigital/Go-Core/xcli/pkg/xcli/xclistruc"
)

// TestTypedDescriptorParseAndValidateSuccessAndErrors verifies the complete single-value
// pipeline executing sequential tracks across parsing, limit validations, and disk checks.
func TestTypedDescriptorParseAndValidateSuccessAndErrors(t *testing.T) {
	// 1. Success track execution check
	descriptorSuccess := xcli.TypedDescriptor[int]{
		Parser: func(val string) (int, error) {
			return 42, nil
		},
		LimitValidator: func(flag string, val int, spec xclistruc.Flag) error {
			if val != 42 {
				return errors.New("limit mismatch simulation")
			}
			return nil
		},
		DiskValidator: func(flag string, val int, spec xclistruc.Flag) error {
			return nil // Not utilized for this primitive success track
		},
	}

	res, err := descriptorSuccess.ParseAndValidate("--port", "42", xclistruc.Flag{})
	if err != nil {
		t.Errorf("expected no validation errors for successful pipeline, got: %v", err)
	}
	if val, ok := res.(int); !ok || val != 42 {
		t.Errorf("expected parsed integer 42 wrapped in interface, got: %v", res)
	}

	// 2. Error track: Core parser failure blocker trigger
	descriptorParserErr := xcli.TypedDescriptor[int]{
		Parser: func(val string) (int, error) {
			return 0, sterr.New().SetMessage("malformed integer layout simulation")
		},
	}

	_, errParser := descriptorParserErr.ParseAndValidate("--port", "abc", xclistruc.Flag{})
	if errParser == nil {
		t.Fatalf("expected early parser error blocker execution, got nil")
	}

	// 3. Error track: Limit validator criteria failure trigger
	descriptorLimitErr := xcli.TypedDescriptor[int]{
		Parser: func(val string) (int, error) {
			return 100, nil
		},
		LimitValidator: func(flag string, val int, spec xclistruc.Flag) error {
			return sterr.New().SetMessage("value out of range simulation")
		},
	}

	_, errLimit := descriptorLimitErr.ParseAndValidate("--port", "100", xclistruc.Flag{})
	if errLimit == nil {
		t.Fatalf("expected limit validator constraint failure trigger, got nil")
	}

	// 4. Error track: Low-level physical disk capability validation failure trigger
	descriptorDiskErr := xcli.TypedDescriptor[string]{
		Parser: func(val string) (string, error) {
			return "/invalid/path", nil
		},
		DiskValidator: func(flag string, val string, spec xclistruc.Flag) error {
			return sterr.New().SetMessage("path resource does not exist simulation")
		},
	}

	_, errDisk := descriptorDiskErr.ParseAndValidate("--config", "/invalid/path", xclistruc.Flag{})
	if errDisk == nil {
		t.Fatalf("expected physical disk capability validator failure trigger, got nil")
	}
}

// TestArrayDescriptorValidateCapacityTargets verifies collection bounds evaluation
// routing length counters down to the package validation infrastructure layers.
func TestArrayDescriptorValidateCapacity(t *testing.T) {
	minItems := 2
	maxItems := 4

	// Instantiate a generic descriptor instance to gain internal visibility access
	descriptor := xcli.ArrayDescriptor[int]{}
	spec := xclistruc.Flag{
		MinItems: &minItems,
		MaxItems: &maxItems,
	}

	// 1. Success validation track
	if err := descriptor.ValidateCapacity("--ids", 3, spec); err != nil {
		t.Errorf("expected no validation errors for size within boundaries, got: %v", err)
	}

	// 2. Error track: Underflow constraint trigger (too few elements)
	errFew := descriptor.ValidateCapacity("--ids", 1, spec)
	if errFew == nil {
		t.Fatalf("expected early collection capability failure for few items, got nil")
	}

	// 3. Error track: Overflow constraint trigger (exceeding capacity limit)
	errMany := descriptor.ValidateCapacity("--ids", 5, spec)
	if errMany == nil {
		t.Fatalf("expected early collection capability failure for many items, got nil")
	}
}

// TestArrayDescriptorParseAndValidateSuccessAndErrors verifies tokenizer splitting,
// iterative primitive parsing, and item-level validation tracks for slices.
func TestArrayDescriptorParseAndValidateSuccessAndErrors(t *testing.T) {
	minItems := 2

	// 1. Success track execution check
	descriptorSuccess := xcli.ArrayDescriptor[int]{
		ElementParser: func(val string) (int, error) {
			if val == "10" {
				return 10, nil
			}
			if val == "20" {
				return 20, nil
			}
			return 0, errors.New("unexpected token layout simulation")
		},
		ItemLimitValidator: func(flag string, val int, spec xclistruc.Flag) error {
			if val < 0 {
				return errors.New("negative value limit simulation")
			}
			return nil
		},
	}

	res, err := descriptorSuccess.ParseAndValidate("--ids", "10,20", xclistruc.Flag{MinItems: &minItems})
	if err != nil {
		t.Errorf("expected no validation errors for successful slice processing, got: %v", err)
	}

	sliceRes, ok := res.([]int)
	if !ok || len(sliceRes) != 2 || sliceRes[0] != 10 || sliceRes[1] != 20 {
		t.Errorf("expected parsed integer slice []int{10, 20}, got: %v", res)
	}

	// 2.a Success track: Completely blank array scenarios handling safely
	descriptorBlank := xcli.ArrayDescriptor[int]{}
	resBlank, errBlank := descriptorBlank.ParseAndValidate("--ids", "  ", xclistruc.Flag{})
	if errBlank != nil {
		t.Errorf("expected blank layout string allocation to succeed, got: %v", errBlank)
	}
	sliceBlank, okBlank := resBlank.([]int)
	if !okBlank || len(sliceBlank) != 0 {
		t.Errorf("expected empty instantiated slice instance, got length: %d", len(sliceBlank))
	}

	// 2.b Error track: Completely blank array violating capacity criteria blocker trigger
	minItemsRequired := 1
	_, errBlankCapacity := descriptorBlank.ParseAndValidate(
		"--ids",
		"   ",
		xclistruc.Flag{MinItems: &minItemsRequired},
	)
	if errBlankCapacity == nil {
		t.Fatalf("expected blank layout capacity validation error trigger to fail early, got nil")
	}

	// 3. Error track: High-level array capacity criteria validation failure trigger
	_, errCapacity := descriptorSuccess.ParseAndValidate("--ids", "10", xclistruc.Flag{MinItems: &minItems})
	if errCapacity == nil {
		t.Fatalf("expected early boundary capacity validation failure blocker, got nil")
	}

	// 4. Error track: Core single token item parsing failure trigger
	_, errParser := descriptorSuccess.ParseAndValidate("--ids", "10,abc", xclistruc.Flag{})
	if errParser == nil {
		t.Fatalf("expected element parser execution block failure, got nil")
	}

	// 5. Error track: Item-level limit validator constraint failure trigger
	descriptorItemErr := xcli.ArrayDescriptor[int]{
		ElementParser: func(val string) (int, error) { return -5, nil },
		ItemLimitValidator: func(flag string, val int, spec xclistruc.Flag) error {
			return errors.New("negative block failure")
		},
	}
	_, errItemLimit := descriptorItemErr.ParseAndValidate("--ids", "-5", xclistruc.Flag{})
	if errItemLimit == nil {
		t.Fatalf("expected item-level validation constraint failure trigger, got nil")
	}

	// 6. Error track: Item-level physical disk capability validation failure trigger
	descriptorDiskErr := xcli.ArrayDescriptor[string]{
		ElementParser: func(val string) (string, error) { return val, nil },
		ItemDiskValidator: func(flag string, val string, spec xclistruc.Flag) error {
			return errors.New("disk resource block failure")
		},
	}
	_, errItemDisk := descriptorDiskErr.ParseAndValidate("--paths", "/ghost", xclistruc.Flag{})
	if errItemDisk == nil {
		t.Fatalf("expected item-level physical disk validator failure trigger, got nil")
	}
}
