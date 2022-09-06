// main.go
package main

import (
	httpTransport "github.com/go-kit/kit/transport/http"
	"microservice/simpleServerDemo/Endpoint"
	"microservice/simpleServerDemo/Server"
	"microservice/simpleServerDemo/Transport"
	"net/http"
)

// 服务发布

func main() {
	// 1.先创建我们最开始定义的Server/server.go
	s := Server.Server{}

	// 2.在用EndPoint/endpoint.go 创建业务服务
	hello := Endpoint.MakeServerEndPointHello(s)
	Bye := Endpoint.MakeServerEndPointBye(s)

	// 3.使用 kit 创建 handler
	// 固定格式
	// 传入 业务服务 以及 定义的 加密解密方法
	helloServer := httpTransport.NewServer(hello, Transport.HelloDecodeRequest, Transport.HelloEncodeResponse)
	sayServer := httpTransport.NewServer(Bye, Transport.ByeDecodeRequest, Transport.ByeEncodeResponse)

	// 使用http包启动服务
	go http.ListenAndServe("0.0.0.0:8000", helloServer)

	go http.ListenAndServe("0.0.0.0:8001", sayServer)
	select {}
}
