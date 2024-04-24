package tests

import (
	"context"
	"fmt"
	"github.com/RevittConsulting/sft/sft"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/tern/v2/migrate"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"os"
	"testing"
	"time"
)

var sftService *sft.Service
var containerCleanup func() // Global cleanup function
var dbPool *pgxpool.Pool    // Global variable for the database connection pool

// TODO: clear db between each test.

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {

	// <<<<<<< PREVIOUS CODE
	//if err := godotenv.Load(); err != nil {
	//	log.Fatal("missing .env files " + err.Error())
	//}
	//
	//dbConfig := config.GetDatabaseConfig()
	//
	//pool, err := ConnectPostgres(dbConfig)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//sftDb := sft.NewDb(pool)
	//sftService = sft.NewService(sftDb, context.Background(), pool)
	// >>>>>>> END

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	container, containerCleanup, err := CreatePostgresContainer(ctx)
	if err != nil {
		log.Fatalf("failed to setup postgres container: %v", err)
	}
	_ = containerCleanup

	dbPool, err := CreatePgxPool(ctx, container)
	if err != nil {
		log.Fatalf("failed to create pgx pool: %v", err)
	}
	//defer dbPool.Close()

	if err = RunMigrations(ctx, dbPool); err != nil {
		log.Fatalf("failed to run db migrations: %v", err)
	}

	sftDb := sft.NewDb(dbPool)
	sftService = sft.NewService(sftDb, context.Background(), dbPool)
}

func teardown() {
	if dbPool != nil {
		dbPool.Close() // Close the database pool
		log.Println("Database connection pool closed.")
	}
	if containerCleanup != nil {
		containerCleanup()
	}
}

const dbName = "postgres"
const dbUser = "user"
const dbPassword = "password"

func CreatePostgresContainer(ctx context.Context) (*postgres.PostgresContainer, func(), error) {
	postresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("docker.io/postgres:latest"),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(10*time.Minute)),
	)
	if err != nil {
		return nil, nil, err
	}

	cleanupFunc := func() {
		if err := postresContainer.Terminate(ctx); err != nil {
			panic(err)
		}
	}

	return postresContainer, cleanupFunc, nil
}

func CreatePgxPool(ctx context.Context, container *postgres.PostgresContainer) (*pgxpool.Pool, error) {
	host, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, err
	}

	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		dbUser,
		dbPassword,
		host,
		port.Int(),
		dbName,
	)

	c, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, c)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func RunMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	ternMigrator, err := migrate.NewMigrator(ctx, conn.Conn(), "sft_test_schema")
	if err != nil {
		return err

	}

	fsys := os.DirFS("../sft/db/migrations/")

	err = ternMigrator.LoadMigrations(fsys)
	if err != nil {
		return err
	}

	err = ternMigrator.Migrate(ctx)
	if err != nil {
		return err
	}

	return nil
}

//func ConnectPostgres(config *config.Config) (*pgxpool.Pool, error) {
//	ctx := context.Background()
//
//	log.Println("connecting to postgres 1")
//	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
//		config.Database.Username,
//		config.Database.Password,
//		config.Database.Host,
//		config.Database.Port,
//		"postgres",
//	)
//	c, err := pgxpool.ParseConfig(connString)
//	if err != nil {
//		return nil, err
//	}
//
//	pool, err := pgxpool.NewWithConfig(ctx, c)
//	if err != nil {
//		log.Fatal(fmt.Sprintf("unable to connect to database: %v\n", err))
//	}

// <<<<<< RANDOM STUFF
//defer pool.Close()
//
//log.Println("connecting to postgres 2")
//newConnString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
//	config.Database.Username,
//	config.Database.Password,
//	config.Database.Host,
//	config.Database.Port,
//	config.Database.TestDBName,
//)
//
//c, err = pgxpool.ParseConfig(newConnString)
//if err != nil {
//	log.Fatal(fmt.Sprintf("unable to parse pool config: %v\n", err))
//}
//
//pool, err = pgxpool.NewWithConfig(context.Background(), c)
//if err != nil {
//	log.Fatal(fmt.Sprintf("unable to connect to database: %v\n", err))
//}
// >>>>>>> END OF RANDOM STUFF

//	return pool, nil
//}
