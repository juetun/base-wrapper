/**
* @Author:changjiang
* @Description:
* @File:dao
* @Version: 1.0.0
* @Date 2020/4/5 8:22 下午
 */
package base

type ServiceDao struct {
	Context *Context
}

func (r *ServiceDao) SetContext(context []*Context) (s *ServiceDao) {
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