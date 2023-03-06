package e

import "pinnacle-primary-be/pkg/ecode"

var (
	ErrNotExist = ecode.New(40400)
)

var ECode = map[ecode.Code]string{
	ErrNotExist: "资源不存在",
}

func init() {
	ecode.Register(ECode)
}
