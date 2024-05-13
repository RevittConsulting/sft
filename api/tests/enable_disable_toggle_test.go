package tests

import (
	"context"
	"fmt"
	"github.com/RevittConsulting/sft/sft"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnableDisableToggle(t *testing.T) {

	tests := []struct {
		name               string
		toggle             sft.ToggleDto
		expectedEnableBool bool
		expectedError      error
	}{
		{
			"disabling enabled feature",
			sft.ToggleDto{
				uuid.New(),
				"originally enabled toggle",
				sft.ToggleMeta{
					"key 1": "value 1",
					"key 2": "value 2",
				},
				true,
			},
			false,
			nil,
		},
		{
			"enabling disabled feature",
			sft.ToggleDto{
				uuid.New(),
				"originally disabled toggle",
				sft.ToggleMeta{
					"key 1": "value 1",
					"key 2": "value 2",
				},
				false,
			},
			true,
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// create toggle
			toggleId, err := sftService.CreateToggle(context.Background(), test.toggle)
			if err != nil {
				t.Errorf("Unexpected error: %s", err)
			}

			// enable/disable
			if test.toggle.Enabled {
				err = sftService.DisableFeature(context.Background(), toggleId.Id)
				if err != nil {
					t.Errorf("Unexpected error: %s", err)
				}
			} else {
				err = sftService.EnableFeature(context.Background(), toggleId.Id)
				if err != nil {
					t.Errorf("Unexpected error: %s", err)
				}
			}
			// find updated toggle and check enabled bool
			allToggles, err := sftService.GetAllToggles(context.Background())
			if err != nil {
				t.Errorf("error getting all toggles to check against: %s", err)
			}
			var currentToggle *sft.Toggle

			for _, toggle := range allToggles {
				if toggle.FeatureName == test.toggle.FeatureName {
					currentToggle = toggle
				}
			}

			fmt.Printf("%+v\n", currentToggle)

			assert.Equal(t, test.expectedEnableBool, currentToggle.Enabled)

		})
	}

	// clear DB of entries after this test
	err := ClearDatabase(context.Background(), dbPool)
	if err != nil {
		fmt.Println("problem clearing DB")
	}

}
