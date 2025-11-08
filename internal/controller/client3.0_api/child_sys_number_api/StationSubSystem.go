package childsysnumber

//获取子系统编号 ldc 20250828 已完成
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
	group.GET("/Basic/StationSubSystem", GetStationSubSystemNumber)

}

// 取redis中的key对应的值，然后根据stationId获取对应的子系统编号json字符串
func GetStationSubSystemNumber(r *ghttp.Request) {
	ctx := context.Background()
	key := "svr_stationSubSystemMatch"

	// 从 URL 参数中获取 StationId，例如 /api/Basic/StationSubSystem?StationId=0101
	stationId := r.Get("StationId").String()
	if stationId == "" {
		r.Response.WriteJson(g.Map{
			"error": "缺少参数 StationId",
		})
		return
	}

	val, err := db.Redis.Get(ctx, key).Result()

	if err != nil {
		if err == redis.Nil { // 注意这里redis.Nil,它是 github.com/redis/go-redis/v9 包里定义
			// key 不存在
			r.Response.WriteJson(g.Map{
				"error": fmt.Sprintf("Redis key '%s' 不存在", key),
			})
			return
		}
		// 其他错误
		r.Response.WriteJson(g.Map{
			"error": err.Error(),
		})
		return
	}

	// Redis 里是 JSON 字符串 → 转成 map
	//取出对应的redis中的key对应的值，然后根据stationId 在结果集中取出stationId_id这个属性值
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		r.Response.WriteJson(g.Map{
			"error": "JSON 解析失败",
			"raw":   val,
		})
		return
	}

	// 拼出 stationId_id 键名
	subKey := fmt.Sprintf("%s_id", stationId)
	subVal, ok := data[subKey]
	if !ok {
		r.Response.WriteJson(g.Map{
			"error": fmt.Sprintf("未找到键 '%s'", subKey),
		})
		return
	}

	// 成功获取值
	r.Response.WriteJson(subVal)
}
