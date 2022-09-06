package Endpoint

// endpoint.go 定义 Request、Response 格式, 并且可以使用闭包来实现各种中间件的嵌套
// 这里了解 protobuf 的比较好理解点
// 就是声明 接收数据和响应数据的结构体 并通过构造函数创建 在创建的过程当然可以使用闭包来进行一些你想要的操作啦

// 这里根据我们Demo来创建一个响应和请求
// 当然你想怎么创建怎么创建 也可以共用 这里我分开写 便于大家看的清楚

// Hello 业务使用的请求和响应格式
// HelloRequest 请求格式
type HelloRequest struct {
	Name string `json:"name"`
}

// HelloResponse 响应格式
type HelloResponse struct {
	Reply string `json:"reply"`
}

// Bye 业务使用的请求和响应格式
// ByeRequest 请求格式
type ByeRequest struct {
	Name string `json:"name"`
}

// ByeResponse 响应格式
type ByeResponse struct {
	Reply string `json:"reply"`
}

// ------------ 当然 也可以通用的写 ----------
// Request 请求格式
type Request struct {
	Name string `json:"name"`
}

// Response 响应格式
type Response struct {
	Reply string `json:"reply"`
}
