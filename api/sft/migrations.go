package sft

import (
	"context"
	"embed"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/tern/v2/migrate"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

//go:embed db/migrations/*.sql
var migrationFiles embed.FS

func RunDbMigrations(ctx context.Context, pool *pgxpool.Pool) error {

	// acquire connection
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	// create tern migrator
	ternMigrator, err := migrate.NewMigrator(ctx, conn.Conn(), "public.sft_schema_version")
	if err != nil {
		return fmt.Errorf("failed to acquire tern migrator: %w", err)
	}

	// identify migration directory
	migrationRoot, _ := fs.Sub(migrationFiles, "db/migrations")

	// Load the migrations from the directory
	if err = ternMigrator.LoadMigrations(migrationRoot); err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	// what is this doing?
	m := len(ternMigrator.Migrations)
	_ = m

	// Run the migrations
	if err = ternMigrator.Migrate(ctx); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// add a callback to OnStart, per the tern Migrate file
	ternMigrator.OnStart = func(seq int32, name string, directionName string, sql string) {
		log.Println(fmt.Sprintf("starting migration %d/%d %s %s", seq, m, name, directionName))
	}

	return nil
}

// Function to dynamically get the path to the seed files
func GetSeedFilesPath() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("cannot get current file path")
	}
	fmt.Println("calling file path: ", filename)
	dir := filepath.Dir(filename)

	fmt.Println("filepath: ", filepath.Join(dir, "..", "db", "seed"))

	return filepath.Join(dir, "..", "sft", "db", "seed"), nil
}

func RunDbSeed(ctx context.Context, pool *pgxpool.Pool) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	seedPath, err := GetSeedFilesPath()
	if err != nil {
		return fmt.Errorf("failed to get seed files path: %w", err)
	}

	fsys := os.DirFS(seedPath)

	err = runSqlFiles(ctx, conn, fsys)
	if err != nil {
		return fmt.Errorf("failed to run sql seed files: %w", err)
	}

	return nil
}

func runSqlFiles(ctx context.Context, conn *pgxpool.Conn, fsys fs.FS) error {
	fmt.Println("file path: ", fsys)

	files, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return fmt.Errorf("error reading seed files: %w", err)
	}

	fmt.Println("number of files: ", len(files))

	for _, file := range files {

		fmt.Println("file name: ", file.Name())

		if file.IsDir() {
			continue
		}

		f, err := fsys.Open(file.Name())
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}

		b, err := io.ReadAll(f)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		_, err = conn.Exec(ctx, string(b))
		if err != nil {
			return fmt.Errorf("failed to execute file: %w", err)
		}
	}
	log.Println("sft db seed complete")
	return nil
}
