package configapi

// 配置管理相关接口 - 透传参数配置管理接口
import (
	"context"
	"gf_api/internal/service"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Register 把当前模块的所有路由注册到 group
func Register(group *ghttp.RouterGroup) {
	group.GET("/Param/ChangeHis", GetChangeHis)
	group.GET("/Param/SetDft", GetSetDft)
	group.GET("/Param/ReadAll", GetReadAll)
	group.GET("/Param/ReadDft", GetReadDft)
	group.GET("/Param/SetTo", GetSetTo)
	group.GET("/Param/ReadView", GetReadView)
	group.GET("/Param/SyncFrom", GetSyncFrom)
	group.POST("/Param/SyncTo", PostSyncTo)
}

// GetChangeHis 获取参数变更历史
// 该接口透传调用第三方接口获取参数变更历史记录
// 请求参数：
//   - cmdId: 命令ID（必填）
//   - pageIndex: 页码，默认为1（可选）
//   - pageSize: 每页大小，默认为20（可选）
func GetChangeHis(r *ghttp.Request) {
	ctx := context.Background()

	// 从URL参数中获取cmdId（命令ID）
	cmdId := r.Get("cmdId").String()
	if cmdId == "" {
		r.Response.WriteJson(g.Map{
			"code":    400,
			"message": "缺少参数 cmdId",
			"data":    nil,
		})
		return
	}

	// 获取可选的分页参数
	pageIndex := r.Get("pageIndex", "1").Int()
	pageSize := r.Get("pageSize", "20").Int()

	// 创建外部服务实例（使用参数服务的基础URL）
	externalService := service.NewParamService(ctx)

	// 调用外部服务获取参数变更历史
	result, err := externalService.GetChangeHis(ctx, cmdId, pageIndex, pageSize)
	if err != nil {
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

// GetSetDft 设置默认参数
// 该接口透传调用第三方接口设置默认参数
// 请求参数：
//   - id: 参数ID（必填）
func GetSetDft(r *ghttp.Request) {
	ctx := context.Background()

	// 从URL参数中获取id（参数ID）
	id := r.Get("id").String()
	if id == "" {
		r.Response.WriteJson(g.Map{
			"code":    400,
			"message": "缺少参数 id",
			"data":    nil,
		})
		return
	}

	// 创建外部服务实例（使用参数服务的基础URL）
	externalService := service.NewParamService(ctx)

	// 调用外部服务设置默认参数
	result, err := externalService.GetSetDft(ctx, id)
	if err != nil {
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

// GetReadAll 读取所有参数
// 该接口透传调用第三方接口读取所有参数
// 无请求参数
func GetReadAll(r *ghttp.Request) {
	ctx := context.Background()

	// 创建外部服务实例（使用参数服务的基础URL）
	externalService := service.NewParamService(ctx)

	// 调用外部服务读取所有参数
	result, err := externalService.GetReadAll(ctx)
	if err != nil {
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

// GetReadDft 读取默认参数
// 该接口透传调用第三方接口读取默认参数
// 无请求参数
func GetReadDft(r *ghttp.Request) {
	ctx := context.Background()

	// 创建外部服务实例（使用参数服务的基础URL）
	externalService := service.NewParamService(ctx)

	// 调用外部服务读取默认参数
	result, err := externalService.GetReadDft(ctx)
	if err != nil {
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

// GetSetTo 设置参数到指定值
// 该接口透传调用第三方接口设置参数到指定值
// 请求参数：
//   - ids: 参数ID，接口"读取默认配置"结果的id（必填），多个时逗号隔开，例如：ids=0 或 ids=0,1,2
func GetSetTo(r *ghttp.Request) {
	ctx := context.Background()

	// 从URL参数中获取ids（参数ID，可以是单个或多个，多个时逗号隔开）
	ids := r.Get("ids").String()
	if ids == "" {
		r.Response.WriteJson(g.Map{
			"code":    400,
			"message": "缺少参数 ids",
			"data":    nil,
		})
		return
	}

	// 创建外部服务实例（使用参数服务的基础URL）
	externalService := service.NewParamService(ctx)

	// 调用外部服务设置参数到指定值
	result, err := externalService.GetSetTo(ctx, ids)
	if err != nil {
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

// GetReadView 读取视图参数
// 该接口透传调用第三方接口读取视图参数
// 无请求参数
func GetReadView(r *ghttp.Request) {
	ctx := context.Background()

	// 创建外部服务实例（使用参数服务的基础URL）
	externalService := service.NewParamService(ctx)

	// 调用外部服务读取视图参数
	result, err := externalService.GetReadView(ctx)
	if err != nil {
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

// GetSyncFrom 从外部同步参数
// 该接口透传调用第三方接口从外部同步参数
// 无请求参数
func GetSyncFrom(r *ghttp.Request) {
	ctx := context.Background()

	// 创建外部服务实例（使用参数服务的基础URL）
	externalService := service.NewParamService(ctx)

	// 调用外部服务从外部同步参数
	result, err := externalService.GetSyncFrom(ctx)
	if err != nil {
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

// PostSyncTo 同步参数到外部
// 该接口透传调用第三方接口同步参数到外部
// 请求方式：POST
// 请求参数（JSON格式）：比对数据，JSON格式
func PostSyncTo(r *ghttp.Request) {
	ctx := context.Background()

	// 从请求体中获取JSON数据（比对数据）
	bodyBytes := r.GetBody()
	if len(bodyBytes) == 0 {
		r.Response.WriteJson(g.Map{
			"code":    400,
			"message": "缺少请求体数据",
			"data":    nil,
		})
		return
	}

	// 解析JSON请求体
	var bodyData map[string]interface{}
	if err := gjson.DecodeTo(bodyBytes, &bodyData); err != nil {
		r.Response.WriteJson(g.Map{
			"code":    400,
			"message": "请求体格式错误，需要JSON格式: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 创建外部服务实例（使用参数服务的基础URL）
	externalService := service.NewParamService(ctx)

	// 调用外部服务同步参数到外部
	result, err := externalService.PostSyncTo(ctx, bodyData)
	if err != nil {
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

