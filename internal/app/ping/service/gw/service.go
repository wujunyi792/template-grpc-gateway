package gw

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/emptypb"
	"math/rand"
	"pinnacle-primary-be/core/errorx"
	"pinnacle-primary-be/core/tracex"
	pingV1 "pinnacle-primary-be/gen/proto/ping/v1"
	"pinnacle-primary-be/internal/app/ping/dao"
	"pinnacle-primary-be/internal/app/ping/model"
	"sync"
	"time"
)

type S struct {
	pingV1.UnimplementedPingServiceServer
}

func (s S) Ping(ctx context.Context, in *pingV1.PingRequest) (*pingV1.PingResponse, error) {
	ctx, span := tracex.TracerFromContext(ctx).Start(ctx, "Ping")
	defer span.End()
	wg := sync.WaitGroup{}
	wg.Add(4)
	for i := 0; i < 4; i++ {
		go func(index int) {
			defer wg.Done()
			_, span := tracex.TracerFromContext(ctx).Start(ctx, fmt.Sprintf("Ping-%d", index+1),
				trace.WithAttributes(attribute.Int("key", index+1)))
			defer span.End()
			time.Sleep(time.Duration(rand.Intn(120-20)+20) * time.Millisecond)
		}(i)
	}
	time.Sleep(23 * time.Millisecond)
	for i := 0; i < 10; i++ {
		dao.Ping.GetOrm().WithContext(ctx).Create(&model.Ping{})
	}
	_err := dao.Ping.GetOrm().WithContext(ctx).Where("uid = ?", 1).First(&model.Ping{}).Error
	wg.Wait()
	if _err != nil {
		span.RecordError(_err)
		return nil, errorx.InternalError
	}
	return &pingV1.PingResponse{Value: in.Value}, nil
}

func (s S) PingErr(ctx context.Context, _ *pingV1.ABitOfEverything) (*emptypb.Empty, error) {
	_, span := tracex.TracerFromContext(ctx).Start(ctx, "PingErr")
	defer span.End()
	span.RecordError(errors.New("PingErr not implemented"))
	return nil, errors.Wrapf(errorx.NewErrCodeMsg(502000, "这是不可以的操作"), "PingErr not implemented")
}

func (s S) PingAuth(ctx context.Context, in *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, errorx.UnAuthorizedError
}
