package logger

import "github.com/rs/zerolog"

type ContextHook struct{}

func (h ContextHook) Run(e *zerolog.Event, _ zerolog.Level, _ string) {
	/*
		ctx := e.GetCtx()

		serviceCtx, _ := keys.GetCtxValue(ctx, keys.LogInfo).(ServiceContext)

		e.Interface("serviceContext", serviceCtx)
	*/
}
