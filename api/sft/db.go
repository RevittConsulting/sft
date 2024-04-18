package sft

import (
	"context"
	"fmt"
	"github.com/RevittConsulting/sft/sft/utils"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Db struct {
	pool *pgxpool.Pool
}

func NewDb(pool *pgxpool.Pool) *Db {
	return &Db{pool: pool}
}

func (db *Db) TestDbFunc() {
	fmt.Println("hello from the db func")
}

func (db *Db) CheckToggleExists(ctx context.Context, featureName string) (bool, error) {
	var exists bool
	sql := `select exists (select 1 from sft.feature_toggles where feature_name = $1)`

	err := pgxscan.Get(ctx, db.pool, &exists, sql, featureName)

	if err != nil {
		return false, err
	}
	return exists, nil
}

func (db *Db) GetAllToggles(ctx context.Context) ([]*Toggle, error) {
	sql := `select * from sft.feature_toggles`

	tx, err := utils.TxBegin(ctx, db.pool)
	defer utils.TxDefer(tx, ctx)

	if err != nil {
		return nil, err
	}

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
