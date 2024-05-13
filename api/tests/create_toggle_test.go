package tests

import (
	"context"
	"fmt"
	"github.com/RevittConsulting/sft/sft"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestCreateToggle(t *testing.T) {
	// tests:
	// normal creation
	// creation of pre-existing toggle
	// creation with incorrect fields?

	duplicateToggle := sft.ToggleDto{
		FeatureName: "duplicate feature",
		ToggleMeta: sft.ToggleMeta{
			"key 1": "value 1",
			"key 2": "value 2",
		},
		Enabled: true,
	}

	// create a duplicate entry to test against.
	initialToggle, err := sftService.CreateToggle(context.Background(), duplicateToggle)
	if err != nil {
		t.Errorf("error creating initial toggle: %s", err)
	}
	_ = initialToggle

	var tests = []struct {
		name          string
		toggle        sft.ToggleDto
		expectedError error
	}{
		{"creating a normal toggle",
			sft.ToggleDto{
				FeatureName: "test feature 1",
				ToggleMeta: sft.ToggleMeta{
					"key 1": "value 1",
					"key 2": "value 2",
				},
				Enabled: true,
			},
			nil,
		},
		{"creating a pre-existing toggle",
			duplicateToggle,
			fmt.Errorf("toggle of that name already exists"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			toggleId, err := sftService.CreateToggle(context.Background(), test.toggle)
			if test.expectedError != nil {
				if err == nil {
					t.Errorf("Expected error but got nil")
				} else if !strings.Contains(err.Error(), test.expectedError.Error()) {
					t.Errorf("Expected error %q, got %q", test.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %s", err)
				} else {
					assert.NotNil(t, toggleId)
				}
			}
		})
	}

	// clear DB of entries after this test
	err = ClearDatabase(context.Background(), dbPool)
	if err != nil {
		fmt.Println("problem clearing DB")
	}

}
