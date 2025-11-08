package controlsysapi

//台站客户端-操作命令下发  ldc 20251023
import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Register 把当前模块的所有路由注册到 group
func Register(group *ghttp.RouterGroup) {
	group.POST("/Resource/IssueOperateNew", IssueOperate)
}

func IssueOperate(r *ghttp.Request) {
	ctx := context.Background()

	// 读取 POST 请求体的 JSON 数据示例：
	var reqData map[string]interface{}
	if err := r.Parse(&reqData); err != nil {
		r.Response.WriteJson(g.Map{
			"result":  "error",
			"message": fmt.Sprintf("请求参数解析失败: %v", err),
		})
		return
	}

	fmt.Println("收到的请求数据:", reqData)

	// 取出字段
	positionId := reqData["positionId"]
	name := reqData["name"]
	para := reqData["para"]
	paranew := reqData["paranew"]
	frequency := reqData["frequency"]
	clientIp := reqData["clientIp"]
	userCode := reqData["userCode"]
	UserName := reqData["UserName"]
	realName := reqData["realName"]
	AgentType := reqData["AgentType"]
	fmt.Println("positionId:", positionId)
	fmt.Println("name:", name)
	fmt.Println("para:", para)
	fmt.Println("paranew:", paranew)
	fmt.Println("frequency:", frequency)
	fmt.Println("clientIp:", clientIp)
	fmt.Println("userCode:", userCode)
	fmt.Println("UserName:", UserName)
	fmt.Println("realName:", realName)
	fmt.Println("AgentType:", AgentType)

	// 3️⃣ 组织要转发的请求数据
	postData := g.Map{
		"positionId": positionId,
		"name":       name,
		"para":       para,
		"paranew":    paranew,
		"frequency":  frequency,
		"clientIp":   clientIp,
		"userCode":   userCode,
		"UserName":   UserName,
		"realName":   realName,
		"AgentType":  AgentType,
	}

	// 4️⃣ 发起 POST 请求到目标接口
	targetURL := "http://111.111.8.242:8005/api/Resource/IssueOperateNew"

	resp, err := g.Client().Post(ctx, targetURL, postData)
	if err != nil {
		r.Response.WriteJson(g.Map{
			"result":  "error",
			"message": fmt.Sprintf("转发接口请求失败: %v", err),
		})
		return
	}
	defer resp.Close()

	// 5️⃣ 读取目标接口返回结果
	body := resp.ReadAll()

	var responseData interface{}
	if err := json.Unmarshal(body, &responseData); err != nil {
		r.Response.WriteJson(g.Map{
			"result":  "error",
			"message": fmt.Sprintf("解析返回 JSON 失败: %v", err),
			"data":    string(body), // 解析失败则返回原始字符串
		})
		return
	}

	// 成功解析 JSON，直接返回
	r.Response.WriteJson(g.Map{
		"result":  "success",
		"message": "转发成功",
		"data":    responseData, // 这里是结构化 JSON
	})
}
