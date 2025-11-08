package childsysdataapi

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

//获取子系统数据 ldc 20250901 已完成，但是因为获取总览数据的接口比较慢，因此本接口也慢

// Register 把当前模块的所有路由注册到 group
func Register(group *ghttp.RouterGroup) {
	group.GET("/Basic/ProgramSystemDataSubscribe", GetProgramSystemDataSubscribe)
}

// 从台站总览数据接口的子系统节点取数据，必须先有台站总览数据接口。
func GetProgramSystemDataSubscribe(r *ghttp.Request) {
	ctx := context.Background()
	timeStart := time.Now()

	// 从 URL 参数中获取 StationId和SubSystem
	stationId := r.Get("StationId").String()
	subSystem := r.Get("SubSystem").String()
	if stationId == "" && subSystem == "" {
		r.Response.WriteJson(g.Map{
			"Result":      false,
			"Message":     "缺少参数 StationId和SubSystem",
			"Content":     nil,
			"MachineName": "",
		})

		return
	}

	// 2调用内部接口 /api/alldata 获取总览数据
	url := fmt.Sprintf("http://127.0.0.1:8001/api/Basic/OverViewData?stationId=%s", stationId)
	fmt.Println("URL:", url)
	resp, err := g.Client().Get(ctx, url)
	if err != nil {
		r.Response.WriteJson(g.Map{
			"Result":      false,
			"Message":     fmt.Sprintf("请求 OverViewData 接口失败: %v", err),
			"Content":     nil,
			"MachineName": "",
		})
		return
	}
	defer resp.Close()

	// 3解析返回的 JSON
	body := resp.ReadAll()

	// 4解析 JSON 数据变为Go的map对象
	var overviewData map[string]interface{}
	if err := json.Unmarshal(body, &overviewData); err != nil {
		r.Response.WriteJson(g.Map{
			"Result":      false,
			"Message":     fmt.Sprintf("解析 JSON 失败: %v", err),
			"Content":     nil,
			"MachineName": "",
		})
		return
	}

	// 5查找 SubSystem 节点（不区分大小写）
	contentRaw, ok := overviewData["Content"]
	if !ok {
		r.Response.WriteJson(g.Map{
			"Result":      false,
			"Message":     "响应数据中缺少 Content 节点",
			"Content":     nil,
			"MachineName": "",
		})
		return
	}

	// Content 通常是一个 map[string]interface{}
	contentMap, ok := contentRaw.(map[string]interface{})
	if !ok {
		r.Response.WriteJson(g.Map{
			"Result":      false,
			"Message":     "Content 节点不是合法的 JSON 对象",
			"Content":     nil,
			"MachineName": "",
		})
		return
	}

	// 开始在 Content 里查找 subSystem
	var target interface{}
	found := false

	for k, v := range contentMap {

		if strings.EqualFold(k, subSystem) {
			target = v
			found = true
			break
		}
	}

	if !found {
		r.Response.WriteJson(g.Map{
			"Result":      false,
			"Message":     fmt.Sprintf("在 Content 中未找到子系统 %s", subSystem),
			"Content":     nil,
			"MachineName": "",
		})
		return
	}

	// 6返回结果
	if !found {
		r.Response.WriteJson(g.Map{
			"Result":      false,
			"Message":     fmt.Sprintf("未找到 SubSystem=%s 的数据", subSystem),
			"Content":     nil,
			"MachineName": "",
		})
		return
	}
	timeEnd := time.Since(timeStart)
	r.Response.WriteJson(g.Map{
		"Result":      true,
		"time":        timeEnd.Milliseconds(),
		"Message":     "",
		"station_id":  stationId,
		"subSystem":   subSystem,
		"Content":     target,
		"MachineName": "",
	})
}
