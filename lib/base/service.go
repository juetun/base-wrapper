package base

type (
	ServiceBase struct {
		Context *Context
	}
)


func (r *ServiceBase) SetContext(context ...*Context) (s *ServiceBase) {
	for _, cont := range context {
		cont.InitContext()
	}
	switch len(context) {
	case 0:
		r.Context = NewContext()
		break
	case 1:
		r.Context = context[0]
		break
	default:
		panic("你传递的参数当前不支持")
	}

	return r
}

// RecordLog 记录日志使用
func (r *ServiceBase) RecordLog(locKey string, logContent map[string]interface{}, err error, needRecordInfo ...bool) {
	if err != nil {
		logContent["err"] = err.Error()
		r.Context.Error(logContent, locKey)
		return
	}
	if len(needRecordInfo) > 0 {
		r.Context.Info(logContent, locKey)
	}
}
