package logic

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
)

// LogInsertParams 日志插入参数
type LogInsertParams struct {
	LogType     string // 日志类型，例如："info"
	ReqNum      string // 请求编号，例如："123"
	Level       string // 日志级别，例如："info"
	Status      string // 状态，例如："NORMAL"
	ChildSystem string // 子系统，例如："UserSystem"
	Module      string // 模块，例如："AuthModule"
	PositionId  string // 位置ID，例如："test"
	LogContent  string // 日志内容，例如："test"
	LogUser     string // 日志用户，例如："admin"
	LogTime     string // 日志时间，格式：YYYY-MM-DD HH:mm:ss，例如："2025-12-20 00:00:00"
}

// InsertLog 插入日志
// 调用第三方日志服务接口写入日志
// 返回: JSON响应数据和错误信息
func InsertLog(ctx context.Context, params LogInsertParams) (map[string]interface{}, error) {
	// 从配置中获取日志服务的基础URL，如果没有配置则使用默认值
	logServiceURL := g.Cfg().MustGet(ctx, "external.logService.baseURL", "http://111.111.8.89:30800").String()

	// 构建完整的接口URL
	requestURL := fmt.Sprintf("%s/api/log/insert", logServiceURL)

	// 构建 form-data 参数
	formData := map[string]string{
		"logType":     params.LogType,
		"reqNum":      params.ReqNum,
		"level":       params.Level,
		"status":      params.Status,
		"childSystem": params.ChildSystem,
		"module":      params.Module,
		"positionId":  params.PositionId,
		"logContent":  params.LogContent,
		"logUser":     params.LogUser,
		"logTime":     params.LogTime,
	}

	// 调用第三方接口（POST请求，发送form-data）
	resp, err := g.Client().Post(ctx, requestURL, formData)
	if err != nil {
		return nil, fmt.Errorf("调用日志服务接口失败: %w", err)
	}
	defer resp.Close()

	// 读取响应内容
	body := resp.ReadAll()

	// 解析JSON响应
	var result map[string]interface{}
	if err := gjson.DecodeTo(body, &result); err != nil {
		return nil, fmt.Errorf("解析日志服务响应失败: %w, 原始响应: %s", err, string(body))
	}

	return result, nil
}

// InsertLogSimple 简化版日志插入方法
// 使用最常用的参数，其他参数使用默认值
func InsertLogSimple(ctx context.Context, level, module, logContent, logUser string) (map[string]interface{}, error) {
	params := LogInsertParams{
		LogType:     "info",
		ReqNum:      "",
		Level:       level,
		Status:      "NORMAL",
		ChildSystem: "UserSystem",
		Module:      module,
		PositionId:  "",
		LogContent:  logContent,
		LogUser:     logUser,
		LogTime:     time.Now().Format("2006-01-02 15:04:05"),
	}
	return InsertLog(ctx, params)
}
