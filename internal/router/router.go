package router

import (
	alarmhisapi "gf_api/internal/controller/alarm_his_api"
	api "gf_api/internal/controller/api"
	childsysdataapi "gf_api/internal/controller/client3.0_api/child_sys_data_api"
	childsysnumber "gf_api/internal/controller/client3.0_api/child_sys_number_api"
	controlsysapi "gf_api/internal/controller/client3.0_api/control_sys_api"
	gethikdataapi "gf_api/internal/controller/client3.0_api/get_hik_data_api"
	getstationfrqprogramapi "gf_api/internal/controller/client3.0_api/get_station_frq_program_api"
	getstationnoteapi "gf_api/internal/controller/client3.0_api/get_station_note_api"
	getsyslogapi "gf_api/internal/controller/client3.0_api/get_sys_log_api"
	gettimeapi "gf_api/internal/controller/client3.0_api/get_time_api"

	"github.com/gogf/gf/v2/net/ghttp"
)

// RegisterAllRoutes 注册所有路由
// 统一管理所有API路由，方便查找和维护
//
// 📋 所有路由地址列表请查看项目根目录的 ROUTES.md 文件
//
// 路由分类：
// - Basic 相关接口：基础功能接口（时间、编号、台站信息等）
// - Resource 相关接口：资源相关接口（历史数据、日志、控制等）
func RegisterAllRoutes(group *ghttp.RouterGroup) {
	// ==================== Basic 相关接口 ====================
	// GET /api/Basic/ServiceTime - 获取服务器时间
	gettimeapi.Register(group)

	// GET /api/Basic/StationSubSystem - 获取子系统编号
	childsysnumber.Register(group)

	// GET /api/Basic/OverViewData - 获取台站总览数据
	// GET /api/Basic/AllStation - 获取所有台站信息
	// GET /api/Basic/AllStationId - 获取所有台站ID
	api.Register(group)

	// GET /api/Basic/ProgramSystemDataSubscribe - 获取子系统信息
	childsysdataapi.Register(group)

	// GET /api/Basic/GetStationFrq - 获取台站的所有频率和节目名称
	getstationfrqprogramapi.Register(group)

	// ==================== Resource 相关接口 ====================
	// GET /api/DevHis - 获取设备历史数据
	alarmhisapi.Register(group)

	// GET /api/Resource/GetNotes - 获取台站注意事项
	getstationnoteapi.Register(group)

	// GET /api/Resource/HIKRec - 获取台站海康威视数据
	gethikdataapi.Register(group)

	// GET /api/Resource/GetOpLog - 获取用户操作日志信息
	getsyslogapi.Register(group)

	// POST /api/Resource/IssueOperateNew - 台站客户端的下发控制
	controlsysapi.Register(group)

	// ==================== 预留扩展区域 ====================
	// 后续新增接口请在此处添加，并添加相应注释说明
	// 同时请在项目根目录的 ROUTES.md 文件中添加路由信息
}
