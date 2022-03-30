package main

import (
	"jhr.com/apirelay/global"
	"jhr.com/apirelay/initialize"
	"jhr.com/apirelay/web"
)

func main() {
	initialize.InitLogger()
	// 初始化配置文件
	initialize.InitConfig()
	// 根据配置动态注册路由
	r := web.RegisterRouter()
	// 测试服务
	if global.ServerConfig.Test {
		web.TestServer(r)
		return
	}
	// 启动服务端
	web.RunServer(r)
}
