//nolint:gochecknoglobals
package main

import (
	"context"

	"github.com/pkg/errors"
	"github.com/webdevelop-pro/go-common/logger"
	"github.com/webdevelop-pro/go-common/verser"
)

var ErrNotFound = errors.New("row not found")

func FindUser() error {
	return ErrNotFound
}

var (
	service    string
	version    string
	repository string
	revisionID string
)

type Context struct{}

func (ctx Context) Get(key string) interface{} {
	if key == "user" {
		return "user-uuid"
	}
	return ""
}

func main() {
	verser.SetServiVersRepoRevis(service, version, repository, revisionID)
	ctx := context.Background()
	defaultLogger := logger.NewComponentLogger(ctx, "main-cli")
	err := FindUser()
	defaultLogger.Error().Stack().Err(err).Msg(ErrNotFound.Error())
}
