package base

const (
	ErrorParameterCode = 1 // 参数错误
)
const (
	SuccessCode            = 0            // 操作成功返回状态码
	ErrorSystem            = iota + 11000 // 系统错误
	ErrorBUSSSINESS                       // 业务系统错误
	ErrorSqlCode                          // 数据库错误信息
	ErrorUpgradeCode                      // 强制用户升级App状态
	ErrorRedisCode                        // Redis错误信息
	ErrorElasticSearchCode                // ElasticSearch错误信息
	ErrorOssCode                          // Oss错误信息
	ErrorNotLogin          = 401          //未登录
)

type ErrorRuntimeStruct struct {
	Code int   `json:"code"`
	err  error `json:"err"`
}

// NewErrorRuntime SQL错误信息
func NewErrorRuntime(err error, code ...int) (res error) {
	cd := ErrorSqlCode
	if len(code) > 0 {
		cd = code[0]
	}
	res = &ErrorRuntimeStruct{
		Code: cd,
		err:  err,
	}
	return
}

func (r *ErrorRuntimeStruct) Error() (res string) {
	res = r.err.Error()
	return
}
