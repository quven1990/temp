package getstationnoteapi

import (
	"context"
	"fmt"

	"gf_api/internal/db"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	_ "github.com/lib/pq" // PostgreSQL 驱动
)

// 获取台站注意事项，从数据库中的notes表中取 ldc 20251017 已完成
// Register 把当前模块的所有路由注册到 group
func Register(group *ghttp.RouterGroup) {
	group.GET("/Resource/GetNotes", GetStationNote)

}

func GetStationNote(r *ghttp.Request) {
	ctx := context.Background()

	// 从 URL 参数中获取 StationId
	stationId := r.Get("StationId").String()
	if stationId == "" {
		r.Response.WriteJson(g.Map{
			"error": "缺少参数 StationId",
		})
		return
	}

	// 查询 note 表的所有字段
	sql := `SELECT * FROM note WHERE station_id =? limit 20`
	fmt.Printf("PgDB 是否为 nil？ %v\n", db.PgDB == nil)
	results, err := db.PgDB.Query(ctx, sql, stationId)
	if err != nil {
		r.Response.WriteJson(g.Map{
			"error": fmt.Sprintf("查询失败: %v", err),
		})
		return
	}

	// 如果没查到数据
	if len(results) == 0 {
		r.Response.WriteJson(g.Map{
			"message": fmt.Sprintf("未找到 StationId=%s 的数据", stationId),
		})
		return
	}

	// 返回查询结果
	r.Response.WriteJson(g.Map{
		"StationId": stationId,
		"data":      results,
	})
}
