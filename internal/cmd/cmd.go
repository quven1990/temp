package cmd

import (
	"context"

	alarmhisapi "gf_api/internal/controller/alarm_his_api"
	api "gf_api/internal/controller/api"
	childsysdataapi "gf_api/internal/controller/client3.0_api/child_sys_data_api"
	childsysnumber "gf_api/internal/controller/client3.0_api/child_sys_number_api"
	configapi "gf_api/internal/controller/config_api"
	controlsysapi "gf_api/internal/controller/client3.0_api/control_sys_api"
	gethikdataapi "gf_api/internal/controller/client3.0_api/get_hik_data_api"
	getstationfrqprogramapi "gf_api/internal/controller/client3.0_api/get_station_frq_program_api"
	getstationnoteapi "gf_api/internal/controller/client3.0_api/get_station_note_api"
	getsyslogapi "gf_api/internal/controller/client3.0_api/get_sys_log_api"
	gettimeapi "gf_api/internal/controller/client3.0_api/get_time_api"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
)

// 注册路由api
var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "启动 HTTP 服务",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			// 注册路由组
			s.Group("/api", func(group *ghttp.RouterGroup) {
				//鉴权中间件，为/api所有路由添加。测试阶段先不用
				//group.Middleware(ghttp.MiddlewareHandlerResponse, middleware.AuthMiddleware)

				// Basic 相关接口
				gettimeapi.Register(group)
				childsysnumber.Register(group)
				api.Register(group)
				childsysdataapi.Register(group)
				getstationfrqprogramapi.Register(group)

				// Resource 相关接口
				alarmhisapi.Register(group)
				getstationnoteapi.Register(group)
				gethikdataapi.Register(group)
				getsyslogapi.Register(group)
				controlsysapi.Register(group)

				// Config 相关接口（配置管理）
				configapi.Register(group)

				// 转发到配置服务（已注释，如需使用请取消注释）
				// group.Group("/config", func(g *ghttp.RouterGroup) {
				// 	g.ALL("/*any", proxy.Proxy("http://config-service"))
				// })

				// 转发到日志服务（已注释，如需使用请取消注释）
				// group.Group("/logs", func(g *ghttp.RouterGroup) {
				// 	//本地开发测试电脑ip是10.170.0.96，谁开发改为谁的。
				// 	// 访问代理的日志服务接口http://10.170.0.96:8001/api/logs/log-manage会跳到http://localhost:80/api/log/log-manage
				// 	//但是跳转前会先判断是否有鉴权，否则跳转到鉴权平台。
				// 	g.ALL("/log-manage", proxy.Proxy("http://localhost:80/api/log/log-manage"))
				// })
			})

			// 启动服务
			s.Run()
			return nil
		},
	}
)
