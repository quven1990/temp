package getstationfrqprogramapi

import (
	"context"
	"encoding/json"
	"fmt"
	"gf_api/internal/db"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/redis/go-redis/v9"
)

// 获取某个台站的所有频率和节目名称 ldc 20251021 已完成
// Register 把当前模块的所有路由注册到 group
func Register(group *ghttp.RouterGroup) {
	group.GET("/Basic/GetStationFrq", GetStationFrq)

}

func GetStationFrq(r *ghttp.Request) {
	ctx := context.Background()
	key := "svr_stationFrqAndProgram"
	var result = "success"
	var message string

	// 从 URL 参数中获取 StationId
	stationId := r.Get("StationId").String()
	if stationId == "" {
		r.Response.WriteJson(g.Map{
			"stationId":   "",
			"result":      "error",
			"message":     "缺少参数 StationId",
			"MachineName": "",
			"Content":     nil,
		})
		return
	}

	val, err := db.Redis.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil { // 注意这里redis.Nil,它是 github.com/redis/go-redis/v9 包里定义
			// key 不存在
			message = fmt.Sprintf("Redis key '%s' 不存在", key)
		} else {
			message = err.Error()
		}

		r.Response.WriteJson(g.Map{
			"stationId":   "",
			"result":      result,
			"message":     message,
			"MachineName": "",
			"Content":     nil,
		})
		return
	}

	// 解析 Redis 返回的 JSON
	var resultMap map[string]interface{}
	if err := json.Unmarshal([]byte(val), &resultMap); err != nil {
		r.Response.WriteJson(g.Map{
			"stationId":   "",
			"result":      "error",
			"message":     fmt.Sprintf("Redis value 不是合法 JSON: %v", err),
			"MachineName": "",
			"Content":     nil,
		})
		return
	}

	// 查找对应 stationId 子集
	stationData, ok := resultMap[stationId]
	if !ok {
		r.Response.WriteJson(g.Map{
			"stationId":   stationId,
			"result":      "error",
			"message":     fmt.Sprintf("未找到 stationId=%s 的数据", stationId),
			"MachineName": "",
			"Content":     nil,
		})
		return
	}

	// 成功获取值
	r.Response.WriteJson(g.Map{
		"stationId":   stationId,
		"result":      result,
		"message":     message, //异常信息
		"MachineName": "",
		"Content":     stationData,
	})
}
