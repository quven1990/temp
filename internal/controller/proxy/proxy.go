package proxy

import (
	"fmt"
	"net/http"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Proxy 将当前请求代理到目标服务
func Proxy(target string) ghttp.HandlerFunc {
	return func(r *ghttp.Request) {
		path := r.Request.URL.Path
		query := r.Request.URL.RawQuery
		url := fmt.Sprintf("%s%s", target, path)
		if query != "" {
			url += "?" + query
		}

		// 构造请求头
		headers := make(map[string]string)
		for k, v := range r.Request.Header {
			if len(v) > 0 {
				headers[k] = v[0]
			}
		}

		client := g.Client()
		req := client.SetHeaderMap(headers)

		var resp *gclient.Response
		var err error

		switch r.Method {
		case http.MethodGet:
			resp, err = req.Get(r.Context(), url)
		case http.MethodPost:
			resp, err = req.Post(r.Context(), url, r.GetBody())
		case http.MethodPut:
			resp, err = req.Put(r.Context(), url, r.GetBody())
		case http.MethodDelete:
			resp, err = req.Delete(r.Context(), url)
		default:
			r.Response.WriteStatus(http.StatusMethodNotAllowed)
			return
		}

		if err != nil {
			r.Response.WriteStatusExit(http.StatusBadGateway, g.Map{
				"error": fmt.Sprintf("调用后端服务失败: %v", err),
			})
			return
		}
		defer resp.Close()

		data := resp.ReadAll()
		r.Response.Write(data)
	}
}
