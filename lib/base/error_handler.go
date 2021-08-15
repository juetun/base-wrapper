package base

const (
	ErrorSqlCode       = iota + 11000 // 数据库错误信息
	ErrorRedisCode                    // Redis错误信息
	ErrorParameterCode                // 参数错误
)

type ErrorRuntime struct {
	Code int   `json:"code"`
	err  error `json:"err"`
}

// NewErrorRuntime SQL错误信息
func NewErrorRuntime(err error, code ...int) (res error) {
	cd := ErrorSqlCode
	if len(code) > 0 {
		cd = code[0]
	}
	res = &ErrorRuntime{
		Code: cd,
		err:  err,
	}
	return
}

func (r *ErrorRuntime) Error() (res string) {
	res = r.err.Error()
	return
}
