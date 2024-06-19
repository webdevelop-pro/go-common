package tests

import "context"

type FixturesManager interface {
	CleanAndApply() error
	GetCTX() context.Context
}
