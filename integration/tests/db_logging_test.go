//nolint:paralleltest,thelper
package tests

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/webdevelop-pro/go-common/tests"
	"github.com/webdevelop-pro/go-logger"

	// . "github.com/webdevelop-pro/go-common/tests"

	dbDriver "github.com/webdevelop-pro/go-common/db"
)

func CreateTestContext() context.Context {
	logInfo := logger.ServiceContext{
		Service: "test-db",
		Version: "v1.0.0",
		SourceReference: &logger.SourceReference{
			Repository: "https://github",
			RevisionID: "1111",
		},
		User:      "test-test-test",
		RequestID: "req-1",
		MSGID:     "msg-1",
		HttpRequest: &logger.HttpRequestContext{
			Method:    "POST",
			RemoteIp:  "0.0.0.0",
			URL:       "https://test",
			UserAgent: "test-agent",
			Referrer:  "test-Referrer",
		},
	}

	return context.WithValue(context.Background(), logger.ServiceContextInfo, logInfo)
}

func TestLogger_DBQuery(t *testing.T) {
	ctx := CreateTestContext()
	resultInt, resultStr, resultTime := 0, "", time.Time{}

	os.Setenv("LOG_LEVEL", "info")
	os.Setenv("DB_LOG_LEVEL", "info")

	stdout := ConnectToStdout()

	db := dbDriver.New()
	err := db.QueryRow(ctx, "select 1, 'b', now();").Scan(&resultInt, &resultStr, &resultTime)
	assert.Nil(t, err)

	actualLogs := ReadStdout(stdout)
	expected := `
		{
			"level": "info",
			"component": "db",
			"data": {
				"args": [],
				"commandTag": "SELECT 1",
				"pid": "%any%",
				"sql": "select 1, 'b', now();",
				"time": "%any%"
			},
			"severity": "INFO",
			"serviceContext": {
				"service": "test-db",
				"version": "v1.0.0",
				"user": "test-test-test",
				"request_id": "req-1",
				"msg_id": "msg-1",
				"httpRequest": {
					"method": "POST",
					"url": "https://test",
					"userAgent": "test-agent",
					"referrer": "test-Referrer",
					"responseStatusCode": 0,
					"remoteIp": "0.0.0.0"
				},
				"sourceReference": {
					"repository": "https://github",
					"revisionId": "1111"
				}
			},
			"time": "%any%",
			"message": "Query"
		}
	`
	actual := actualLogs[len(actualLogs)-1]

	tests.CompareJSONBody(t, []byte(actual), []byte(expected))
}

func TestLogger_DBExec(t *testing.T) {
	ctx := CreateTestContext()

	os.Setenv("LOG_LEVEL", "info")
	os.Setenv("DB_LOG_LEVEL", "info")

	stdout := ConnectToStdout()

	db := dbDriver.New()
	_, err := db.Exec(ctx, "update user_users set email='' where id=-1;")
	assert.Nil(t, err)

	actualLogs := ReadStdout(stdout)
	expected := `
		{
			"level": "info",
			"component": "db",
			"data": {
				"args": [],
				"commandTag": "UPDATE 0",
				"pid": "%any%",
				"sql": "update user_users set email='' where id=-1;",
				"time": "%any%"
			},
			"severity": "INFO",
			"serviceContext": {
				"service": "test-db",
				"version": "v1.0.0",
				"user": "test-test-test",
				"request_id": "req-1",
				"msg_id": "msg-1",
				"httpRequest": {
					"method": "POST",
					"url": "https://test",
					"userAgent": "test-agent",
					"referrer": "test-Referrer",
					"responseStatusCode": 0,
					"remoteIp": "0.0.0.0"
				},
				"sourceReference": {
					"repository": "https://github",
					"revisionId": "1111"
				}
			},
			"time": "%any%",
			"message": "Query"
		}
	`
	actual := actualLogs[len(actualLogs)-1]

	tests.CompareJSONBody(t, []byte(actual), []byte(expected))
}

func TestLogger_DBQuery_ERROR(t *testing.T) {
	ctx := CreateTestContext()
	resultInt, resultStr, resultTime := 0, "", time.Time{}

	os.Setenv("LOG_LEVEL", "info")
	os.Setenv("DB_LOG_LEVEL", "info")

	stdout := ConnectToStdout()

	db := dbDriver.New()
	err := db.QueryRow(ctx, "select asd;").Scan(&resultInt, &resultStr, &resultTime)
	assert.NotNil(t, err)

	actualLogs := ReadStdout(stdout)
	expected := `
		{
			"level": "error",
			"component": "db",
			"error": "ERROR: column \"asd\" does not exist (SQLSTATE 42703)",
			"data": {
				"args": [],
				"err": {
					"Severity": "ERROR",
					"Code": "42703",
					"Message": "column \"asd\" does not exist",
					"Detail": "",
					"Hint": "",
					"Position": 8,
					"InternalPosition": 0,
					"InternalQuery": "",
					"Where": "",
					"SchemaName": "",
					"TableName": "",
					"ColumnName": "",
					"DataTypeName": "",
					"ConstraintName": "",
					"File": "parse_relation.c",
					"Line": 3633,
					"Routine": "errorMissingColumn"
				},
				"pid": "%any%",
				"sql": "select asd;",
				"time": "%any%"
			},
			"severity": "ERROR",
			"serviceContext": {
				"service": "test-db",
				"version": "v1.0.0",
				"user": "test-test-test",
				"request_id": "req-1",
				"msg_id": "msg-1",
				"httpRequest": {
					"method": "POST",
					"url": "https://test",
					"userAgent": "test-agent",
					"referrer": "test-Referrer",
					"responseStatusCode": 0,
					"remoteIp": "0.0.0.0"
				},
				"sourceReference": {
					"repository": "https://github",
					"revisionId": "1111"
				}
			},
			"time": "%any%",
			"message": "column \"asd\" does not exist"
		}
	`
	actual := actualLogs[len(actualLogs)-1]

	tests.CompareJSONBody(t, []byte(actual), []byte(expected))
}
