package gw

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func IncomeMatcher(key string) (string, bool) {
	switch key {
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}

func OutgoingMatcher(key string) (string, bool) {
	switch key {
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}
