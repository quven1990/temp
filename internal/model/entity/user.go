package entity

import "github.com/gogf/gf/v2/os/gtime"

// User 用户实体
type User struct {
	Id        uint        `json:"id"          description:"用户ID"`
	Username  string      `json:"username"    description:"用户名"`
	Password  string      `json:"password"    description:"密码"`
	CreatedAt *gtime.Time `json:"created_at"  description:"创建时间"`
	UpdatedAt *gtime.Time `json:"updated_at"  description:"更新时间"`
}
