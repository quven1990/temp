package gethikdataapi

//获取台站所有海康威视接口信息 ldc 20251022 已完成
import (
	"context"
	"encoding/json"
	"fmt"
	"gf_api/internal/db"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/redis/go-redis/v9"
)

// Register 把当前模块的所有路由注册到 group
func Register(group *ghttp.RouterGroup) {
	group.GET("/Resource/HIKRec", GetHIKData)
}

func GetHIKData(r *ghttp.Request) {
	ctx := context.Background()
	key := "svr_stationHIKRec"
	var message = "success"
	var return_code = 1 // 默认返回成功

	// 从 URL 参数中获取 StationId
	stationId := r.Get("StationId").String()
	if stationId == "" {
		r.Response.WriteJson(g.Map{
			"stationId":   "",
			"return_code": 0,
			"return_msg":  "缺少参数 StationId",
			"node_name":   "",
			"data":        nil,
		})
		return
	}

	// 获取 Redis key 类型
	keyType, err := db.Redis.Type(ctx, key).Result()
	if err != nil {
		r.Response.WriteJson(g.Map{
			"stationId":   "",
			"return_code": 0,
			"return_msg":  fmt.Sprintf("查询 Redis key 类型出错: %v", err),
			"node_name":   "",
			"data":        nil,
		})
		return
	}

	var val string
	switch keyType {
	case "string":
		val, err = db.Redis.Get(ctx, key).Result()

	case "hash":
		// 获取整个 hash，转成 JSON
		hashData, errH := db.Redis.HGetAll(ctx, key).Result()
		if errH != nil {
			err = errH
			break
		}
		jsonBytes, _ := json.Marshal(hashData)
		val = string(jsonBytes)

	default:
		r.Response.WriteJson(g.Map{
			"stationId":   "",
			"return_code": 0,
			"return_msg":  fmt.Sprintf("Redis key 类型不支持: %s", keyType),
			"node_name":   "",
			"data":        nil,
		})
		return
	}

	// 如果取值出错
	if err != nil {
		if err == redis.Nil {
			message = fmt.Sprintf("Redis key '%s' 不存在", key)
		} else {
			message = err.Error()
		}
		r.Response.WriteJson(g.Map{
			"stationId":   "",
			"return_code": 0,
			"return_msg":  message,
			"node_name":   "",
			"data":        nil,
		})
		return
	}

	// 解析 JSON
	var resultMap map[string]interface{}
	if err := json.Unmarshal([]byte(val), &resultMap); err != nil {
		r.Response.WriteJson(g.Map{
			"stationId":   "",
			"return_code": 0,
			"return_msg":  fmt.Sprintf("Redis value 不是合法 JSON: %v", err),
			"node_name":   "",
			"data":        nil,
		})
		return
	}

	// 查找对应 stationId 子集
	stationData, ok := resultMap[stationId]
	if !ok {
		r.Response.WriteJson(g.Map{
			"stationId":   stationId,
			"return_code": 0,
			"return_msg":  fmt.Sprintf("未找到 stationId=%s 的数据", stationId),
			"node_name":   "",
			"data":        nil,
		})
		return
	}

	// 判断 stationData 是否是字符串类型（即 JSON 字符串）
	var parsedData interface{}
	switch v := stationData.(type) {
	case string:
		// 如果是字符串类型的 JSON，再解析一次
		if err := json.Unmarshal([]byte(v), &parsedData); err != nil {
			// 如果解析失败，就直接返回原字符串
			parsedData = v
		}
	default:
		parsedData = v
	}

	// 正常返回
	r.Response.WriteJson(g.Map{
		"stationId":   stationId,
		"return_msg":  message,
		"return_code": return_code,
		"node_name":   "",
		"data":        parsedData,
	})
}
