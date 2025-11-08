package main

import (
	"fmt"
	//"net/url"
	"gf_api/internal/db"
	_ "gf_api/internal/packed"
	//"github.com/gogf/gf/v2/database/gdb"

	"gf_api/internal/cmd"
	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"github.com/gogf/gf/v2/os/gctx"
)

func main() {
	/*GoFrame 的 cmd.Main.Run(ctx) 是阻塞调用，它会直接启动 HTTP 服务器，并占用当前 goroutine。
	所以初始化 Redis 和 PgDB 的代码 可能在 HTTP server 启动之后才执行，或者因为 goroutine 调度导致接口处理时 PgDB 尚未完成初始化。
	这就解释了为什么接口请求会报 nil pointer，但程序启动日志显示 PgDB 初始化成功。*/

	ctx := gctx.New() //携带 GoFrame 特有的全局对象，比如 Request、Response、日志、配置等。

	db.InitRedis()

	//初始化pgsql数据库
	if err := db.InitPostgresNew(ctx); err != nil {
		fmt.Println(fmt.Errorf("pgsql数据库初始化失败: %w", err))
	}
	//fmt.Println("PgDB 是否为空？", db.PgDB == nil)

	cmd.Main.Run(ctx) //这个必须在数据库初始化后才能执行，否则会阻塞

	// 初始化pgsql数据库
	// db, err := db.InitPostgres(ctx)
	// if err != nil {
	// 	fmt.Println("failed to connect postgres: %v", err)
	// }
	// // 测试查询
	// res, err := db.Query(ctx, "SELECT * FROM note where station_id='0101'")
	// if err != nil {
	// 	fmt.Println("query failed: %v", err)
	// }
	// if len(res) == 0 {
	// 	fmt.Println("no rows found")
	// 	return
	// }
	// // 遍历结果
	// for _, record := range res {
	// 	fmt.Println("Row:", record.Map())
	// }

	// 初始化 Redis
	//db.InitRedis()
	// 测试redis写入和读取
	// err := db.Redis.Set(ctx, "gf_redis_test", "hello world", 0).Err()
	// if err != nil {
	// 	fmt.Println("写入失败:", err)
	// 	return
	// }
	// val, err := db.Redis.Get(ctx, "gf_redis_test").Result()
	// if err != nil {
	// 	fmt.Println("读取失败:", err)
	// 	return
	// }
	// fmt.Println("读取 Redis 值:", val)
}
