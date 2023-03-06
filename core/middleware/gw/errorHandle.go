package gw

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"pinnacle-primary-be/core/jsonx"
	"pinnacle-primary-be/core/middleware/response"
)

func GrpcGatewayError(_ context.Context, _ *runtime.ServeMux, _ runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	s, ok := status.FromError(err)
	if !ok {
		s = status.New(codes.Internal, err.Error())
	}

	// TODO 可用性测试
	httpError := response.JsonResponse{Code: int32(s.Code()), Message: s.Message(), Error: s.Details()}

	resp, _ := jsonx.Marshal(httpError)
	w.Header().Set("Content-type", "application/json")

	code := int(s.Code()) / 1000
	if isGrpcErr(s.Code()) {
		code = runtime.HTTPStatusFromCode(s.Code())
	} else {
		if code == 0 {
			code = 500
		}
	}

	w.WriteHeader(code)
	_, _ = w.Write(resp)
}

func isGrpcErr(code codes.Code) bool {
	return code < 16
}
