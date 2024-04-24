package sft

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// interface for the service
type ISimpleFeatureToggleDb interface {
	CheckToggleExists(ctx context.Context, featureName string) (bool, error)
	CreateToggle(ctx context.Context, toggleDto ToggleDto) (*ToggleId, error)
	DisableFeature(ctx context.Context, toggleId uuid.UUID) error
	EnableFeature(ctx context.Context, toggleId uuid.UUID) error
	GetAllToggles(ctx context.Context) ([]*Toggle, error)
}

type Service struct {
	db ISimpleFeatureToggleDb
}

// based on srt, takes a db, context, and pool, all of which that are created in the calling application, returns a service
func NewService(db ISimpleFeatureToggleDb, ctx context.Context, pool *pgxpool.Pool) *Service {
	fmt.Println("new service!")
	//if err := RunDbMigrations(ctx, pool); err != nil {
	//	panic(err)
	//}
	//
	//if err := RunDbSeed(ctx, pool); err != nil {
	//	log.Println("db seed failed")
	//}

	s := &Service{
		db: db,
	}

	return s
}

func (s *Service) CreateToggle(ctx context.Context, toggleDto ToggleDto) (*ToggleId, error) {
	return s.db.CreateToggle(ctx, toggleDto)
}

func (s *Service) DisableFeature(ctx context.Context, toggleId uuid.UUID) error {
	return s.db.DisableFeature(ctx, toggleId)
}

func (s *Service) EnableFeature(ctx context.Context, toggleId uuid.UUID) error {
	return s.db.DisableFeature(ctx, toggleId)
}

func (s *Service) GetAllToggles(ctx context.Context) ([]*Toggle, error) {
	return s.db.GetAllToggles(ctx)
}
