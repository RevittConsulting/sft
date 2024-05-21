package sft

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

// interface for the service
type ISimpleFeatureToggleDb interface {
	CheckToggleExists(ctx context.Context, featureName string) (bool, error)
	CreateToggle(ctx context.Context, toggleDto ToggleDto) (*ToggleId, error)
	ToggleFeature(ctx context.Context, toggleId uuid.UUID) error
	GetAllToggles(ctx context.Context) ([]*Toggle, error)
	DeleteToggle(ctx context.Context, toggleId uuid.UUID) error
	CheckFeatureIsEnabled(ctx context.Context, featureName string) (*Enabled, error)
}

type Service struct {
	db ISimpleFeatureToggleDb
}

// based on srt, takes a db, context, and pool, all of which that are created in the calling application, returns a service
func NewService(db ISimpleFeatureToggleDb, ctx context.Context, pool *pgxpool.Pool) *Service {

	s := &Service{
		db: db,
	}

	err := RunDbMigrations(ctx, pool)
	if err != nil {
		log.Println("Failed to run db migrations")
	}

	return s
}

func (s *Service) CreateToggle(ctx context.Context, toggleDto ToggleDto) (*ToggleId, error) {
	return s.db.CreateToggle(ctx, toggleDto)
}

func (s *Service) ToggleFeature(ctx context.Context, toggleId uuid.UUID) error {
	return s.db.ToggleFeature(ctx, toggleId)
}

func (s *Service) GetAllToggles(ctx context.Context) ([]*Toggle, error) {
	return s.db.GetAllToggles(ctx)
}

func (s *Service) DeleteToggle(ctx context.Context, toggleId uuid.UUID) error {
	return s.db.DeleteToggle(ctx, toggleId)
}

func (s *Service) CheckFeatureIsEnabled(ctx context.Context, featureName string) (*Enabled, error) {
	log.Println("checking feature is enabled: ", featureName)
	return s.db.CheckFeatureIsEnabled(ctx, featureName)
}
