package tests

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/webdevelop-pro/go-common/configurator"
	"github.com/webdevelop-pro/lib/db"
)

type FixturesManager struct {
	fixtures map[string]string
	db       *db.DB
}

func NewFixturesManager(fixturesList map[string]string) FixturesManager {
	cfg := db.Config{}

	// Fix for timezones
	_ = os.Setenv("TZ", "America/Central Time")

	err := configurator.NewConfiguration(&cfg, "DB")
	if err != nil {
		log.Fatalln(err)
	}

	configurator := configurator.NewConfigurator()

	configurator.New("postgres", &cfg, "db")

	db := db.New(configurator)

	return FixturesManager{
		db:       db,
		fixtures: fixturesList,
	}
}

func (f FixturesManager) ExecQuery(query string) error {
	_, err := f.db.Exec(context.TODO(), query)

	return err
}

func (f FixturesManager) SelectQuery(query string) (string, error) {
	var result string

	query = "select row_to_json(q)::text from (" + query + ") as q"

	err := f.db.QueryRow(context.TODO(), query).Scan(&result)

	return result, err
}

func (f FixturesManager) CleanAndApply(fixturePath string) error {
	for table, _ := range f.fixtures {
		err := f.Clean(table)
		if err != nil {
			return err
		}

		if err = f.LoadFixtures(); err != nil {
			return err
		}
	}
	return nil
}

func (f FixturesManager) Clean(table string) error {
	query := fmt.Sprintf("TRUNCATE TABLE %s", table)
	_, err := f.db.Exec(context.TODO(), query)
	if err != nil {
		return fmt.Errorf("failed delete fixtures: %w", err)
	}

	return err
}
