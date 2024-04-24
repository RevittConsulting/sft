package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAllTestToggles(t *testing.T) {

	// TODO: create some toggles here rather than relying on seeded data

	toggles, err := sftService.GetAllToggles(context.Background())

	_ = toggles
	assert.NoError(t, err)
	assert.Equal(t, 2, len(toggles))

}
