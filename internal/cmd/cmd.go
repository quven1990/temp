package cmd

import (
	"context"

	api "gf_api/internal/controller/api"
	childsysdataapi "gf_api/internal/controller/client3.0_api/child_sys_data_api"
	childsysnumber "gf_api/internal/controller/client3.0_api/child_sys_number_api"
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
				//roup.Middleware(ghttp.MiddlewareHandlerResponse, middleware.AuthMiddleware)

				// 转发到配置服务
				// group.Group("/config", func(g *ghttp.RouterGroup) {
				// 	g.ALL("/*any", proxy.Proxy("http://config-service"))
				// })

				// 转发到日志服务
				// group.Group("/logs", func(g *ghttp.RouterGroup) {
				// 	//本地开发测试电脑ip是10.170.0.96，谁开发改为谁的。
				// 	// 访问代理的日志服务接口http://10.170.0.96:8001/api/logs/log-manage会跳到http://localhost:80/api/log/log-manage
				// 	//但是跳转前会先判断是否有鉴权，否则跳转到鉴权平台。
				// 	g.ALL("/log-manage", proxy.Proxy("http://localhost:80/api/log/log-manage"))

				//获取服务器时间
				gettimeapi.Register(group)

				//获取子系统编号
				childsysnumber.Register(group)

				//获取台站信息、台站总览数据
				api.Register(group)

				//获取台站注意事项
				getstationnoteapi.Register(group)

				//获取台站的所有频率和节目名称
				getstationfrqprogramapi.Register(group)

				//获取台站海康威视数据
				gethikdataapi.Register(group)

				//获取用户操作日志信息
				getsyslogapi.Register(group)

				//获取子系统信息
				childsysdataapi.Register(group)

				//台站客户端的下发控制
				controlsysapi.Register(group)

				// })
			})

			// 启动服务
			s.Run()
			return nil
		},
	}
)
