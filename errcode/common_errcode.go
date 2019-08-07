package errcode

// All common ecode
var (
	OK = add(0) // 正确

	Unauthorized = add(-401) // 未认证
	Canceled     = add(-498) // 客户端取消请求
	ServerErr    = add(-500) // 服务器错误
	Deadline     = add(-504) // 服务调用超时
)
