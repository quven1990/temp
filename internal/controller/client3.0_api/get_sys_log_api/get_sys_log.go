package getsyslogapi

//获取用户操作信息  ldc 20251022 已完成
import (
	"context"
	"fmt"

	"gf_api/internal/db"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	_ "github.com/lib/pq" // PostgreSQL 驱动
)

// Register 把当前模块的所有路由注册到 group
func Register(group *ghttp.RouterGroup) {
	group.GET("/Resource/GetOpLog", GetLogData)
}

// GetLogData
func GetLogData(r *ghttp.Request) {
	ctx := context.Background()

	// 从 URL 参数中获取
	positionId := r.Get("positionId").String()
	logType := r.Get("logType").String()
	if positionId == "" || logType == "" {
		r.Response.WriteJson(g.Map{
			"return_code": 0,
			"return_msg":  "缺少参数 positionId或logType",
			"node_name":   "",
			"data":        nil,
		})
		return
	}

	// 查询 operation_log 表的所有字段
	sql := `select user_name, ip_addr, station_id, postion_id, operate_detail, operate_time, real_name, frequency, para_data, remarks from operation_log WHERE postion_id =? and log_type=? limit 20`
	fmt.Printf("PgDB 是否为 nil？ %v\n", db.PgDB == nil)
	results, err := db.PgDB.Query(ctx, sql, positionId, logType)
	if err != nil {
		r.Response.WriteJson(g.Map{
			"error": fmt.Sprintf("查询失败: %v", err),
		})
		return
	}

	// 如果没查到数据
	if len(results) == 0 {
		r.Response.WriteJson(g.Map{
			"message": fmt.Sprintf("未找到 postion_id=%s 的数据", positionId),
			"logData": "",
		})
		return
	}

	// 返回查询结果
	r.Response.WriteJson(g.Map{
		"postion_id": positionId,
		"logData":    results,
	})

}
