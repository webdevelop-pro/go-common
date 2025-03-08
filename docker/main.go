//nolint:paralleltest,thelper
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/webdevelop-pro/go-common/context/keys"
	"github.com/webdevelop-pro/go-common/db"
	"github.com/webdevelop-pro/go-common/logger"
	"github.com/webdevelop-pro/go-common/queue/pclient"
)

func CreateTestContext(ctx context.Context) context.Context {
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
		HTTPRequest: &logger.HTTPRequestContext{
			Method:    "POST",
			RemoteIP:  "0.0.0.0",
			URL:       "https://test",
			UserAgent: "test-agent",
			Referrer:  "test-Referrer",
		},
	}

	return keys.SetCtxValue(ctx, keys.LogInfo, logInfo)
}

func main() {
	ctx := context.TODO()
	ctx = CreateTestContext(ctx)
	resultInt, resultStr, resultTime := 0, "", time.Time{}

	os.Setenv("LOG_LEVEL", "info")
	os.Setenv("DB_LOG_LEVEL", "info")

	ddb := db.New(ctx)
	ddb.QueryRow(ctx, "select 1, 'b', now();").Scan(&resultInt, &resultStr, &resultTime)

	pubsubClient, err := pclient.New(ctx)
	fmt.Println(pubsubClient, err)
}
