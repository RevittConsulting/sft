package sft

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

// interface for the service
type ISimpleFeatureToggleDb interface {
	TestDbFunc()
	CheckToggleExists(ctx context.Context, featureName string) (bool, error)
	GetAllToggles(ctx context.Context) ([]*Toggle, error)
}

type Service struct {
	db ISimpleFeatureToggleDb
}

// based on srt, takes a db, context, and pool, all of which that are created in the calling application, returns a service
func NewService(db ISimpleFeatureToggleDb, ctx context.Context, pool *pgxpool.Pool) *Service {
	fmt.Println("new service!")
	if err := RunDbMigrations(ctx, pool); err != nil {
		panic(err)
	}

	if err := RunDbSeed(ctx, pool); err != nil {
		log.Println("db seed failed")
	}

	s := &Service{
		db: db,
	}

	return s
}

func (s *Service) CreateToggle(ctx context.Context, featureName string, enabled bool) error {
	// finish off create toggle
	// check whether toggle already exists
	// if not, create!

	exists, err := s.db.CheckToggleExists(ctx, featureName)
	if err != nil {
		return fmt.Errorf("error checking toggle exists: %w", err)
	}

	if exists {
		return fmt.Errorf("toggle already exists for feature: %s", featureName)
	}

	// TODO: add create toggle functionality

	return nil
}

func (s *Service) GetAllToggles(ctx context.Context) ([]*Toggle, error) {
	return s.db.GetAllToggles(ctx)
}
