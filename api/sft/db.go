package sft

import (
	"context"
	"fmt"
	"github.com/RevittConsulting/sft/sft/utils"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Db struct {
	pool *pgxpool.Pool
}

func NewDb(pool *pgxpool.Pool) *Db {
	return &Db{pool: pool}
}

func (db *Db) CheckToggleExists(ctx context.Context, featureName string) (bool, error) {
	var exists bool
	sql := `select exists (select 1 from sft.feature_toggles where feature_name = $1)`

	err := pgxscan.Get(ctx, db.pool, &exists, sql, featureName)

	if err != nil {
		return false, fmt.Errorf("error checking whether toggle exists: %w", err)
	}
	return exists, nil
}

func (db *Db) CreateToggle(ctx context.Context, toggleDto ToggleDto) (*ToggleId, error) {

	tx, err := utils.TxBegin(ctx, db.pool)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	exists, err := db.CheckToggleExists(ctx, toggleDto.FeatureName)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("toggle of that name already exists")
	}

	sql := `insert into sft.feature_toggles (feature_name, toggle_meta, enabled) values ($1, $2, $3) returning id`

	toggleId := &ToggleId{}

	err = tx.QueryRow(ctx, sql, toggleDto.FeatureName, toggleDto.ToggleMeta, toggleDto.Enabled).Scan(&toggleId.Id)
	if err != nil {
		return nil, fmt.Errorf("error inserting toggle: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return toggleId, nil
}

func (db *Db) ToggleFeature(ctx context.Context, toggleId uuid.UUID) error {
	sql := `update sft.feature_toggles set enabled = NOT enabled where id = $1`

	tx, err := utils.TxBegin(ctx, db.pool)
	if err != nil {
		return err
	}
	defer utils.TxDefer(tx, ctx)

	result, err := tx.Exec(ctx, sql, toggleId)
	if err != nil {
		return fmt.Errorf("error updating toggle: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no rows updated")
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

func (db *Db) GetAllToggles(ctx context.Context) ([]*Toggle, error) {
	sql := `select id, feature_name, toggle_meta, enabled from sft.feature_toggles`

	tx, err := utils.TxBegin(ctx, db.pool)
	if err != nil {
		return nil, err
	}
	defer utils.TxDefer(tx, ctx)

	var toggles []*Toggle
	var rows pgx.Rows

	// Note: the query gets sent to rows first, then this is sent to our slice vis pgxscan
	rows, err = tx.Query(ctx, sql)
	if err != nil {
		return nil, err
	}

	err = pgxscan.ScanAll(&toggles, rows)
	if err != nil {
		return nil, err
	}

	return toggles, nil
}

func (db *Db) DeleteToggle(ctx context.Context, toggleId uuid.UUID) error {
	sql := `delete from sft.feature_toggles where id = $1`

	tx, err := utils.TxBegin(ctx, db.pool)
	if err != nil {
		return err
	}
	defer utils.TxDefer(tx, ctx)

	result, err := tx.Exec(ctx, sql, toggleId)
	_ = result
	if err != nil {
		return fmt.Errorf("error removing toggle: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error commiting db transaction: %w", err)
	}

	return nil

}

func (db *Db) CheckFeatureIsEnabled(ctx context.Context, featureName string) (*Enabled, error) {

	enabled := &Enabled{}

	// I've added this so that if it errors out before completion, a feature isn't automatically disabled.
	// TODO: maybe this isn't the best approach - come up with a better one.

	enabled.Enabled = true
	// find toggle by feature name
	sql := `select enabled from sft.feature_toggles where feature_name = $1`

	tx, err := utils.TxBegin(ctx, db.pool)
	if err != nil {
		return enabled, err
	}
	defer utils.TxDefer(tx, ctx)

	row := tx.QueryRow(ctx, sql, featureName)

	err = row.Scan(&enabled.Enabled)
	if err != nil {
		if err == pgx.ErrNoRows {
			enabled.Enabled = true
			return enabled, nil
		}
		return enabled, nil
	}

	return enabled, nil
}

//// toggle check code for within parent:
//enabled, err := h.sft.CheckFeatureIsEnabled(r.Context(), "Add task")
//if err != nil {
//log.Println("error checking feature: ", err)
//}
//if enabled.Enabled != true {
//return
//}
