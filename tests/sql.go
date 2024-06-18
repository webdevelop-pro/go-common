package tests

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type SQLTestCase struct {
	ExpectedSelectQuery string
	ExpectedResult      string
}

// ToDo
// Do we still need it?
func AssertSQL(t *testing.T, fManager FixturesManager, testCase SQLTestCase) {
	t.Helper()

	actualResult, err := fManager.SelectQuery(testCase.ExpectedSelectQuery)
	if err != nil {
		query := strings.ReplaceAll(testCase.ExpectedSelectQuery, "\t", " ")
		query = strings.ReplaceAll(query, "\n", " ")
		query = strings.ReplaceAll(query, "  ", " ")
		query = strings.ReplaceAll(query, "  ", " ")
		assert.FailNow(t, err.Error()+" failed sql: "+query)
	}

	CompareJSONBody(t, []byte(actualResult), []byte(testCase.ExpectedResult))
}
