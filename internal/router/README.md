# 路由管理说明

## 路由统一管理

所有API路由现在统一在 `internal/router/router.go` 文件中管理，方便查找和维护。

## 文件结构

```
internal/
├── router/
│   └── router.go          # 统一路由注册文件
├── cmd/
│   └── cmd.go             # 服务启动入口，调用 router.RegisterAllRoutes()
└── controller/            # 各个接口控制器
    ├── alarm_his_api/
    ├── api/
    └── ...
```

## 如何添加新路由

### 步骤1：在对应的controller中实现接口

例如：`internal/controller/your_api/your_handler.go`

```go
package yourapi

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

func Register(group *ghttp.RouterGroup) {
	group.GET("/YourPath", YourHandler)
}

func YourHandler(r *ghttp.Request) {
	// 处理逻辑
}
```

### 步骤2：在 router.go 中注册

编辑 `internal/router/router.go`：

```go
import (
	yourapi "gf_api/internal/controller/your_api"
	// ... 其他导入
)

func RegisterAllRoutes(group *ghttp.RouterGroup) {
	// ... 现有路由
	
	// 你的新接口
	yourapi.Register(group)
}
```

## 路由分类

当前路由按功能分类：

- **Basic 相关接口**：基础功能接口（时间、编号等）
- **Resource 相关接口**：资源相关接口（历史数据、日志等）
- **预留扩展区域**：后续新增接口的位置

## 优势

1. ✅ **集中管理**：所有路由在一个文件中，一目了然
2. ✅ **易于维护**：新增路由只需在一个地方修改
3. ✅ **清晰分类**：按功能分类，方便查找
4. ✅ **代码简洁**：cmd.go 文件更简洁，只负责启动服务

