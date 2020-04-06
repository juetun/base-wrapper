package base

type ServiceBase struct {
	Context *Context
}

func (r *ServiceBase) SetContext(context []*Context) (s *ServiceBase) {
	switch len(context) {
	case 0:
		r.Context = NewContext()
		break
	case 1:
		r.Context = context[0]
		r.Context.Init() // 初始化一些没有初始化的对象
		break
	default:
		panic("你传递的参数当前不支持")
	}

	return r
}
