package api

import (
	"context"
	"encoding/json"
	"fmt"
	"gf_api/internal/db"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/redis/go-redis/v9"
)

// 台站数据总览 ldc 20250901

// Register 把当前模块的所有路由注册到 group
func Register(group *ghttp.RouterGroup) {
	group.GET("/Basic/OverViewData", GetOverViewData)
	group.GET("/Basic/AllStation", GetAllStaitonInfo)
	group.GET("/Basic/AllStationId", GetAllStationId)
}

// 台站总览的结构体
type OverViewDataRow struct {
	ID                 int64   `json:"id"`
	StationID          string  `json:"station_id"`
	Remarks            string  `json:"remarks"`
	RelationPositionId string  `json:"relation_position_id"`
	SetItemModelId     []g.Map `json:"setitem_model_id"` //查数据库
	OperateModelId     []g.Map `json:"operate_model_id"` //查数据库
	DynamicModelId     []g.Map `json:"dynamic_model_id"` //查数据库
	StaticModelId      []g.Map `json:"static_model_id"`  //查数据库
	ParentNodeId       int64   `json:"parent_node_id"`
	PositionId         string  `json:"position_id"`
	DeviceTypeId       string  `json:"device_type_id"`
	SubSystemId        string  `json:"sub_system_id"`
	NodeName           string  `json:"node_name"`
}

type StationIdResp struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// flattenKeys 用来指定哪些字段需要再去查 Redis
var modelKeyMapping = map[string]string{
	"dynamic_model_id": "svr_dynamic_model",
	"operate_model_id": "svr_operations_model",
	"setitem_model_id": "svr_setitem_model",
	"static_model_id":  "svr_static_model",
}

// 构造台站总览数据接口。
// 参数：stationId，查数据库台站模型表、静态属性表和动态属性表，组成台站模型结构，然后去redis中取字段的值
// 必须先取所有台站信息数据列表得到 stationId, 接口：http://127.0.0.1:8001/api/Basic/AllStationId
func GetOverViewDataOld(r *ghttp.Request) {
	ctx := context.Background()

	// Step 1. 请求 /api/Basic/AllStationId 获取 stationId 列表，只取第一个
	allURL := "http://127.0.0.1:8001/api/Basic/AllStationId" // 建议写到 config.yaml
	resp, err := g.Client().Timeout(5*time.Second).Get(ctx, allURL)
	if err != nil {
		r.Response.WriteJson(g.Map{
			"error": fmt.Sprintf("请求 AllStationId 接口失败: %v", err),
		})
		return
	}
	defer resp.Close()

	if resp.StatusCode != http.StatusOK {
		r.Response.WriteJson(g.Map{
			"error": fmt.Sprintf("AllStationId 接口返回非200: %d", resp.StatusCode),
		})
		return
	}

	// //这里假设接口返回的是数组 [{key:..., value:...}, {...}]
	// var stationList []StationIdResp
	// // 读取响应体
	// body := resp.ReadAll()
	// // 用 gjson.DecodeTo 解析 JSON
	// if err := gjson.DecodeTo(body, &stationList); err != nil {
	// 	r.Response.WriteJson(g.Map{
	// 		"error": fmt.Sprintf("解析 AllStationId 响应失败: %v", err),
	// 	})
	// 	return
	// }
	// if len(stationList) == 0 {
	// 	r.Response.WriteJson(g.Map{
	// 		"error": "未获取到任何 stationId",
	// 	})
	// 	return
	// }
	// // 这里取第一个 stationId
	// stationId := stationList[0].Value
	// if stationId == "" {
	// 	r.Response.WriteJson(g.Map{
	// 		"error": "未获取到 stationId",
	// 	})
	// 	return
	// }

	var stationList []StationIdResp
	body := resp.ReadAll()

	// 适配接口返回结果为数组或单个对象
	if err := gjson.DecodeTo(body, &stationList); err != nil || len(stationList) == 0 {
		// 如果不是数组，再解析成单个对象
		var station StationIdResp
		if err2 := gjson.DecodeTo(body, &station); err2 == nil && station.Value != "" {
			stationList = append(stationList, station)
		}
	}

	if len(stationList) == 0 {
		r.Response.WriteJson(g.Map{
			"error": "未获取到任何 stationId",
		})
		return
	}

	stationId := stationList[0].Value

	//db.InitRedis() //初始化Redis

	// Step 2. 查询数据库台站模型表station_node,关联到静态属性表static_model、动态属性表dynamic_model、设置属性表setitem_model和操作属性表operations_model，组成最终的json模型
	//var rows []g.Map
	sql := `
		SELECT node_id, parent_node_id, node_name,dynamic_model_id, static_model_id,setitem_model_id,relation_position_id,position_id
		FROM station_node
		WHERE station_id = $1
		ORDER BY node_id;
		`

	res, err := g.DB().GetAll(ctx, sql, stationId)
	if err != nil {
		fmt.Println("查询台站信息失败:", err)
	}
	// else {
	// 	g.Dump(res) // 打印结果
	// }

	if len(res) == 0 {
		r.Response.Write("{}")
		return
	}

	// 1) 定义总览数据接口的json结构体
	type nodeInfo struct {
		NodeId             int    //节点id
		parent             int    //父节点id
		name               string //节点名字
		dynamic            string //动态节点
		static             string //静态节点
		setitem            string //设置项
		relationPositionID string //关联工位号
		positionID         string //工位号
	}
	//根节点node_id，parent_node_id指向node_id
	nodes := make(map[int]*nodeInfo)   // node_id -> nodeInfo
	childrenMap := make(map[int][]int) // parent_id -> []child_id
	minNodeID := math.MaxInt32         // 找最小的 node_id，当成根

	//遍历数据库结果，填充 nodes 和 childrenMap
	for _, row := range res {
		NodeId := gconv.Int(row["node_id"])
		parent := gconv.Int(row["parent_node_id"])
		name := gconv.String(row["node_name"])
		dynamic := gconv.String(row["dynamic_model_id"])
		static := gconv.String(row["static_model_id"])
		setitem := gconv.String(row["setitem_model_id"])
		relationPositionID := gconv.String(row["relation_position_id"])
		positionID := gconv.String(row["position_id"])

		nodes[NodeId] = &nodeInfo{NodeId: NodeId, parent: parent, name: name, dynamic: dynamic, static: static, setitem: setitem, relationPositionID: relationPositionID, positionID: positionID}
		childrenMap[parent] = append(childrenMap[parent], NodeId)

		if NodeId < minNodeID {
			minNodeID = NodeId
		}
	}

	// 如果没有找到最小节点，就给一个空的根节点，并返回
	if minNodeID == math.MaxInt32 {
		r.Response.Write("{}")
		return
	}

	// 2) 以显式栈做后序遍历（避免递归），并带有环检测
	type frame struct {
		NodeId  int
		visited bool // 标记是否已经处理过子节点。false先处理子节点，true再处理自己
	}

	processed := make(map[int]bool) // 已处理好的节点
	onStack := make(map[int]bool)   //  当前路径上的节点，用于检测环
	result := make(map[int]g.Map)   // 每个 node_id 的子节点树

	stack := []frame{{NodeId: minNodeID, visited: false}}

	// 用于收集检测到的环（仅做日志）
	cycles := make([]int, 0)

	//每次从栈里取一个节点 NodeId，找到对应的 nodeInfo。
	for len(stack) > 0 {
		// pop
		f := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		NodeId := f.NodeId

		nodeInfo, ok := nodes[NodeId]

		// 若parent 指向了一个不存在的节点，即节点不存在，这种情况就生成空 map。
		if !ok {
			result[NodeId] = g.Map{}
			processed[NodeId] = true
			continue
		}

		//当第一次访问时
		if !f.visited {
			if processed[NodeId] {
				// 已构建，跳过
				continue
			}
			if onStack[NodeId] {
				// 如果再次遇到在路径上的节点，说明有环；记录并标记为空，避免死循环
				cycles = append(cycles, NodeId)
				result[NodeId] = g.Map{}
				processed[NodeId] = true
				continue
			}
			// 若还没处理过，就先标记 onStack，再把自己 二次入栈（visited=true），同时把所有子节点先入栈。这样保证 后序遍历（先处理子节点，再处理自己）。
			onStack[NodeId] = true
			stack = append(stack, frame{NodeId: NodeId, visited: true})
			// 子节点入栈
			for _, childID := range childrenMap[NodeId] {
				if !processed[childID] {
					stack = append(stack, frame{NodeId: childID, visited: false})
				}
			}
		} else { //不是第一次访问
			// 遍历子节点，把每个子节点的子树挂到 m 上
			m := g.Map{}
			for _, childID := range childrenMap[NodeId] {
				// child 可能因环或缺失未被 processed，确保 result 有值
				if _, ok := result[childID]; !ok {
					result[childID] = g.Map{}
				}
				// 如果 child 在 nodes 中不存在，则用占位名 "unknown"
				childName := "unknown"
				if n, ok := nodes[childID]; ok {
					childName = n.name
				}
				m[childName] = result[childID]
			}

			// 处理redis动态属性：再取 dynamic_model_id 对应的 redis 值
			// 查找对应的 dynamic_model_id，动态属性。
			// staton_node表中有工位号position_id和关联工位号relation_position_id，在dynamic_model_id查出来的结果中，取key作为节点字段属性放到父节点下，而值是根据下面的规则取值：
			// 1.每个节点属性下面都有属性名panro和关联属性名relation_parno；
			// 2.设变量Num:当关联工位号relation_position_id和关联属性relation_parno都有值，则使用Num=relation_position_id；否则使用Num=position_id。
			// 3.当is_enable=1时，表示有设备，则Num=position_id
			// 4.去redis中查key=svr_DATA_Num的值，作为节点属性值。
			if nodeInfo.dynamic != "" {
				// 获取指定 key=svr_dynamic_model 下的所有字段和值
				allDynCmd := db.Redis.HGetAll(ctx, "svr_dynamic_model")
				allDyn, err := allDynCmd.Result()
				if err != nil {
					fmt.Println(ctx, "获取 redis svr_dynamic_model 失败:", err)
				} else {
					// 找到当前 dynamic_model_id 对应的值
					dynVal, ok := allDyn[nodeInfo.dynamic]
					if ok && dynVal != "" {
						// 解析成 map[string]any
						var dyn map[string]any
						if err := json.Unmarshal([]byte(dynVal), &dyn); err != nil {
							// 解析失败 → 挂原始字符串
							fmt.Println("dynamic json unmarshal failed:", nodeInfo.dynamic, "err:", err)
							m[nodeInfo.name] = dynVal
						} else {
							// 遍历 dynamic_model_id 下的每个属性（key:属性名，value:属性值）
							for attrKey, attrVal := range dyn {
								if attrObj, ok := attrVal.(map[string]any); ok {
									// 提取 para_name、relation_parno、is_enable
									//parno, _ := attrObj["parno"].(string)
									relationParno, _ := attrObj["relation_parno"].(string)
									isEnable := "0"                        //默认
									if v, ok := attrObj["is_enable"]; ok { //is_enable这个字段不一定在redis中存在
										switch vv := v.(type) {
										case string:
											isEnable = vv
										case float64: // JSON 里数字默认是 float64
											if vv == 1 {
												isEnable = "1"
											}
										}
									}
									// fmt.Println("parno:", parno)
									// fmt.Println("relation_parno:", relationParno)
									// fmt.Println("is_enable:", isEnable)

									// =============== 规则逻辑开始 ===============
									// 设变量 Num = position_id，默认情况
									num := nodeInfo.positionID

									// 如果 relation_position_id 和 relation_parno 都有值 → Num = relation_position_id
									if nodeInfo.relationPositionID != "" && relationParno != "" {
										num = nodeInfo.relationPositionID
									}

									// 如果 is_enable = "1" → 强制 Num = position_id
									if isEnable == "1" {
										num = nodeInfo.positionID
									}
									// =============== 规则逻辑结束 ===============

									// 拼接 redis的key = svr_DATA_Num
									redisKey := fmt.Sprintf("svr_DATA_%s", num)
									valCmd := db.Redis.Get(ctx, redisKey)
									val, err := valCmd.Result()
									if err != nil {
										//fmt.Println("获取 redis key 失败:", redisKey, "err:", err)
										// 兜底：把 para_value 挂上去
										pv := ""
										if v, ok := attrObj["para_value"]; ok {
											switch vv := v.(type) {
											case string:
												pv = vv
											case float64: // 如果 JSON 里是数字
												pv = fmt.Sprintf("%v", vv)
											default:
												// 其他类型，转成字符串
												pv = fmt.Sprintf("%v", vv)
											}
										}
										m[attrKey] = pv
									} else {
										// 正常挂 redis 里的值
										m[attrKey] = val
									}

									// 调试输出
									// fmt.Printf("节点:%s 属性:%s Num:%s redisKey:%s => %v\n",
									// 	nodeInfo.name, paraName, num, redisKey, m[attrKey])
								}
							}
						}
					} else {
						fmt.Println("查询不到对应的 dynamic_model_id:", nodeInfo.dynamic)
					}
				}
			}

			// 处理静态属性
			if nodeInfo.static != "" {
				allSticCmd := db.Redis.HGetAll(ctx, "svr_static_model")
				allStic, _ := allSticCmd.Result()
				sticVal, ok := allStic[nodeInfo.static] //去结果中查key的值
				if ok && sticVal != "" {
					var stic map[string]any
					if err := json.Unmarshal([]byte(sticVal), &stic); err != nil {
						// 解析失败 → 只打日志，不污染结果
						fmt.Println("static json unmarshal failed:", nodeInfo.static, "err:", err)
					} else {
						for attrKey, attrVal := range stic {
							if attrObj, ok := attrVal.(map[string]any); ok {
								// 安全取 para_value
								pv := getStringAttr(attrObj, "para_value", "")
								m[attrKey] = pv
							} else {
								// 如果不是 map，直接转成字符串
								m[attrKey] = fmt.Sprintf("%v", attrVal)
							}
						}
					}
				} else {
					fmt.Println("查询不到对应的 static_model_id:", nodeInfo.static)
				}
			}

			//处理设置项属性，规则与动态属性表一致
			// staton_node表中有工位号position_id和关联工位号relation_position_id，在svr_setitem_model查出来的结果中，取key作为节点字段属性放到父节点下，而值是根据下面的规则取值：
			// 1.每个节点属性下面都有属性名panro和关联属性名relation_parno；
			// 2.设变量Num:当关联工位号relation_position_id和关联属性relation_parno都有值，则使用Num=relation_position_id；否则使用Num=position_id。
			// 3.当is_enable=1时，表示有设备，则Num=position_id
			// 4.去redis中查key=svr_DATA_Num的值，作为节点属性值。
			if nodeInfo.setitem != "" {
				// 获取指定 key=svr_setitem_model 下的所有字段和值
				allSetCmd := db.Redis.HGetAll(ctx, "svr_setitem_model")
				allSet, _ := allSetCmd.Result()
				if err != nil {
					fmt.Println(ctx, "获取 redis svr_setitem_model 失败:", err)
				} else {
					// 找到当前 setitem_model_id 对应的值
					setVal, ok := allSet[nodeInfo.setitem]
					if ok && setVal != "" {
						// 解析成 map[string]any
						var set map[string]any
						if err := json.Unmarshal([]byte(setVal), &set); err != nil {
							// 解析失败 → 挂原始字符串
							fmt.Println("setitem json unmarshal failed:", nodeInfo.setitem, "err:", err)
							m[nodeInfo.name] = setVal
						} else {
							// 遍历 setitem_model_id 下的每个属性（key:属性名，value:属性值）
							for attrKey, attrVal := range set {
								if attrObj, ok := attrVal.(map[string]any); ok {
									// 提取 para_name、relation_parno、is_enable
									relationParno, _ := attrObj["relation_parno"].(string)
									isEnable := "0"                        //默认
									if v, ok := attrObj["is_enable"]; ok { //is_enable这个字段不一定在redis中存在
										switch vv := v.(type) {
										case string:
											isEnable = vv
										case float64: // JSON 里数字默认是 float64
											if vv == 1 {
												isEnable = "1"
											}
										}
									}
									// fmt.Println("parno:", parno)
									// fmt.Println("relation_parno:", relationParno)
									// fmt.Println("is_enable:", isEnable)

									// =============== 规则逻辑开始 ===============
									// 设变量 Num = position_id，默认情况
									num := nodeInfo.positionID

									// 如果 relation_position_id 和 relation_parno 都有值 → Num = relation_position_id
									if nodeInfo.relationPositionID != "" && relationParno != "" {
										num = nodeInfo.relationPositionID
									}

									// 如果 is_enable = "1" → 强制 Num = position_id
									if isEnable == "1" {
										num = nodeInfo.positionID
									}
									// =============== 规则逻辑结束 ===============

									// 拼接 redis的key = svr_DATA_Num
									redisKey := fmt.Sprintf("svr_DATA_%s", num)
									valCmd := db.Redis.Get(ctx, redisKey)
									val, err := valCmd.Result()
									keyLower := strings.ToLower(attrKey) // key 转小写
									if err != nil {
										//fmt.Println("获取 redis key 失败:", redisKey, "err:", err)
										// 兜底：把 para_value 挂上去
										pv := ""
										if v, ok := attrObj["para_value"]; ok {
											switch vv := v.(type) {
											case string:
												pv = strings.ToLower(vv)
											case float64: // 如果 JSON 里是数字
												pv = fmt.Sprintf("%v", vv)
											default:
												// 其他类型，转成字符串
												pv = fmt.Sprintf("%v", vv)
											}
										}
										m[keyLower] = strings.ToLower(pv) // 挂到结果前再转小写
									} else {
										// 正常挂 redis 里的值
										// 如果不是 map，直接转成字符串并转小写
										m[keyLower] = strings.ToLower(fmt.Sprintf("%v", val))
									}

									// 调试输出
									// fmt.Printf("节点:%s 属性:%s Num:%s redisKey:%s => %v\n",
									// 	nodeInfo.name, paraName, num, redisKey, m[attrKey])
								}
							}
						}
					} else {
						fmt.Println("查询不到对应的 setitem_model_id:", nodeInfo.setitem)
					}
				}
			}

			//当前节点处理完成，保存结果并解除 onStack
			result[NodeId] = m
			processed[NodeId] = true
			onStack[NodeId] = false
		}
	}

	// 3) 最终结果：以 minNodeID 作为根（“node_id 最小为第一节点”）
	finalJSON := result[minNodeID]

	// 4) 如果检测到环，记录日志（可选：也可以把这些信息返回给前端）
	if len(cycles) > 0 {
		fmt.Println(ctx, "检测到环形引用（已自动忽略），示例节点id：", cycles)
	}

	// 5) 转 JSON 并返回
	jsonStr := gjson.New(finalJSON).MustToJsonString()
	//fmt.Println(jsonStr)
	r.Response.Write(jsonStr)
}

func GetStationIdInfo(r *ghttp.Request) {
	ctx := context.Background()

	// 请求 /api/Basic/AllStationId 获取 stationId 列表
	allURL := "http://127.0.0.1:8001/api/Basic/AllStationId" // 建议写到 config.yaml
	resp, err := g.Client().Timeout(5*time.Second).Get(ctx, allURL)
	if err != nil {
		r.Response.WriteJson(g.Map{
			"error": fmt.Sprintf("请求 AllStationId 接口失败: %v", err),
		})
		return
	}
	defer resp.Close()

	if resp.StatusCode != http.StatusOK {
		r.Response.WriteJson(g.Map{
			"error": fmt.Sprintf("AllStationId 接口返回非200: %d", resp.StatusCode),
		})
		return
	}

	var stationList []StationIdResp
	body := resp.ReadAll()

	// 适配接口返回结果为数组或单个对象
	if err := gjson.DecodeTo(body, &stationList); err != nil || len(stationList) == 0 {
		// 如果不是数组，再解析成单个对象
		var station StationIdResp
		if err2 := gjson.DecodeTo(body, &station); err2 == nil && station.Value != "" {
			stationList = append(stationList, station)
		}
	}

	if len(stationList) == 0 {
		r.Response.WriteJson(g.Map{
			"error": "未获取到任何 stationId",
		})
		return
	}

	//stationId := stationList[0].Value //只取第一个
}

// 台站总览数据接口已完成，但是速度较慢，还需要优化 ldc 20251017W
func GetOverViewData(r *ghttp.Request) {
	ctx := context.Background()
	start := time.Now() // 记录开始时间

	//db.InitRedis() //初始化Redis

	// 加载模型缓存（只读一次 Redis）
	cache, err := LoadModelCache(ctx)
	if err != nil {
		r.Response.WriteJson(g.Map{"error": fmt.Sprintf("加载模型缓存失败: %v", err)})
		return
	}

	stepStart := start // 每个阶段的起点时间

	// 从请求参数中获取 stationId，例如 /api/GetOverViewData?stationId=0101
	stationId := r.Get("stationId").String()
	if stationId == "" {
		r.Response.WriteJson(g.Map{
			"error": "缺少参数 stationId",
		})
		return
	}

	// // 1. 取 Basic
	// stepStart = time.Now()
	// basicVal, err := db.Redis.HGet(ctx, "svr_stationNodeModelBasic", stationId).Result()
	// fmt.Printf("阶段1 耗时: %v ms\n", time.Since(stepStart).Milliseconds())
	// if err == redis.Nil {
	// 	r.Response.WriteJson(g.Map{
	// 		"error": fmt.Sprintf("Redis中未找到 svr_stationNodeModelBasic,[%s]", stationId),
	// 	})
	// 	return
	// } else if err != nil {
	// 	r.Response.WriteJson(g.Map{
	// 		"error": fmt.Sprintf("读取 svr_stationNodeModelBasic 出错: %v", err),
	// 	})
	// 	return
	// }

	// stepStart = time.Now()
	// var basic map[string]interface{}
	// if err := json.Unmarshal([]byte(basicVal), &basic); err != nil {
	// 	r.Response.WriteJson(g.Map{
	// 		"error": fmt.Sprintf("解析 Basic JSON 失败: %v", err),
	// 	})
	// 	return
	// }
	// fmt.Printf("阶段2 Basic JSON解析耗时: %v ms\n", time.Since(stepStart).Milliseconds())

	// // 2. 取 Idx
	// stepStart = time.Now()
	// idxVal, _ := db.Redis.HGet(ctx, "svr_stationNodeModelIdx", stationId).Result()
	// fmt.Printf("阶段3 耗时: %v ms\n", time.Since(stepStart).Milliseconds())
	// if err == redis.Nil {
	// 	r.Response.WriteJson(g.Map{
	// 		"error": fmt.Sprintf("Redis中未找到 svr_stationNodeModelIdx,[%s]", stationId),
	// 	})
	// 	return
	// } else if err != nil {
	// 	r.Response.WriteJson(g.Map{
	// 		"error": fmt.Sprintf("读取 svr_stationNodeModelIdx 出错: %v", err),
	// 	})
	// 	return
	// }

	// stepStart = time.Now()
	// var idx map[string]interface{}
	// if err := json.Unmarshal([]byte(idxVal), &idx); err != nil {
	// 	r.Response.WriteJson(g.Map{
	// 		"error": fmt.Sprintf("解析 Idx JSON 失败: %v", err),
	// 	})
	// 	return
	// }
	// fmt.Printf("阶段4 Idx JSON解析耗时: %v ms\n", time.Since(stepStart).Milliseconds())

	// 1.读取 Basic 和 Idx（并行）
	stepStart = time.Now()
	var (
		basicStr, idxStr string
		basicErr, idxErr error
	)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		basicStr, basicErr = db.Redis.HGet(ctx, "svr_stationNodeModelBasic", stationId).Result()
	}()

	go func() {
		defer wg.Done()

		idxStr, idxErr = db.Redis.HGet(ctx, "svr_stationNodeModelIdx", stationId).Result()
	}()

	wg.Wait()

	if basicErr != nil {
		r.Response.WriteJson(g.Map{"error": fmt.Sprintf("读取 svr_stationNodeModelBasic 出错: %v", basicErr)})
		return
	}
	if idxErr != nil {
		r.Response.WriteJson(g.Map{"error": fmt.Sprintf("读取 svr_stationNodeModelIdx 出错: %v", idxErr)})
		return
	}
	fmt.Printf("阶段1 耗时: %v ms\n", time.Since(stepStart).Milliseconds())

	// 2.解析 JSON（并行）
	stepStart = time.Now()
	var basic, idx map[string]interface{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		_ = json.Unmarshal([]byte(basicStr), &basic)
	}()
	go func() {
		defer wg.Done()

		_ = json.Unmarshal([]byte(idxStr), &idx)
	}()
	wg.Wait()
	fmt.Printf("阶段2 JSON解析耗时: %v ms\n", time.Since(stepStart).Milliseconds())

	// // 3.扫描 idx，找出所有 positionId / rPositionId
	// stepStart = time.Now()
	// nums := collectAllPositions(idx)
	// fmt.Printf("阶段3 提取 positionId/rPositionId 数量: %d 耗时: %v ms\n", len(nums), time.Since(stepStart).Milliseconds())

	// // 4.批量预加载 svr_DATA_* 到内存
	// stepStart = time.Now()
	// dataCache, err := preloadDataByNums(ctx, nums)
	// if err != nil {
	// 	fmt.Printf("preloadDataByNums 出错: %v\n", err)
	// }
	// fmt.Printf("阶段4 批量加载 Redis 数据耗时: %v ms\n", time.Since(stepStart).Milliseconds())

	stepStart = time.Now()
	//  5.合并
	//mergeRecursiveCache(ctx, basic, idx, cache, dataCache)
	mergeRecursive(basic, idx, cache) //使用这个函数，就得定义cache
	//mergeRecursiveOld(ctx, basic, idx) //使用这个函数无需定义cache
	fmt.Printf("阶段5 mergeRecursive 耗时: %v ms\n", time.Since(stepStart).Milliseconds())

	// 计算耗时
	duration := time.Since(start) // 计算执行时长
	durationMs := duration.Milliseconds()

	// 输出 JSON 给前端前，包装一下
	response := g.Map{
		"stationId":   stationId,
		"MachineName": "",
		"Result":      "true",
		"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
		"time":        durationMs,
		"Message":     "",
		"Content":     basic, // 这是原本的合并结果
	}

	r.Response.WriteJson(response)
	//r.Response.WriteJson(json.RawMessage(data))
}

// 扫描 idx 中的所有 positionId / rPositionId
func collectAllPositions(node map[string]interface{}) []string {
	numSet := make(map[string]struct{})
	var collect func(map[string]interface{})
	collect = func(m map[string]interface{}) {
		if v, ok := m["positionId"].(string); ok && v != "" {
			numSet[v] = struct{}{}
		}
		if v, ok := m["rPositionId"].(string); ok && v != "" {
			numSet[v] = struct{}{}
		}

		for _, v := range m {
			if sub, ok := v.(map[string]interface{}); ok {
				collect(sub)
			}
		}
	}
	collect(node)

	nums := make([]string, 0, len(numSet))
	for k := range numSet {
		nums = append(nums, k)
	}
	return nums
}

// 批量预加载 Redis 中所有 svr_DATA_*（按 positionId / rPositionId）
func preloadDataByNums(ctx context.Context, nums []string) (map[string]map[string]string, error) {
	dataCache := make(map[string]map[string]string)
	if len(nums) == 0 {
		return dataCache, nil
	}

	pipe := db.Redis.Pipeline()
	cmds := make(map[string]*redis.MapStringStringCmd) // ✅ 改这里

	for _, num := range nums {
		key := fmt.Sprintf("svr_DATA_%s", num)
		cmds[key] = pipe.HGetAll(ctx, key)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	for key, cmd := range cmds {
		dataCache[key] = cmd.Val()
	}

	return dataCache, nil
}

// mergeRecursive 合并 Basic 和 Idx，并展开 rConfig
func mergeRecursiveOld(ctx context.Context, basic, idx map[string]interface{}) {
	// 处理 dynamic_model_id
	if id, ok := idx["dynamic_model_id"].(string); ok && id != "" {
		processDynamicModel(ctx, idx, id)
		basic["dynamic_model_id"] = id
	}

	// 处理 svr_static_model
	if id2, ok := idx["static_model_id"].(string); ok && id2 != "" {
		processStaticModel(ctx, idx, id2)
		basic["static_model_id"] = id2
	}

	// 处理 setitem_model_id
	if id3, ok := idx["setitem_model_id"].(string); ok && id3 != "" {
		processSetitemModel(ctx, idx, "setitem_model_id", id3)
		basic["setitem_model_id"] = id3
	}

	for k, v := range idx {
		// 跳过 dynamic_model_id
		if k == "dynamic_model_id" {
			continue
		}

		// 跳过 svr_static_model
		if k == "static_model_id" {
			continue
		}

		if k == "setitem_model_id" {
			// 跳过 setitem_model_id
			continue
		}

		// 特殊处理 rConfig：跳过节点名，仅合并其子属性
		if k == "rConfig" {
			if sub, ok := v.(map[string]interface{}); ok {
				mergeRecursiveOld(ctx, basic, sub)
			}
			continue
		}

		switch sub := v.(type) {
		case map[string]interface{}:
			// 若 basic 中无该键则创建
			if _, exists := basic[k]; !exists {
				basic[k] = make(map[string]interface{})
			}
			mergeRecursiveOld(ctx, basic[k].(map[string]interface{}), sub)
		default:
			basic[k] = v
		}
	}
}

// mergeRecursive 合并数据总览的基本属性Basic 和 数据总览的详细属性Idx。这个函数是从缓存中取，而不是每次都去查redis，速度快点。20251014 ldc
func mergeRecursive(basic, idx map[string]interface{}, cache *ModelCache) {
	subStart := time.Now()
	if id, ok := idx["dynamic_model_id"].(string); ok && id != "" {
		processDynamicModelCached(idx, id, cache)
		basic["dynamic_model_id"] = id

	}

	if id, ok := idx["static_model_id"].(string); ok && id != "" {
		processStaticModelCached(idx, id, cache)
		basic["static_model_id"] = id

	}

	if id, ok := idx["setitem_model_id"].(string); ok && id != "" {
		processSetitemModelCached(idx, id, cache)
		basic["setitem_model_id"] = id

	}

	for k, v := range idx {
		switch val := v.(type) {
		case map[string]interface{}:
			if k == "rConfig" {
				mergeRecursive(basic, val, cache)
				continue
			}
			if _, exists := basic[k]; !exists {
				basic[k] = make(map[string]interface{})
			}
			mergeRecursive(basic[k].(map[string]interface{}), val, cache)
		default:
			if k == "dynamic_model_id" || k == "static_model_id" || k == "setitem_model_id" {
				continue
			}
			basic[k] = v
		}
	}

	// for k, v := range idx {
	// 	if k == "dynamic_model_id" || k == "static_model_id" || k == "setitem_model_id" {
	// 		continue
	// 	}
	// 	if k == "rConfig" {
	// 		if sub, ok := v.(map[string]interface{}); ok {
	// 			// 递归调用 mergeRecursive，把 rConfig 下的内容合并到 basic
	// 			mergeRecursive(ctx, basic, sub, cache)
	// 		}
	// 		continue
	// 	}
	// 	switch sub := v.(type) {
	// 	case map[string]interface{}: //如果 value 是 map（非 rConfig），递归合并到 basic 对应的 key 下。
	// 		if _, exists := basic[k]; !exists {
	// 			basic[k] = make(map[string]interface{})
	// 		}
	// 		mergeRecursive(ctx, basic[k].(map[string]interface{}), sub, cache)
	// 	default: // 普通值（string、int 等）
	// 		basic[k] = v // 直接赋值到 basic
	// 	}
	// }

	fmt.Printf("mergeRecursive总耗时: %v\n", time.Since(subStart))
}

// 从svr_Data缓存集中取值
func mergeRecursiveCache(ctx context.Context, basic, idx map[string]interface{}, cache *ModelCache, dataCache map[string]map[string]string) {
	subStart := time.Now()
	if id, ok := idx["dynamic_model_id"].(string); ok && id != "" {
		processDynamicModelCachedNew(ctx, idx, id, cache, dataCache)
		basic["dynamic_model_id"] = id
	}

	if id, ok := idx["static_model_id"].(string); ok && id != "" {
		processStaticModelCached(idx, id, cache)
		basic["static_model_id"] = id
	}

	if id, ok := idx["setitem_model_id"].(string); ok && id != "" {
		processSetitemModelCached(idx, id, cache)
		basic["setitem_model_id"] = id
	}

	for k, v := range idx {
		if k == "dynamic_model_id" || k == "static_model_id" || k == "setitem_model_id" {
			continue
		}

		if k == "rConfig" {
			if sub, ok := v.(map[string]interface{}); ok {
				mergeRecursive(basic, sub, cache)
			}
			continue
		}

		switch sub := v.(type) {
		case map[string]interface{}:
			if _, exists := basic[k]; !exists {
				basic[k] = make(map[string]interface{})
			}
			mergeRecursive(basic[k].(map[string]interface{}), sub, cache)
		default:
			basic[k] = v
		}
	}
	fmt.Printf("mergeRecursiveCache 总耗时: %v\n", time.Since(subStart))
}

// 预加载模型，动态、静态和设置项 20251014 ldc
type ModelCache struct {
	Dynamic map[string]map[string]map[string]interface{}
	Static  map[string]map[string]interface{}
	SetItem map[string]map[string]map[string]interface{}
}

// 加载redis模型缓存,以便减少访问redis的次数 20251014 ldc
func LoadModelCache(ctx context.Context) (*ModelCache, error) {
	cache := &ModelCache{
		Dynamic: make(map[string]map[string]map[string]interface{}),
		Static:  make(map[string]map[string]interface{}),
		SetItem: make(map[string]map[string]map[string]interface{}),
	}

	// 动态模型的结果集
	if all, err := db.Redis.HGetAll(ctx, "svr_dynamic_model").Result(); err == nil {
		for id, jsonStr := range all {
			var obj map[string]map[string]interface{}
			if err := json.Unmarshal([]byte(jsonStr), &obj); err == nil {
				cache.Dynamic[id] = obj
			}
		}
	}

	// modelID := "DTVTrans1000W_1_DynamicModel"
	// // 取出这个模型的定义（结果是一个 map）
	// modelDef, ok := cache.Dynamic[modelID]
	// if !ok {
	// 	fmt.Println("找不到 Dynamic 模型:", modelID)
	// }
	// // 打印这个模型下的所有属性定义
	// for attr, def := range modelDef {
	// 	fmt.Printf("属性名: %s, 定义: %v\n", attr, def)
	// }

	// 静态模型的结果集
	if all, err := db.Redis.HGetAll(ctx, "svr_static_model").Result(); err == nil {
		for id, jsonStr := range all {
			var obj map[string]interface{}
			if err := json.Unmarshal([]byte(jsonStr), &obj); err == nil {
				cache.Static[id] = obj
			}
		}
	}

	// 设置模型的结果集
	if all, err := db.Redis.HGetAll(ctx, "svr_setitem_model").Result(); err == nil {
		for id, jsonStr := range all {
			var obj map[string]map[string]interface{}
			if err := json.Unmarshal([]byte(jsonStr), &obj); err == nil {
				cache.SetItem[id] = obj
			}
		}
	}

	return cache, nil
}

// 从缓存中取动态属性  20251014 ldc
/***
查dynamic_model_id的值的时候，DTVSystemSwitcher_1属性中有工位号positionId和关联工位号rpositionId，
在dynamic_model_id查出来的结果中，取key作为节点字段属性放到父节点下，而值是根据下面的规则取值：
1.每个节点属性下面都有属性名panro和关联属性名relation_parno；
2.设变量Num：当关联工位号rpositionId和关联属性relation_parno都有值，则Num=rpositionId；否则Num=positionId。
3.当is_enable=1时，表示有设备，则Num=positionId；
4.去redis中查key=svr_DATA_Num的结果集。
5.使用的是rpositionId得到的结果集，则取relation_parno作为key去结果集中查到具体值；使用positionId得到的结果集，则取parno作为key去结果集中查到具体值。
***/
func processDynamicModelCached(node map[string]interface{}, modelID string, cache *ModelCache) {
	positionId, _ := node["positionId"].(string)
	rPositionId, _ := node["rPositionId"].(string)

	// 从缓存中根据动态属性id取出动态模型的结果集
	modelDef, ok := cache.Dynamic[modelID]
	if !ok {
		return
	}

	// 遍历结果集中的所有属性定义
	for attrName, attrDef := range modelDef {
		//attrName = strings.ToLower(attrName)

		relationParno, _ := attrDef["relation_parno"].(string)
		parno, _ := attrDef["parno"].(string)
		isEnable, _ := attrDef["is_enable"].(float64)

		// 决定使用哪个工位号 Num
		Num := positionId
		useRelation := false

		if rPositionId != "" && relationParno != "" {
			Num = rPositionId
			useRelation = true
		}
		if isEnable == 1 {
			Num = positionId
			useRelation = false
		}

		// 查询 Redis
		dataKey := fmt.Sprintf("svr_DATA_%s", Num)
		dataVal, err := db.Redis.HGetAll(context.Background(), dataKey).Result()
		if err != nil || len(dataVal) == 0 {
			continue
		}

		// 全部 key 转为小写，动态属性先暂时不需要转小写
		// lowerMap := make(map[string]string, len(dataVal))
		// for k, v := range dataVal {
		// 	lowerMap[strings.ToLower(k)] = v
		// }

		// 根据是否是关联取不同的字段值
		var finalValue string
		if useRelation {
			//finalValue = lowerMap[strings.ToLower(relationParno)]
			finalValue = dataVal[relationParno]
		} else {
			//finalValue = lowerMap[strings.ToLower(parno)]
			finalValue = dataVal[parno]
		}

		// 设置节点属性值，即使为空
		node[attrName] = finalValue
	}
}

// 改为从svr_Data缓存集中取数
func processDynamicModelCachedNew(ctx context.Context, node map[string]interface{}, modelID string, cache *ModelCache, dataCache map[string]map[string]string) {
	positionId, _ := node["positionId"].(string)
	rPositionId, _ := node["rPositionId"].(string)

	// 从缓存中根据动态属性id取出动态模型的结果集
	modelDef, ok := cache.Dynamic[modelID]
	if !ok {
		return
	}

	// 遍历结果集中的所有属性定义
	for attrName, attrDef := range modelDef {
		//attrName = strings.ToLower(attrName)

		relationParno, _ := attrDef["relation_parno"].(string)
		parno, _ := attrDef["parno"].(string)
		isEnable, _ := attrDef["is_enable"].(float64)

		// 决定使用哪个工位号 Num
		Num := positionId
		useRelation := false

		if rPositionId != "" && relationParno != "" {
			Num = rPositionId
			useRelation = true
		}
		if isEnable == 1 {
			Num = positionId
			useRelation = false
		}

		// 查询 Redis
		dataKey := fmt.Sprintf("svr_DATA_%s", Num)
		dataVal, ok := dataCache[dataKey]
		if !ok || len(dataVal) == 0 {
			continue
		}

		// 根据是否是关联取不同的字段值
		var finalValue string
		if useRelation {
			//finalValue = lowerMap[strings.ToLower(relationParno)]
			finalValue = dataVal[relationParno]
		} else {
			//finalValue = lowerMap[strings.ToLower(parno)]
			finalValue = dataVal[parno]
		}

		// 设置节点属性值，即使为空
		node[attrName] = finalValue
	}
}

// 从缓存中取静态属性 20251014 ldc
func processStaticModelCached(node map[string]interface{}, modelID string, cache *ModelCache) {
	// 从缓存中获取静态模型定义
	modelDef, ok := cache.Static[modelID]
	if !ok {
		return
	}

	for attrName, attrDef := range modelDef {
		attrName = strings.ToLower(attrName)

		// 如果属性定义是 map，则取 para_value
		if attrMap, ok := attrDef.(map[string]interface{}); ok {
			if paraValue, ok := attrMap["para_value"]; ok {
				node[attrName] = paraValue
			}
			continue
		}

		// 挂载到当前节点
		node[attrName] = attrDef
	}
}

// 从缓存中取设置项属性 20251014 ldc
/***
属性中有工位号positionId和关联工位号rpositionId，
在Setitem_model_id查出来的结果中，取key作为节点字段属性放到父节点下，而值是根据下面的规则取值：
1.每个节点属性下面都有属性名panro和关联属性名relation_parno；
2.设变量Num：当关联工位号rpositionId和关联属性relation_parno都有值，则Num=rpositionId；否则Num=positionId。
3.当is_enable=1时，表示有设备，则Num=positionId；
4.去redis中查key=svr_DATA_Num的结果集。
5.使用的是rpositionId得到的结果集，则取relation_parno作为key去结果集中查到具体值；使用positionId得到的结果集，则取parno作为key去结果集中查到具体值。
***/
func processSetitemModelCached(node map[string]interface{}, modelID string, cache *ModelCache) {
	positionId, _ := node["positionId"].(string)
	rPositionId, _ := node["rPositionId"].(string)

	modelDef, ok := cache.SetItem[modelID]
	if !ok {
		return
	}

	for attrName, attrDef := range modelDef {
		attrName = strings.ToLower(attrName)

		// 取配置字段
		parno, _ := attrDef["parno"].(string)
		relationParno, _ := attrDef["relation_parno"].(string)
		isEnable, _ := attrDef["is_enable"].(float64)

		// ----------- 第2步：确定 Num ----------
		Num := positionId
		useR := false
		if rPositionId != "" && relationParno != "" {
			Num = rPositionId
			useR = true
		}
		if isEnable == 1 {
			Num = positionId
			useR = false
		}

		// ----------- 第3步：从 Redis 获取数据 ----------
		dataKey := fmt.Sprintf("svr_DATA_%s", Num)
		dataVal, err := db.Redis.HGetAll(context.Background(), dataKey).Result()
		if err != nil || len(dataVal) == 0 {
			continue
		}

		// ----------- 第4步：全小写处理 ----------
		lowerMap := make(map[string]string)
		for k, v := range dataVal {
			lowerMap[strings.ToLower(k)] = v
		}

		// ----------- 第5步：根据使用哪种 Num 取值 ----------
		var val string
		if useR && relationParno != "" {
			val = lowerMap[strings.ToLower(relationParno)]
		} else if parno != "" {
			val = lowerMap[strings.ToLower(parno)]
		}

		// ----------- 第6步：赋值到节点 ----------
		node[attrName] = val
	}
}

// 处理动态属性

func processDynamicModel(ctx context.Context, node map[string]interface{}, modelID string) {
	// positionId 和 rPositionId 由 Idx 中挂过来
	positionId, _ := node["positionId"].(string)
	rPositionId, _ := node["rPositionId"].(string)

	// 查 Redis: svr_dynamic_model
	val, err := db.Redis.HGetAll(ctx, "svr_dynamic_model").Result()
	if err != nil {
		fmt.Println("Redis HGetAll svr_dynamic_model error:", err)
		return
	}

	// 通过 modelID 获取该模型定义的 JSON
	modelJSON, ok := val[modelID]
	if !ok || modelJSON == "" {
		//fmt.Println("modelID not found in svr_dynamic_model:", modelID)
		return
	}

	// 解析模型定义 JSON
	var modelDef map[string]map[string]interface{}
	if err := json.Unmarshal([]byte(modelJSON), &modelDef); err != nil {
		fmt.Println("Unmarshal modelDef error:", err)
		return
	}

	// 遍历属性定义
	for attrName, attrDef := range modelDef {
		//全部转为小写
		attrName = strings.ToLower(attrName)

		relationParno, _ := attrDef["relation_parno"].(string)
		isEnable, _ := attrDef["is_enable"].(float64)

		// 规则：决定 Num
		Num := positionId
		if rPositionId != "" && relationParno != "" {
			Num = rPositionId
		}
		if isEnable == 1 {
			Num = positionId
		}

		// 去 Redis 查数据值
		dataKey := fmt.Sprintf("svr_DATA_%s", Num)
		dataVal, err := db.Redis.HGetAll(ctx, dataKey).Result()
		if err != nil {
			fmt.Println("Redis HGetAll", dataKey, "error:", err)
			continue
		}
		// 将数据的值挂载到对应节点上
		node[attrName] = dataVal
	}
}

// 处理静态属性,查redis中的svr_static_model的到结果集，再根据static_model_id查出具体的属性值，得到的是一个json字符串，取paran_value作为该属性值的值
func processStaticModel(ctx context.Context, node map[string]interface{}, modelID string) {
	// 1. 从 Redis 读取所有静态模型定义
	val, err := db.Redis.HGetAll(ctx, "svr_static_model").Result()
	if err != nil {
		fmt.Println("Redis HGetAll svr_static_model error:", err)
		return
	}

	// 2. 根据 modelID 获取具体模型定义
	modelJSON, ok := val[modelID]
	if !ok || modelJSON == "" {
		fmt.Println("modelID not found in svr_static_model:", modelID)
		return
	}

	// 3. 解析 JSON
	var modelDef map[string]interface{}
	if err := json.Unmarshal([]byte(modelJSON), &modelDef); err != nil {
		fmt.Println("Unmarshal svr_static_model JSON error:", err)
		return
	}

	// 4. 遍历属性定义
	for attrName, attrDef := range modelDef {
		//全部转为小写
		attrName = strings.ToLower(attrName)

		// attrDef 可能是一个对象，需要判断类型
		if attrMap, ok := attrDef.(map[string]interface{}); ok {
			if paraValue, ok := attrMap["para_value"]; ok {
				node[attrName] = paraValue
			}
		}
	}
}

// 处理设置项属性，与动态属性的规则逻辑一样。
func processSetitemModel(ctx context.Context, node map[string]interface{}, modelType, modelID string) {
	// positionId 和 rPositionId 由 Idx 中挂过来
	positionId, _ := node["positionId"].(string)
	rPositionId, _ := node["rPositionId"].(string)

	// 查 Redis: svr_setitem_model
	val, err := db.Redis.HGetAll(ctx, "svr_setitem_model").Result()
	if err != nil {
		fmt.Println("Redis HGetAll svr_setitem_model error:", err)
		return
	}

	// 通过 modelID 获取该模型定义的 JSON
	modelJSON, ok := val[modelID]
	if !ok || modelJSON == "" {
		//fmt.Println("modelID not found in svr_setitem_model:", modelID)
		return
	}

	// 解析模型定义 JSON
	var modelDef map[string]map[string]interface{}
	if err := json.Unmarshal([]byte(modelJSON), &modelDef); err != nil {
		fmt.Println("Unmarshal modelDef error:", err)
		return
	}

	// 遍历属性定义
	for attrName, attrDef := range modelDef {
		//全部转为小写
		attrName = strings.ToLower(attrName)

		relationParno, _ := attrDef["relation_parno"].(string)
		isEnable, _ := attrDef["is_enable"].(float64)

		// 规则：决定 Num
		Num := positionId
		if rPositionId != "" && relationParno != "" {
			Num = rPositionId
		}
		if isEnable == 1 {
			Num = positionId
		}

		// 去 Redis 查数据值
		dataKey := fmt.Sprintf("svr_DATA_%s", Num)
		dataVal, err := db.Redis.HGetAll(ctx, dataKey).Result()
		if err != nil {
			fmt.Println("Redis HGetAll", dataKey, "error:", err)
			continue
		}
		// 将数据的值挂载到对应节点上
		node[attrName] = dataVal

	}
}

// 处理 operate的逻辑
func processNormalModel(ctx context.Context, node map[string]interface{}, modelType, modelID string) {
	var redisKey string
	switch modelType {
	case "static_model_id":
		redisKey = "svr_static_model"
	case "operate_model_id":
		redisKey = "svr_operate_model"
	case "setitem_model_id":
		redisKey = "svr_setitem_model"
	default:
		return
	}

	val := db.Redis.Get(ctx, redisKey)

	var modelDef map[string]interface{}
	if err := json.Unmarshal([]byte(val.String()), &modelDef); err != nil {
		return
	}

	for attrName, attrVal := range modelDef {
		node[attrName] = attrVal
	}
}

//处理其他属性

// getStringAttr 安全地从 map[string]any 里取 string 值
func getStringAttr(obj map[string]any, key, def string) string {
	if v, ok := obj[key]; ok {
		switch vv := v.(type) {
		case string:
			return vv
		case float64: // JSON number 默认会被解成 float64
			return fmt.Sprintf("%v", vv)
		default:
			return fmt.Sprintf("%v", vv)
		}
	}
	return def
}

// 获取所有台站信息接口
func GetAllStaitonInfo(r *ghttp.Request) {
	ctx := context.Background()
	key := "svr_stations"

	val, err := db.Redis.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil { // 注意这里redis.Nil,它是 github.com/redis/go-redis/v9 包里定义
			// key 不存在
			r.Response.WriteJson(g.Map{
				"error": fmt.Sprintf("Redis key '%s' 不存在", key),
			})
			return
		}
		// 其他错误
		r.Response.WriteJson(g.Map{
			"error": err.Error(),
		})
		return
	}

	// 成功获取值
	r.Response.WriteJson(g.Map{
		"key":   key,
		"value": val,
	})
}

// 获取所有台站Id接口
func GetAllStationId(r *ghttp.Request) {
	ctx := context.Background()
	key := "svr_station_id"

	val, err := db.Redis.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil { // 注意这里redis.Nil,它是 github.com/redis/go-redis/v9 包里定义
			// key 不存在
			r.Response.WriteJson(g.Map{
				"error": fmt.Sprintf("Redis key '%s' 不存在", key),
			})
			return
		}
		// 其他错误
		r.Response.WriteJson(g.Map{
			"error": err.Error(),
		})
		return
	}

	// 成功获取值
	r.Response.WriteJson(g.Map{
		"key":   key,
		"value": val,
	})
}
