package tests

import "context"

type FixturesManager interface {
	CleanAndApply() error
	SetCTX(context.Context) context.Context
}
