package adbclient

import (
	"io"
	"net"
	"net/http"
	"net/url"
	"os"

	"github.com/astaxie/beego/logs"
	adb "github.com/daimall/goadb"
	"github.com/daimall/goadb/wire"
)

type Options struct {
	Adbsrvport   int    // adbsrvport adb server 的端口号 [5037]
	ServicePort  int    // servicePort 设备侧监听端口号 [9008]
	Method       string // method 请求方法
	SaveFilePath string // 针对文件下载方法有效，保存文件到本地路径
}

type Option func(*Options)

func WithAdbSrvPort(value int) Option {
	return func(o *Options) {
		o.Adbsrvport = value
	}
}

func WithServicePort(value int) Option {
	return func(o *Options) {
		o.ServicePort = value
	}
}

func WithHTTPMethod(value string) Option {
	return func(o *Options) {
		o.Method = value
	}
}

func WithSaveFilePath(value string) Option {
	return func(o *Options) {
		o.SaveFilePath = value
	}
}

// 通过adb 发送 http请求（指定端口）
func HTTPRequest(path2adb, uri, deviceId string, options ...Option) (body []byte, err error) {
	// 设置默认参数
	opt := &Options{
		Adbsrvport:  5037,
		ServicePort: 9008,
		Method:      http.MethodGet,
	}
	// 应用自定义选项
	for _, option := range options {
		option(opt)
	}

	// 重写TCP Dialer，为了拿出底层的通信连接供HTTP协议栈使用
	var td = &tcpDialer{}
	var client *adb.Adb
	if client, err = adb.NewWithConfig(adb.ServerConfig{
		Port:      opt.Adbsrvport,
		PathToAdb: path2adb,
		Dialer:    td,
	}); err != nil {
		logs.Error("Failed to initialize adb client:", err.Error())
		return
	}
	if err = client.StartServer(); err != nil {
		logs.Error("start adb server failed,", err.Error())
		return nil, err
	}
	device := client.Device(adb.DeviceWithSerial(deviceId))
	var conn *wire.Conn
	if conn, err = device.ConnDeviceTCP(opt.ServicePort); err != nil {
		logs.Error("device change tcp(http) port failed,", err.Error())
		return
	}
	defer conn.Close()
	// 创建自定义的 Transport，指定底层的连接为已存在的 ADB TCP 连接
	transport := &http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			return td.NetConn, nil
		},
	}
	httpclient := &http.Client{
		Transport: transport,
	}
	var request *http.Request = &http.Request{
		Method: http.MethodGet,
		URL:    &url.URL{Scheme: "http", Host: "localhost", Path: uri},
	}
	// 发送请求并获取响应
	var response *http.Response
	if response, err = httpclient.Do(request); err != nil {
		logs.Error("failed to send http request,", err.Error())
		return nil, err
	}
	defer response.Body.Close()
	if opt.SaveFilePath != "" {
		// 保存文件的方式
		var file *os.File
		if file, err = os.Create(opt.SaveFilePath); err != nil {
			logs.Error("create save file failed,", err.Error())
			return nil, err
		}
		defer file.Close()
		_, err = io.Copy(file, response.Body)
		return nil, nil
	}
	// 非保存文件，获取请求结果
	body, err = io.ReadAll(response.Body)
	return body, err
}
