package service

import (
	"context"
	"fmt"
	"net/url"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
)

// ExternalService 外部服务调用封装
type ExternalService struct {
	baseURL string
}

// NewExternalService 创建外部服务实例
func NewExternalService(ctx context.Context) *ExternalService {
	baseURL := g.Cfg().MustGet(ctx, "external.hisDataService.baseURL", "http://111.111.8.242:8003").String()
	return &ExternalService{
		baseURL: baseURL,
	}
}

// CallAPI 通用方法：调用外部API
// path: API路径，例如 "/HisData/DevHis"
// params: 查询参数map，例如 map[string]string{"positionId": "0101", "pageIndex": "1"}
// 返回: JSON响应数据
func (s *ExternalService) CallAPI(ctx context.Context, path string, params map[string]string) (map[string]interface{}, error) {
	// 构建完整URL
	requestURL := fmt.Sprintf("%s%s", s.baseURL, path)

	// 构建查询参数
	queryParams := url.Values{}
	for key, value := range params {
		queryParams.Set(key, value)
	}

	// 如果有参数，添加到URL
	fullURL := requestURL
	if len(queryParams) > 0 {
		fullURL = fmt.Sprintf("%s?%s", requestURL, queryParams.Encode())
	}

	// 调用外部接口
	resp, err := g.Client().Get(ctx, fullURL)
	if err != nil {
		return nil, fmt.Errorf("调用外部接口失败: %w", err)
	}
	defer resp.Close()

	// 读取响应内容
	body := resp.ReadAll()

	// 解析JSON响应
	var result map[string]interface{}
	if err := gjson.DecodeTo(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w, 原始响应: %s", err, string(body))
	}

	return result, nil
}

// GetDevHis 获取设备历史数据
// positionId: 设备位置ID
// beginTime: 开始时间，格式：YYYY-MM-DD HH:mm:ss
// endTime: 结束时间，格式：YYYY-MM-DD HH:mm:ss
// pageIndex: 页码，默认为1
// pageSize: 每页大小，默认为20
func (s *ExternalService) GetDevHis(ctx context.Context, positionId, beginTime, endTime string, pageIndex, pageSize int) (map[string]interface{}, error) {
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	// 构建JSON请求体，使用下划线命名格式
	bodyData := map[string]interface{}{
		"position_id": positionId,
		"begin_time":  beginTime,
		"end_time":    endTime,
		"page_index":  fmt.Sprintf("%d", pageIndex),
		"page_size":   fmt.Sprintf("%d", pageSize),
	}

	return s.CallAPIPost(ctx, "/HisData/DevHis", bodyData)
}

// CallAPIPost 通用方法：调用外部API（POST请求，发送JSON数据）
// path: API路径，例如 "/HisData/AlarmHis"
// bodyData: JSON请求体数据
// 返回: JSON响应数据
func (s *ExternalService) CallAPIPost(ctx context.Context, path string, bodyData map[string]interface{}) (map[string]interface{}, error) {
	// 构建完整URL
	requestURL := fmt.Sprintf("%s%s", s.baseURL, path)

	// 调用外部接口（POST请求，发送JSON）
	resp, err := g.Client().Post(ctx, requestURL, bodyData)
	if err != nil {
		return nil, fmt.Errorf("调用外部接口失败: %w", err)
	}
	defer resp.Close()

	// 读取响应内容
	body := resp.ReadAll()

	// 解析JSON响应
	var result map[string]interface{}
	if err := gjson.DecodeTo(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w, 原始响应: %s", err, string(body))
	}

	return result, nil
}

// GetAlarmHis 获取告警历史数据
// positionId: 设备位置ID
// beginTime: 开始时间，格式：YYYY-MM-DD HH:mm:ss
// endTime: 结束时间，格式：YYYY-MM-DD HH:mm:ss
// pageIndex: 页码，默认为1
// pageSize: 每页大小，默认为20
func (s *ExternalService) GetAlarmHis(ctx context.Context, positionId, beginTime, endTime string, pageIndex, pageSize int) (map[string]interface{}, error) {
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	// 构建JSON请求体，使用下划线命名格式
	bodyData := map[string]interface{}{
		"position_id": positionId,
		"begin_time":  beginTime,
		"end_time":    endTime,
		"page_index":  fmt.Sprintf("%d", pageIndex),
		"page_size":   fmt.Sprintf("%d", pageSize),
	}

	return s.CallAPIPost(ctx, "/HisData/AlarmHis", bodyData)
}

// NewParamService 创建参数服务实例
// 用于调用参数相关的第三方接口，使用不同的baseURL
func NewParamService(ctx context.Context) *ExternalService {
	baseURL := g.Cfg().MustGet(ctx, "external.paramService.baseURL", "http://10.70.1.190:30804").String()
	return &ExternalService{
		baseURL: baseURL,
	}
}

// GetChangeHis 获取参数变更历史
// cmdId: 命令ID
// pageIndex: 页码，默认为1
// pageSize: 每页大小，默认为20
func (s *ExternalService) GetChangeHis(ctx context.Context, cmdId string, pageIndex, pageSize int) (map[string]interface{}, error) {
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	// 构建查询参数
	params := map[string]string{
		"cmdId":     cmdId,
		"pageIndex": fmt.Sprintf("%d", pageIndex),
		"pageSize":  fmt.Sprintf("%d", pageSize),
	}

	// 调用外部接口（GET请求，URL查询参数）
	return s.CallAPI(ctx, "/Param/ChangeHis", params)
}

// GetSetDft 设置默认参数
// id: 参数ID
func (s *ExternalService) GetSetDft(ctx context.Context, id string) (map[string]interface{}, error) {
	// 构建查询参数
	params := map[string]string{
		"id": id,
	}

	// 调用外部接口（GET请求，URL查询参数）
	return s.CallAPI(ctx, "/Param/SetDft", params)
}

// GetReadAll 读取所有参数
// 无参数，直接调用第三方接口
func (s *ExternalService) GetReadAll(ctx context.Context) (map[string]interface{}, error) {
	// 无查询参数，传递空的参数map
	params := map[string]string{}

	// 调用外部接口（GET请求，无URL查询参数）
	return s.CallAPI(ctx, "/Param/ReadAll", params)
}

// GetReadDft 读取默认参数
// 无参数，直接调用第三方接口
func (s *ExternalService) GetReadDft(ctx context.Context) (map[string]interface{}, error) {
	// 无查询参数，传递空的参数map
	params := map[string]string{}

	// 调用外部接口（GET请求，无URL查询参数）
	return s.CallAPI(ctx, "/Param/ReadDft", params)
}

// GetSetTo 设置参数到指定值
// ids: 参数ID，可以是单个或多个（多个时逗号隔开），例如：ids=0 或 ids=0,1,2
func (s *ExternalService) GetSetTo(ctx context.Context, ids string) (map[string]interface{}, error) {
	// 构建查询参数
	params := map[string]string{
		"ids": ids,
	}

	// 调用外部接口（GET请求，URL查询参数）
	return s.CallAPI(ctx, "/Param/SetTo", params)
}

// GetReadView 读取视图参数
// 无参数，直接调用第三方接口
func (s *ExternalService) GetReadView(ctx context.Context) (map[string]interface{}, error) {
	// 无查询参数，传递空的参数map
	params := map[string]string{}

	// 调用外部接口（GET请求，无URL查询参数）
	return s.CallAPI(ctx, "/Param/ReadView", params)
}

// GetSyncFrom 从外部同步参数
// 无参数，直接调用第三方接口
func (s *ExternalService) GetSyncFrom(ctx context.Context) (map[string]interface{}, error) {
	// 无查询参数，传递空的参数map
	params := map[string]string{}

	// 调用外部接口（GET请求，无URL查询参数）
	return s.CallAPI(ctx, "/Param/SyncFrom", params)
}

// PostSyncTo 同步参数到外部
// bodyData: 比对数据，JSON格式
func (s *ExternalService) PostSyncTo(ctx context.Context, bodyData map[string]interface{}) (map[string]interface{}, error) {
	// 调用外部接口（POST请求，发送JSON数据）
	return s.CallAPIPost(ctx, "/Param/SyncTo", bodyData)
}
