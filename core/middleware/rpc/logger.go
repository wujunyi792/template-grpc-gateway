package rpc

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"pinnacle-primary-be/core/errorx"
	"pinnacle-primary-be/core/logx"
)

func LoggerInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	resp, err = handler(ctx, req)
	if err != nil {
		causeErr := errors.Cause(err)                  // err类型
		if e, ok := causeErr.(*errorx.CodeError); ok { //自定义错误类型
			logx.WithContext(ctx).Errorf("【RPC-SRV-ERR】 %+v", err)

			//转成grpc err
			err = status.Error(codes.Code(e.GetErrCode()), e.GetErrMsg())
		} else {
			logx.WithContext(ctx).Errorf("【RPC-SRV-ERR】 %+v", err)
		}

	}

	return resp, err
}
