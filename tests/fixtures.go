package tests

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

type FixturesManager struct{}

func (fm *FixturesManager) CleanAndApply(f []Fixture) error {
	return nil
}

func NewFixturesManager() FixturesManager {
	return FixturesManager{}
}
