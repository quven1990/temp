# API路由列表

> 所有API路由地址统一在此文件中管理，方便查找和维护
> 
> 基础路径：`http://localhost:8001/api`

---

## 📋 Basic 相关接口

### 1. 获取服务器时间
- **路径**: `GET /api/Basic/ServiceTime`
- **说明**: 返回当前服务器时间
- **参数**: 无
- **Controller**: `internal/controller/client3.0_api/get_time_api/get_time.go`

### 2. 获取子系统编号
- **路径**: `GET /api/Basic/StationSubSystem`
- **说明**: 获取台站子系统编号
- **参数**: 
  - `StationId` (必填): 台站ID
- **示例**: `/api/Basic/StationSubSystem?StationId=0101`
- **Controller**: `internal/controller/client3.0_api/child_sys_number_api/StationSubSystem.go`

### 3. 获取台站总览数据
- **路径**: `GET /api/Basic/OverViewData`
- **说明**: 获取台站总览数据
- **参数**: 
  - `stationId` (必填): 台站ID
- **示例**: `/api/Basic/OverViewData?stationId=0101`
- **Controller**: `internal/controller/api/OverViewData.go`

### 4. 获取所有台站信息
- **路径**: `GET /api/Basic/AllStation`
- **说明**: 获取所有台站信息
- **参数**: 无
- **Controller**: `internal/controller/api/OverViewData.go`

### 5. 获取所有台站ID
- **路径**: `GET /api/Basic/AllStationId`
- **说明**: 获取所有台站ID列表
- **参数**: 无
- **Controller**: `internal/controller/api/OverViewData.go`

### 6. 获取台站的所有频率和节目名称
- **路径**: `GET /api/Basic/GetStationFrq`
- **说明**: 获取台站的所有频率和节目名称
- **参数**: 
  - `StationId` (必填): 台站ID
- **示例**: `/api/Basic/GetStationFrq?StationId=0101`
- **Controller**: `internal/controller/client3.0_api/get_station_frq_program_api/get_station_frq_program.go`

### 7. 获取子系统信息
- **路径**: `GET /api/Basic/ProgramSystemDataSubscribe`
- **说明**: 获取子系统信息
- **参数**: 
  - `StationId` (必填): 台站ID
  - `SubSystem` (必填): 子系统名称
- **示例**: `/api/Basic/ProgramSystemDataSubscribe?StationId=0101&SubSystem=发射机`
- **Controller**: `internal/controller/client3.0_api/child_sys_data_api/ProgramSystemDataSubscribe.go`

---

## 📦 Resource 相关接口

### 8. 获取设备历史数据
- **路径**: `GET /api/DevHis`
- **说明**: 获取设备历史数据，调用外部服务
- **参数**: 
  - `positionId` (必填): 设备位置ID
  - `pageIndex` (可选): 页码，默认1
  - `pageSize` (可选): 每页大小，默认20
- **示例**: `/api/DevHis?positionId=0101_0x0702_2&pageIndex=1&pageSize=20`
- **Controller**: `internal/controller/alarm_his_api/alarm.go`

### 9. 获取台站注意事项
- **路径**: `GET /api/Resource/GetNotes`
- **说明**: 获取台站注意事项，从数据库notes表查询
- **参数**: 
  - `StationId` (必填): 台站ID
- **示例**: `/api/Resource/GetNotes?StationId=0101`
- **Controller**: `internal/controller/client3.0_api/get_station_note_api/get_station_note.go`

### 10. 获取台站海康威视数据
- **路径**: `GET /api/Resource/HIKRec`
- **说明**: 获取台站所有海康威视接口信息
- **参数**: 
  - `StationId` (必填): 台站ID
- **示例**: `/api/Resource/HIKRec?StationId=0101`
- **Controller**: `internal/controller/client3.0_api/get_hik_data_api/get_hik_data.go`

### 11. 获取用户操作日志信息
- **路径**: `GET /api/Resource/GetOpLog`
- **说明**: 获取用户操作日志信息，从数据库operation_log表查询
- **参数**: 
  - `positionId` (必填): 位置ID
  - `logType` (必填): 日志类型
- **示例**: `/api/Resource/GetOpLog?positionId=0101&logType=操作`
- **Controller**: `internal/controller/client3.0_api/get_sys_log_api/get_sys_log.go`

### 12. 台站客户端的下发控制
- **路径**: `POST /api/Resource/IssueOperateNew`
- **说明**: 台站客户端操作命令下发
- **参数**: JSON Body
  - `positionId`: 位置ID
  - `name`: 名称
  - `para`: 参数
  - `paranew`: 新参数
  - `frequency`: 频率
  - `clientIp`: 客户端IP
  - `userCode`: 用户代码
  - `UserName`: 用户名
  - `realName`: 真实姓名
  - `AgentType`: 代理类型
- **Controller**: `internal/controller/client3.0_api/control_sys_api/control_sys.go`

### 13. 获取台站管理信息
- **路径**: `GET /api/Resource/StationManager`
- **说明**: 获取台站联系人信息
- **参数**: 
  - `StationId` (必填): 台站ID
- **示例**: `/api/Resource/StationManager?StationId=0101`
- **Controller**: `internal/controller/client3.0_api/get_station_manager_api/get_station_manager.go`

---

## 📝 使用说明

### 添加新路由

1. 在对应的Controller中实现接口
2. 在 `internal/router/router.go` 中注册路由
3. **在此文件中添加路由信息**（重要！）

### 格式示例

```markdown
### N. 接口名称
- **路径**: `GET/POST/PUT/DELETE /api/路径`
- **说明**: 接口功能说明
- **参数**: 
  - `param1` (必填/可选): 参数说明
- **示例**: `/api/path?param1=value1`
- **Controller**: `internal/controller/路径/文件.go`
```

---

## 🔍 快速查找

- **按功能分类**: Basic、Resource
- **按HTTP方法**: GET、POST、PUT、DELETE
- **按路径前缀**: `/Basic/`、`/Resource/`

---

**最后更新**: 2025-01-XX
**维护者**: 开发团队

