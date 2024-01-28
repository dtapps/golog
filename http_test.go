package golog

import (
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestHttp(t *testing.T) {
	url := "https://api.dtapp.net/ip"

	// 发送 GET 请求
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// 输出响应内容
	fmt.Println("Response Body:", string(body))

	// 输出请求信息
	fmt.Println("Request URL 请求的 URL 对象:", resp.Request.URL)
	fmt.Println("Request Method 请求的 HTTP 方法:", resp.Request.Method)
	fmt.Println("Request Headers 请求的头部信息:", resp.Request.Header)

	// 输出响应信息
	fmt.Println("Response Status 响应的状态行，包括状态码和描述:", resp.Status)
	fmt.Println("Response Status Code 响应的状态码:", resp.StatusCode)
	fmt.Println("Response Proto 响应的协议版本:", resp.Proto)
	fmt.Println("Response ProtoMajor 响应协议版本的主版本:", resp.ProtoMajor)
	fmt.Println("Response ProtoMinor 响应协议版本的次版本:", resp.ProtoMinor)
	fmt.Println("Response Headers: 响应的头部信息", resp.Header)
	fmt.Println("Response Body: 响应体的读取接口", resp.Body)
	fmt.Println("Response ContentLength: 响应体的长度（如果已知）", resp.ContentLength)
	fmt.Println("Response TransferEncoding: 传输编码方式", resp.TransferEncoding)
	fmt.Println("Response Close: 指示是否在读取完响应体后关闭连接", resp.Close)
}
