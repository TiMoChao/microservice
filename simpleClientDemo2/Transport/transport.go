package Transport

import (
	"context"
	"encoding/json"
	"errors"
	"microservice/simpleClientDemo2/Endpoint"
	"net/http"
	"net/url"
	"strconv"
)

// Transport/transport.go 主要负责HTTP、gRpc、thrift等相关的逻辑

// 这里有两个关键函数
// DecodeRequest & EncodeResponse 函数签名是固定的哟
// func EncodeRequestFunc (context.Context, *http.Request, interface{}) error
// func DecodeResponseFunc (context.Context, *http.Response) (response interface{}, err error)

// HelloEncodeRequestFunc: 处理请求数据符合服务方要求的数据
func HelloEncodeRequestFunc(c context.Context, request *http.Request, r interface{}) error {
	// r就是我们在EndPoint中定义的请求响应对象
	req, ok := r.(Endpoint.HelloRequest)
	if !ok {
		return errors.New("断言失败")
	}
	// 拿到自定义的请求对象对url做业务处理
	request.URL.Path += "/hello"
	data := url.Values{}
	data.Set("name", req.Name)
	request.URL.RawQuery = data.Encode()
	// 实际上这里做的就是增加url参数 body之类的一些事情,简而言之就是构建http请求需要的一些资源
	return nil
}

// HelloDecodeResponseFunc: 解密服务方传回的数据
func HelloDecodeResponseFunc(c context.Context, res *http.Response) (response interface{}, err error) {
	// 判断响应
	if res.StatusCode != 200 {
		return nil, errors.New("异常的响应码" + strconv.Itoa(res.StatusCode))
	}
	// body中的内容需要我们解析成我们通用定义好的内容
	var r Endpoint.HelloResponse
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// ByeEncodeRequestFunc: 处理请求数据符合服务方要求的数据
func ByeEncodeRequestFunc(c context.Context, request *http.Request, r interface{}) error {
	// r就是我们在EndPoint中定义的请求响应对象
	req, ok := r.(Endpoint.HelloRequest)
	if !ok {
		return errors.New("断言失败")
	}
	// 拿到自定义的请求对象对url做业务处理
	request.URL.Path += "/bye"
	data := url.Values{}
	data.Set("name", req.Name)
	request.URL.RawQuery = data.Encode()
	// 实际上这里做的就是增加url参数 body之类的一些事情,简而言之就是构建http请求需要的一些资源
	return nil
}

// ByeDecodeResponseFunc: 解密服务方传回的数据
func ByeDecodeResponseFunc(c context.Context, res *http.Response) (response interface{}, err error) {
	// 判断响应
	if res.StatusCode != 200 {
		return nil, errors.New("异常的响应码" + strconv.Itoa(res.StatusCode))
	}
	// body中的内容需要我们解析成我们通用定义好的内容
	var r Endpoint.HelloResponse
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return r, nil
}
