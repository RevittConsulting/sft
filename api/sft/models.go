package sft

import "github.com/google/uuid"

type ToggleMeta map[string]interface{}

type Toggle struct {
	Id          uuid.UUID  `db:"id" json:"id"`
	FeatureName string     `db:"feature_name" json:"feature_name"`
	ToggleMeta  ToggleMeta `db:"toggle_meta" json:"toggle_meta"`
	Enabled     bool       `db:"enabled" json:"enabled"`
}

type ToggleDto struct {
	Id          uuid.UUID  `db:"id" json:"id"`
	FeatureName string     `db:"feature_name" json:"feature_name"`
	ToggleMeta  ToggleMeta `db:"toggle_meta" json:"toggle_meta"`
	Enabled     bool       `db:"enabled" json:"enabled"`
}

type ToggleId struct {
	Id uuid.UUID `db:"id" json:"id"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
