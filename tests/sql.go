package tests

import (
	"fmt"
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
		query := strings.Replace(testCase.ExpectedSelectQuery, "\t", " ", -1)
		query = strings.Replace(query, "\n", " ", -1)
		query = strings.Replace(query, "  ", " ", -1)
		query = strings.Replace(query, "  ", " ", -1)
		assert.FailNow(t, err.Error()+fmt.Sprintf(" failed sql: %s", query))
	}

	CompareJsonBody(t, []byte(actualResult), []byte(testCase.ExpectedResult))
}
