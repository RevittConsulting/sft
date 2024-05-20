package tests

import (
	"context"
	"fmt"
	"github.com/RevittConsulting/sft/sft"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckFeatureIsEnabled(t *testing.T) {

	scenarios := map[string]struct {
		toggle          sft.ToggleDto
		expectedEnabled bool
		expectedError   error
	}{
		"feature is enabled": {
			toggle: sft.ToggleDto{
				uuid.New(),
				"test enabled feature",
				sft.ToggleMeta{
					"key 1": "value 1",
					"key 2": "value 2",
				},
				true,
			},
			expectedEnabled: true,
			expectedError:   nil,
		},
		"feature is disabled": {
			toggle: sft.ToggleDto{
				uuid.New(),
				"test disabled feature",
				sft.ToggleMeta{
					"key 1": "value 1",
					"key 2": "value 2",
				},
				false,
			},
			expectedEnabled: false,
			expectedError:   nil,
		},
		"there is no such toggle": {
			toggle: sft.ToggleDto{
				uuid.New(),
				"no such toggle",
				sft.ToggleMeta{
					"key 1": "value 1",
					"key 2": "value 2",
				},
				false,
			},
			expectedEnabled: true,
			expectedError:   nil,
		},
	}

	for name, scenario := range scenarios {
		t.Run(name, func(t *testing.T) {
			// create toggle
			toggle, err := sftService.CreateToggle(context.Background(), scenario.toggle)
			if err != nil {
				t.Errorf("Unexpected error: %s", err)
			}

			if name == "there is no such toggle" {
				err = sftService.DeleteToggle(context.Background(), toggle.Id)
				if err != nil {
					fmt.Println("problem deleting toggle: ", err)
				}
			}
			// check enabled
			enabled, err := sftService.CheckFeatureIsEnabled(context.Background(), scenario.toggle.FeatureName)

			assert.Equal(t, scenario.expectedEnabled, enabled.Enabled)
		})
	}

}
