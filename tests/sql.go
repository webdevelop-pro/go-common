package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type SQLTestCase struct {
	ExpectedSelectQuery string
	ExpectedResult      string
}

func AssertSQL(t *testing.T, fManager FixturesManager, testCase SQLTestCase) {
	t.Helper()

	actualResult, err := fManager.SelectQuery(testCase.ExpectedSelectQuery)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	CompareJsonBody(t, []byte(actualResult), []byte(testCase.ExpectedResult))
}
