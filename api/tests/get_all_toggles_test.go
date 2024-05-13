package tests

import (
	"context"
	"github.com/RevittConsulting/sft/sft"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAllTestToggles(t *testing.T) {

	// create some toggles
	var toggles = []sft.ToggleDto{
		{
			Id:          uuid.New(),
			FeatureName: "toggle 1",
			ToggleMeta: sft.ToggleMeta{
				"key 1": "value 1",
				"key 2": "value 2",
			},
			Enabled: true,
		},
		{
			Id:          uuid.New(),
			FeatureName: "toggle 2",
			ToggleMeta: sft.ToggleMeta{
				"key 1": "value 1",
				"key 2": "value 2",
			},
			Enabled: true,
		},
		{
			Id:          uuid.New(),
			FeatureName: "toggle 3",
			ToggleMeta: sft.ToggleMeta{
				"key 1": "value 1",
				"key 2": "value 2",
			},
			Enabled: true,
		},
	}

	for _, toggle := range toggles {
		_, err := sftService.CreateToggle(context.Background(), toggle)
		if err != nil {
			t.Fatal("Failed to add initial toggles to db")
		}
	}

	allTogglesResponse, err := sftService.GetAllToggles(context.Background())

	_ = toggles
	assert.NoError(t, err)
	assert.Equal(t, len(toggles), len(allTogglesResponse))

	// clear DB of entries after this test
	err = ClearDatabase(context.Background(), dbPool)
	if err != nil {
		t.Errorf("problem clearing DB")
	}

}
