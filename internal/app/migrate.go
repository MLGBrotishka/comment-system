package app

import (
	"errors"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	// migrate tools
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	_defaultAttempts = 20
	_defaultTimeout  = time.Second
)

func Migrate(databaseURL string) {
	databaseURL += "?sslmode=disable"

	for _, migration := range []string{"/comments", "/posts"} {
		var (
			attempts = _defaultAttempts
			err      error
			m        *migrate.Migrate
		)
		for attempts > 0 {
			m, err = migrate.New("file://migrations"+migration, databaseURL)
			if err == nil {
				break
			}

			log.Printf("Migrate: postgres is trying to connect, attempts left: %d", attempts)
			time.Sleep(_defaultTimeout)
			attempts--
		}

		if err != nil {
			log.Fatalf("Migrate: postgres connect error: %s", err)
		}

		err = m.Up()
		defer m.Close()
		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("Migrate: up error: %s", err)
		}

		if errors.Is(err, migrate.ErrNoChange) {
			log.Printf("Migrate: no change")
			return
		}

		log.Printf("Migrate: up success")
	}
}
