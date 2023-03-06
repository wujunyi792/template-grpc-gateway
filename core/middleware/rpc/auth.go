package rpc

import (
	"context"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"pinnacle-primary-be/pkg/jwt"
)

func AuthInterceptor(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return ctx, nil
	}
	entry, err := jwt.ParseToken(token)
	if err == nil {
		return context.WithValue(ctx, "uid", entry.Info.Uid), nil
	}
	return ctx, nil
}
