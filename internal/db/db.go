package db

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

var DB *gorm.DB

// 初始化mysql
func Init() {
	var err error
	dsn := "user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
}

// 初始化redis
//var Redis *redis.ClusterClient //集群模式使用

var Redis *redis.Client //单级模式使用

func InitRedis() {
	ctx := context.Background()

	//集群模式
	// Redis = redis.NewClusterClient(&redis.ClusterOptions{
	// 	Addrs: []string{
	// 		"10.170.1.207:6379",
	// 		"10.170.1.211:6379",
	// 		"10.170.1.219:6379",
	// 		"10.170.1.223:6379",
	// 		"10.170.1.249:6379",
	// 		"10.170.1.250:6379",
	// 		// 可继续添加其他集群节点
	// 	},
	// 	Password: "gxrtbtc", // 设置集群节点的密码
	// })

	// 单机模式
	// Redis = redis.NewClient(&redis.Options{
	// 	Addr:     "111.111.8.242:6379",
	// 	Password: "123456",
	// 	DB:       0, // 默认库
	// })

	//集群和单击模式用到的连接方式
	// if err := Redis.Ping(ctx).Err(); err != nil {
	// 	fmt.Println("Redis 集群连接失败: " + err.Error())
	// } else {
	// 	fmt.Println("Redis 集群连接成功")
	// }

	//config.yaml提取配置的连接方式
	// 读取配置
	addr := g.Cfg().MustGet(ctx, "redis.default.addr").String()
	password := g.Cfg().MustGet(ctx, "redis.default.password").String()
	db := g.Cfg().MustGet(ctx, "redis.default.db").Int()
	poolSize := g.Cfg().MustGet(ctx, "redis.default.poolSize").Int()
	// 初始化 Redis
	Redis = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
		PoolSize: poolSize,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := Redis.Ping(ctx).Err(); err != nil {
		fmt.Println("Redis connect failed: %v", err)
	}
	fmt.Println("Redis connected successfully!")

}

var PgDB gdb.DB // 全局可复用的数据库连接

// InitPgsql 初始化 PostgreSQL 数据库,ldc 20250904
func InitPostgresNew(ctx context.Context) error {
	//在函数顶部先声明 err
	var err error

	// 从配置中读取完整的连接字符串
	link := g.Cfg().MustGet(ctx, "database.default.link").String()

	if link == "" {
		return fmt.Errorf("database.default.link 配置为空，请检查 config.yaml")
	}

	// 拼接 gdb 配置节点
	node := gdb.ConfigNode{
		Link: fmt.Sprintf("pgsql:%s", link), // 注意这里要加上 "pgsql:" 前缀
	}

	PgDB, err = gdb.New(node)
	if err != nil {
		return fmt.Errorf("failed to connect postgres: %w", err)
	}

	return nil
}

func InitPostgres(ctx context.Context) (gdb.DB, error) {
	// 从配置中读取完整的连接字符串
	link := g.Cfg().MustGet(ctx, "database.default.link").String()

	if link == "" {
		return nil, fmt.Errorf("database.default.link 配置为空，请检查 config.yaml")
	}

	// 拼接 gdb 配置节点
	node := gdb.ConfigNode{
		Link: fmt.Sprintf("pgsql:%s", link), // 注意这里要加上 "pgsql:" 前缀
	}

	db, err := gdb.New(node)
	if err != nil {
		return nil, fmt.Errorf("failed to connect postgres: %w", err)
	}

	return db, nil
}
