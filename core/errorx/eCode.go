package errorx

const OK uint32 = 200

const SERVER_COMMON_ERROR uint32 = 500000
const REUQEST_PARAM_ERROR uint32 = 400000
const TOKEN_EXPIRE_ERROR uint32 = 403000

var (
	UnAuthorizedError = NewErrCodeMsg(TOKEN_EXPIRE_ERROR, "登录过期失效，请重新登陆")
	InternalError     = NewErrCodeMsg(SERVER_COMMON_ERROR, "内部异常")
)
