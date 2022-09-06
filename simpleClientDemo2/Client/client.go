package Client

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/transport/http"
	"github.com/hashicorp/consul/api"

	"io"
	"math/rand"
	"net/url"
	"os"
	"strings"
)

// Direct: 直接调用服务端
// method:方法 fullUrl: 完整的url http://localhost:8000
// enc: http.EncodeRequestFunc dec: http.DecodeResponseFunc 这两个函数具体等一下会在Transport中进行详细解释
// requestStruct: 根据EndPoint定义的request结构体传参
func Direct(method string, fullUrl string, enc http.EncodeRequestFunc, dec http.DecodeResponseFunc, requestStruct interface{}) (interface{}, error) {
	// 1.解析url
	target, err := url.Parse(fullUrl)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// kit调用服务端拿到Client对象
	client := http.NewClient(strings.ToUpper(method), target, enc, dec)
	// 调用服务 client.Endpoint()返回一个可执行函数 传入context 和 请求数据结构体
	return client.Endpoint()(context.Background(), requestStruct)
}

// ServiceDiscovery: 通过服务发现的形式调用服务
// registryAddress: 注册中心的地址
// servicesName: 注册的服务名称
// tags: 可用标签
// passingOnly: true 只返回通过健康监测的实例
// method:方法
// enc: http.EncodeRequestFunc dec: http.DecodeResponseFunc 这两个函数具体等一下会在Transport中进行详细解释
// requestStruct: 根据EndPoint定义的request结构体传参
func ServiceDiscovery(method string, registryAddress string, enc http.EncodeRequestFunc, dec http.DecodeResponseFunc, requestStruct interface{}, servicesName string, passingOnly bool, tags ...string) (interface{}, error) {
	// 1.通过consul api创建一个client, 使用go-kit sd 封装一个专门的client 用于获取服务对象
	config := api.DefaultConfig()
	// registryAddress 注册中心的地址
	config.Address = registryAddress

	// 这里是拿到 consul 中的client
	apiClient, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	// kit封装client
	client := consul.NewClient(apiClient)

	logger := log.NewLogfmtLogger(os.Stdout)

	// 创建实例
	// logger 可以使用自己的logger对象 例如zap等
	// instances: 实例对象
	instances := consul.NewInstancer(client, logger, servicesName, tags, passingOnly)

	// f: Factory
	// servicesUrl: 传入服务的url, 通过这个url根据直连一样去获取endpoint.Endpoint对象
	f := func(servicesUrl string) (endpoint.Endpoint, io.Closer, error) {
		fmt.Println(servicesUrl)
		// 解析url
		target, err := url.Parse("http://" + servicesUrl)
		if err != nil {
			return nil, nil, err
		}
		return http.NewClient(strings.ToUpper(method), target, enc, dec).Endpoint(), nil, nil
	}

	// 获取endpoint可执行对象 与我们直连client.Endpoint()返回的一样
	// 传入 instances: 实例对象  Factory:工厂模式  logger: 日志对象
	endpointer := sd.NewEndpointer(instances, f, logger)

	// 获取所有实例 endpoints
	endpoints, err := endpointer.Endpoints()
	if err != nil {
		return nil, err
	}
	// 随机选择一个执行
	l := len(endpoints)
	if l == 0 {
		return nil, errors.New("len(endpoints) == 0")
	}
	return endpoints[rand.Intn(len(endpoints))](context.Background(), requestStruct)
}
