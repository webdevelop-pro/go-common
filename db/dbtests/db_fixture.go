package dbtests

import (
	"context"
	"fmt"
	"os"

	"github.com/webdevelop-pro/go-common/configurator"
	"github.com/webdevelop-pro/go-common/db"
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
	db       *db.DB
	cfg      db.Config
	fixtures []Fixture
}

func NewFixturesManager(ctx context.Context, fixtures ...Fixture) FixturesManager {
	cfg := db.Config{}

	// Fix for timezones
	_ = os.Setenv("TZ", "UTC")

	err := configurator.NewConfiguration(&cfg, "DB")
	if err != nil {
		panic(err)
	}

	db := db.New(ctx)
	return FixturesManager{
		db:       db,
		cfg:      cfg,
		fixtures: fixtures,
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

func (f FixturesManager) CleanAndApply(fixtures []Fixture) error {
	for _, fixture := range fixtures {
		err := f.Clean(fixture.table)
		if err != nil {
			return err
		}
	}
	return f.LoadFixtures(fixtures)
}

func (f FixturesManager) GetCTX() context.Context {
	return context.WithValue(context.Background(), "db", f.db)
}

func (f FixturesManager) Clean(table string) error {
	query := fmt.Sprintf(
		"DELETE FROM %s; select setval('%s_id_seq',(select max(id)+1 from %s));",
		table, table, table,
	)

	_, err := f.db.Exec(context.TODO(), query)
	if err != nil {
		return fmt.Errorf("failed delete fixtures: %w", err)
	}

	return err
}