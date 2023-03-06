package gw

import (
	"context"
	"encoding/json"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"net/http"
	"pinnacle-primary-be/core/jsonx"
	"pinnacle-primary-be/core/middleware/response"
)

func GrpcGatewaySuccess(ctx context.Context, w http.ResponseWriter, m proto.Message) error {
	marshaller := runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			EmitUnpopulated: true,
			UseProtoNames:   true,
		},
	}
	jsonBytes, _ := marshaller.Marshal(m)
	buf, err := jsonx.Marshal(response.JsonResponse{
		Message: "success",
		Data:    json.RawMessage(jsonBytes),
	})
	if err != nil {
		return err
	}
	_, _ = w.Write(buf)
	return nil
}

type CustomMarshaller struct {
	runtime.JSONPb
}

func (m *CustomMarshaller) Marshal(v interface{}) ([]byte, error) {
	return nil, nil
}
