package gettimeapi

//包名不建议用下划线
//获取当前服务器时间
import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Register 把当前模块的所有路由注册到 group
func Register(group *ghttp.RouterGroup) {
	group.GET("/Basic/ServiceTime", GetServerTime)
}

// GetServerTime 返回当前服务器时间
func GetServerTime(r *ghttp.Request) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	r.Response.WriteJson(g.Map{
		"server_time": currentTime,
	})
}
