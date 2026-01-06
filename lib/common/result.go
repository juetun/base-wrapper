package common

type HttpResult struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}
type Object struct {
}

func NewHttpResult() *HttpResult {
	return &HttpResult{
		Code:    0,
		Data:    new(Object),
		Message: "ok",
	}
}
func (r *HttpResult) SetCode(code int) *HttpResult {
	r.Code = code
	return r
}
func (r *HttpResult) SetData(data interface{}) *HttpResult {
	r.Data = data
	return r
}
func (r *HttpResult) SetMessage(message string) *HttpResult {
	r.Message = message
	return r
}
