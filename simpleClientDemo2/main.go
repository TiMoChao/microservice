// main.go
package main

import (
	"fmt"
	"microservice/simpleClientDemo2/Client"
	"microservice/simpleClientDemo2/Endpoint"
	"microservice/simpleClientDemo2/Transport"
)

// 服务发布

// 调用我们在client封装的函数就好了
func main() {
	//i, err := Client.Direct("GET", "http://127.0.0.1:8000", Transport.HelloEncodeRequestFunc, Transport.HelloDecodeResponseFunc, Endpoint.HelloRequest{Name: "songzhibin"})
	i, err := Client.ServiceDiscovery("GET", "http://127.0.0.1:8500", Transport.ByeEncodeRequestFunc, Transport.ByeDecodeResponseFunc, Endpoint.HelloRequest{Name: "songzhibin"}, "测试", true, "test")
	if err != nil {
		fmt.Println(err)
		return
	}
	res, ok := i.(Endpoint.HelloResponse)
	if !ok {
		fmt.Println("no ok")
		return
	}
	fmt.Println(res)
}
