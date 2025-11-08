package entity

import "time"

// Log 日志实体,和数据库表结构对应。
type Log struct {
	Id          uint      `json:"id"          description:"日志ID"  gorm:"primaryKey"`
	LogType     string    `json:"logType" description:"日志类型"`    // 日志类型
	ReqNum      string    `json:"reqNum" description:"请求编号"`     // 请求编号
	Level       string    `json:"level" description:"日志级别"`      //日志级别
	Status      string    `json:"status" description:"是否异常"`     //是否异常
	ChildSystem string    `json:"childSystem" description:"子系统"` //子系统
	Module      string    `json:"module" description:"模块"`       //模块
	PositionId  string    `json:"positionId" description:"工位号"`  //工位号
	LogContent  string    `json:"logContent" description:"日志内容"` //日志内容
	LogUser     string    `json:"logUser" description:"用户"`      //用户
	LogTime     time.Time `json:"logTime" description:"日志时间"`    //日志时间
}
