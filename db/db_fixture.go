package db

import (
	"log"
	"os"

	"github.com/webdevelop-pro/go-common/configurator"
)

type Fixture struct {
	table    string
	filePath string
}

func NewFixture(table, filePath string) Fixture {
	return Fixture{
		table:    table,
		filePath: filePath,
	}
}

type FixturesManager struct {
	db  *DB
	cfg Config
}

func NewFixturesManager() FixturesManager {
	cfg := Config{}

	// Fix for timezones
	_ = os.Setenv("TZ", "America/Central Time")

	err := configurator.NewConfiguration(&cfg, "DB")
	if err != nil {
		log.Fatalln(err)
	}

	configurator := configurator.NewConfigurator()

	configurator.New("postgres", &cfg, "db")

	db := New()

	return FixturesManager{
		db:  db,
		cfg: cfg,
	}
}

func (f FixturesManager) ExecQuery(query string) error {
	// ToDo: FixMe
	/*
		_, err := f.db.Exec(context.TODO(), query)

		return err
	*/
	return nil
}

func (f FixturesManager) SelectQuery(query string) (string, error) {
	var result string

	// ToDo: FixMe
	/*
		query = "select row_to_json(q)::text from (" + query + ") as q"

		err := f.db.QueryRow(context.TODO(), query).Scan(&result)

		return result, err
	*/
	return result, nil
}

func (f FixturesManager) CleanAndApply(fixtures []Fixture) error {
	for _, fixture := range fixtures {
		err := f.Clean(fixture.table)
		if err != nil {
			return err
		}
	}
	// ToDo: Fix me
	return nil
	// return f.LoadFixtures(fixtures)
}

func (f FixturesManager) Clean(table string) error {
	// ToDo: FixMe
	/*
		query := fmt.Sprintf(
			"DELETE FROM %s; select setval('%s_id_seq',(select max(id)+1 from %s));",
			table, table, table,
		)

			_, err := f.db.Exec(context.TODO(), query)
			if err != nil {
				return fmt.Errorf("failed delete fixtures: %w", err)
			}

			return err
	*/
	return nil
}
