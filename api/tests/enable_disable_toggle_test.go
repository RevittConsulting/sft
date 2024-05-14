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

	scenarios := map[string]struct {
		toggle             sft.ToggleDto
		expectedEnableBool bool
		expectedError      error
	}{

		"disabling enabled feature": {
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

		"enabling disabled feature": {
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

	for name, scenario := range scenarios {
		t.Run(name, func(t *testing.T) {
			// create toggle
			toggleId, err := sftService.CreateToggle(context.Background(), scenario.toggle)
			if err != nil {
				t.Errorf("Unexpected error: %s", err)
			}

			// enable/disable
			err = sftService.ToggleFeature(context.Background(), toggleId.Id)
			if err != nil {
				t.Errorf("Unexpected error: %s", err)
			}
			// find updated toggle and check enabled bool
			allToggles, err := sftService.GetAllToggles(context.Background())
			if err != nil {
				t.Errorf("error getting all toggles to check against: %s", err)
			}
			var currentToggle *sft.Toggle

			for _, toggle := range allToggles {
				if toggle.FeatureName == scenario.toggle.FeatureName {
					currentToggle = toggle
				}
			}

			fmt.Printf("%+v\n", currentToggle)

			assert.Equal(t, scenario.expectedEnableBool, currentToggle.Enabled)

		})
	}

	// clear DB of entries after this test
	err := ClearDatabase(context.Background(), dbPool)
	if err != nil {
		t.Errorf("problem clearing DB")
	}

}
