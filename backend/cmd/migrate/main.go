package main

import (
    "context"
    "errors"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "sort"
    "time"

    "github.com/jackc/pgx/v5"
)

type Config struct {
    DatabaseURL   string
    MigrationsDir string
}

type MigrationFile struct {
    Name string
    Path string
}

func main() {
    if err := run(); err != nil {
        log.Fatal(err)
    }
}

func run() error {
    cfg := loadConfig()
    if err := validateConfig(cfg); err != nil {
        return err
    }
    conn, err := connect(cfg.DatabaseURL)
    if err != nil {
        return err
    }
    defer conn.Close(context.Background())
    if err := ensureMigrationsTable(conn); err != nil {
        return err
    }
    files, err := listMigrationFiles(cfg.MigrationsDir)
    if err != nil {
        return err
    }
    return applyFiles(conn, files)
}

func loadConfig() Config {
    cfg := Config{}
    cfg.DatabaseURL = os.Getenv("DATABASE_URL")
    cfg.MigrationsDir = getenv("MIGRATIONS_DIR", "migrations")
    return cfg
}

func validateConfig(cfg Config) error {
    if cfg.DatabaseURL == "" {
        return errors.New("DATABASE_URL is required")
    }
    return nil
}

func connect(databaseURL string) (*pgx.Conn, error) {
    parsed, err := pgx.ParseConfig(databaseURL)
    if err != nil {
        return nil, err
    }
    parsed.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    return pgx.ConnectConfig(ctx, parsed)
}

func ensureMigrationsTable(conn *pgx.Conn) error {
    query := `CREATE TABLE IF NOT EXISTS public.schema_migrations (
        version TEXT PRIMARY KEY,
        applied_at TIMESTAMP NOT NULL DEFAULT now()
    )`
    _, err := conn.Exec(context.Background(), query)
    return err
}

func listMigrationFiles(dir string) ([]MigrationFile, error) {
    pattern := filepath.Join(dir, "*.sql")
    paths, err := filepath.Glob(pattern)
    if err != nil {
        return nil, err
    }
    sort.Strings(paths)
    return toMigrationFiles(paths), nil
}

func toMigrationFiles(paths []string) []MigrationFile {
    files := make([]MigrationFile, 0, len(paths))
    for _, path := range paths {
        files = append(files, MigrationFile{Name: filepath.Base(path), Path: path})
    }
    return files
}

func applyFiles(conn *pgx.Conn, files []MigrationFile) error {
    for _, file := range files {
        if err := applyOne(conn, file); err != nil {
            return err
        }
    }
    log.Println("migrations finished")
    return nil
}

func applyOne(conn *pgx.Conn, file MigrationFile) error {
    already, err := isApplied(conn, file.Name)
    if err != nil {
        return err
    }
    if already {
        log.Printf("skip %s", file.Name)
        return nil
    }
    log.Printf("apply %s", file.Name)
    sqlText, err := os.ReadFile(file.Path)
    if err != nil {
        return err
    }
    if _, err := conn.Exec(context.Background(), string(sqlText)); err != nil {
        return fmt.Errorf("migration %s failed: %w", file.Name, err)
    }
    return markApplied(conn, file.Name)
}

func isApplied(conn *pgx.Conn, version string) (bool, error) {
    row := conn.QueryRow(context.Background(),
        `SELECT EXISTS(SELECT 1 FROM public.schema_migrations WHERE version=$1)`, version)
    exists := false
    return exists, row.Scan(&exists)
}

func markApplied(conn *pgx.Conn, version string) error {
    _, err := conn.Exec(context.Background(),
        `INSERT INTO public.schema_migrations (version) VALUES ($1)`, version)
    return err
}

func getenv(key, fallback string) string {
    value := os.Getenv(key)
    if value == "" {
        return fallback
    }
    return value
}
