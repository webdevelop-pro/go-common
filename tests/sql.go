package tests

import (
	"fmt"
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
		assert.FailNow(t, err.Error()+fmt.Sprintf(" failed sql: %s", testCase.ExpectedSelectQuery))
	}

	CompareJsonBody(t, []byte(actualResult), []byte(testCase.ExpectedResult))
}
