package tests

import (
	"context"
	"fmt"
	"github.com/RevittConsulting/sft/config"
	"github.com/RevittConsulting/sft/sft"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
)

var sftService *sft.Service

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("missing .env files " + err.Error())
	}

	dbConfig := config.GetDatabaseConfig()

	pool, err := ConnectPostgres(dbConfig)
	if err != nil {
		log.Fatal(err)
	}

	sftDb := sft.NewDb(pool)
	sftService = sft.NewService(sftDb, context.Background(), pool)

}

func ConnectPostgres(config *config.Config) (*pgxpool.Pool, error) {
	ctx := context.Background()

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		config.Database.Username,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		"postgres",
	)
	c, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, c)
	if err != nil {
		log.Fatal(fmt.Sprintf("unable to connect to database: %v\n", err))
	}
	defer pool.Close()

	newConnString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		config.Database.Username,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.TestDBName,
	)

	c, err = pgxpool.ParseConfig(newConnString)
	if err != nil {
		log.Fatal(fmt.Sprintf("unable to parse pool config: %v\n", err))
	}

	pool, err = pgxpool.NewWithConfig(context.Background(), c)
	if err != nil {
		log.Fatal(fmt.Sprintf("unable to connect to database: %v\n", err))
	}

	return pool, nil
}
