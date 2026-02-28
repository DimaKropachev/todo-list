package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/DimaKropachev/todo-list/internal/config"
	"github.com/golang-migrate/migrate/v4"
)

func main() {
	var (
		migrationPath string
		command       string
		version       int
		steps         int
	)

	flag.StringVar(&migrationPath, "path", "./migrations", "path to migration directory")
	flag.StringVar(&command, "command", "up", "migration command: up, down, version, force")
	flag.IntVar(&version, "version", 0, "migration version (used with force)")
	flag.IntVar(&steps, "steps", 0, "number of migration steps")
	flag.Parse()

	envPath := os.Getenv("ENV_PATH")

	cfg, err := config.ParseConfig(envPath)
	if err != nil {
		log.Fatalf("error parsing config: %v\n", err)
	}

	// sslmode отвечает за шифрование
	// sslmode=disable - без шифрования, только локально или внутри docker-сети
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.DB.UserName,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.DBName,
		cfg.DB.SSLMode,
	)

	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationPath),
		connStr,
	)
	if err != nil {
		log.Fatalf("failed to create instance migrate: %v\n", err)
	}
	defer m.Close()

	switch command {
	case "up":
		if steps > 0 {
			err = m.Steps(steps)
		} else {
			err = m.Up()
		}
		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("migrations up failed: %v\n", err)
		}
		log.Println("migrations applied successfully!")

	case "down":
		if steps > 0 {
			err = m.Steps(-steps)
		} else {
			log.Fatalln("down without -steps is not allowed (dangerous)")
		}
		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("migrations down failed: %v\n", err)
		}
		log.Println("migrations rolled back successfully!")

	case "version":
		version, dirty, err := m.Version()
		if errors.Is(err, migrate.ErrNilVersion) {
			log.Println("no migrations applied yet")
			return
		}
		if err != nil {
			log.Fatalf("failed to get version: %v\n", err)
		}
		log.Printf("Current version: %d, Dirty: %v\n", version, dirty)
		if dirty {
			log.Println("WARNING: database is in dirty state")
		}

	case "force":
		if version <= 0 {
			log.Fatalln("force requires -version > 0")
		}
		if err := m.Force(version); err != nil {
			log.Fatalf("force failed: %v\n", err)
		}
		log.Printf("database version forcibly set to %d", version)

	default:
		log.Fatalf("unknown command: %s (available: up, down, version, force)\n", command)
	}
}
