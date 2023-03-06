package sentryx

import (
	"github.com/getsentry/sentry-go"
	"pinnacle-primary-be/core/logx"
)

func NewSentry(r Config) {
	if !r.Available() {
		return
	}
	if err := sentry.Init(sentry.ClientOptions{
		Dsn: r.Dsn,
	}); err != nil {
		logx.Error(err)
	}
}
