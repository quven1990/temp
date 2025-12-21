package alarmhisapi

// 设备历史数据接口 - 调用外部服务获取设备历史记录
import (
	"context"
	"gf_api/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Register 把当前模块的所有路由注册到 group
func Register(group *ghttp.RouterGroup) {
	group.GET("/DevHis", GetDevHis)
	group.GET("/Resource/AlarmHis", GetAlarmHis) // 该接口依赖方待定
}

// GetDevHis 获取设备历史数据
func GetDevHis(r *ghttp.Request) {
	ctx := context.Background()
	// 从URL参数中获取positionId
	positionId := r.Get("positionId").String()
	if positionId == "" {
		r.Response.WriteJson(g.Map{
			"code":    400,
			"message": "缺少参数 positionId",
			"data":    nil,
		})
		return
	}

	// 获取时间参数
	beginTime := r.Get("beginTime").String()
	endTime := r.Get("endTime").String()
	if beginTime == "" {
		r.Response.WriteJson(g.Map{
			"code":    400,
			"message": "缺少参数 beginTime",
			"data":    nil,
		})
		return
	}
	if endTime == "" {
		r.Response.WriteJson(g.Map{
			"code":    400,
			"message": "缺少参数 endTime",
			"data":    nil,
		})
		return
	}

	// 获取可选的分页参数
	pageIndex := r.Get("pageIndex", "1").Int()
	pageSize := r.Get("pageSize", "20").Int()

	// 创建外部服务实例
	externalService := service.NewExternalService(ctx)

	// 调用外部服务
	result, err := externalService.GetDevHis(ctx, positionId, beginTime, endTime, pageIndex, pageSize)
	if err != nil {
		// 如果是业务错误响应（外部接口返回的错误），原样返回
		if errorResponse, ok := err.(*service.BusinessError); ok {
			r.Response.WriteJson(errorResponse.Response)
			return
		}
		// 其他错误（网络错误、解析错误等）返回500
		r.Response.WriteJson(g.Map{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	// 透传返回结果
	r.Response.WriteJson(result)
}

// GetAlarmHis 获取告警历史数据
// 该接口用于获取资源告警历史记录，通过调用第三方接口进行数据透传
// 注意：当前第三方接口还未确定，暂时返回模拟数据
func GetAlarmHis(r *ghttp.Request) {
	ctx := context.Background()
	_ = ctx // 预留用于后续第三方接口调用

	// 从URL参数中获取positionId（设备位置ID）
	positionId := r.Get("positionId").String()
	if positionId == "" {
		r.Response.WriteJson(g.Map{
			"code":    400,
			"message": "缺少参数 positionId",
			"data":    nil,
		})
		return
	}

	// 获取可选的分页参数
	pageIndex := r.Get("pageIndex", "1").Int()
	pageSize := r.Get("pageSize", "20").Int()

	// TODO: 调用第三方接口获取告警历史数据
	// 第三方接口地址和参数格式待确定，确定后需要：
	// 1. 在 internal/service/external.go 中添加对应的调用方法（如 GetAlarmHis）
	// 2. 创建外部服务实例：externalService := service.NewExternalService(ctx)
	// 3. 调用外部服务：result, err := externalService.GetAlarmHis(ctx, positionId, pageIndex, pageSize)
	// 4. 处理错误并透传返回结果

	// 临时返回模拟数据，待第三方接口确定后替换为实际调用
	r.Response.WriteJson(g.Map{
		"code":    200,
		"message": "success",
		"data": g.Map{
			"positionId": positionId,
			"pageIndex":  pageIndex,
			"pageSize":   pageSize,
			"total":      0,
			"list":       []interface{}{},
			"note":       "该接口的第三方调用逻辑待实现，当前返回模拟数据",
		},
	})
}
