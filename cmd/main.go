package main

import (
	"fmt"
	"qqbot/internal/config"
	"qqbot/internal/handler"
	_ "qqbot/internal/plug"
	"qqbot/internal/service"

	newlog "github.com/Lqingzhe/go-common-tool/log"
)

//	func main() {
//		// 1. 连接 NapCat 的正向 WS 服务端
//		conn, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:3001?access_token=token", nil)
//		if err != nil {
//			log.Fatal("连接失败:", err)
//		}
//		defer conn.Close()
//
//		// 2. 发消息的辅助函数
//		send := func(action string, params map[string]interface{}) {
//			req := map[string]interface{}{
//				"action": action,
//				"params": params,
//				"echo":   "req-1", // 用于匹配响应
//			}
//			data, _ := json.Marshal(req)
//			conn.WriteMessage(websocket.TextMessage, data)
//		}
//
//		// 3. 开一个协程持续读消息
//		go func() {
//			for {
//				messagetype, msg, err := conn.ReadMessage()
//				if err != nil {
//					log.Println("读取结束:", err)
//					return
//				}
//
//				// 先解析出 post_type，判断是事件还是 API 响应
//				var raw map[string]any
//				json.Unmarshal(msg, &raw)
//				fmt.Println("messagetype:", messagetype)
//				fmt.Println("raw message1: ", raw)
//				var resp model.OneBotResp
//				json.Unmarshal(msg, &resp)
//				fmt.Println("resp2: ", resp)
//				// API 响应（含 echo 字段）
//				if echo, ok := raw["echo"].(string); ok {
//					fmt.Printf("[API响应]3 echo=%s status=%v\n", echo, raw["status"])
//					continue
//				}
//
//				// 事件推送（含 post_type 字段）
//				if postType, ok := raw["post_type"].(string); ok {
//					fmt.Printf("[事件]4 post_type=%s\n", postType)
//
//					// 如果是私聊消息，回复一句话
//					if postType == "message" {
//						if msgType, _ := raw["message_type"].(string); msgType == "private" {
//							userID := int64(raw["user_id"].(float64))
//							fmt.Printf("  5收到私聊消息，来自 %d，内容: %v\n", userID, raw["raw_message"])
//							send("send_private_msg", map[string]any{
//								"user_id": userID,
//								"message": "Hello from Go!",
//							})
//						}
//					}
//				}
//			}
//		}()
//
//		// 4. 主动调用一次 API 测试（获取登录信息）
//		send("get_login_info", map[string]any{})
//
//		// 5. 阻塞主协程
//		fmt.Scanln()
//		fmt.Println("退出")
//	}
func main() {
	newlog.SetIgnoreModel()

	Config := config.GetConfig()
	logger := newlog.InitLog("QQbot", Config.Global.DeviceID)
	ServiceOperator := service.GetAndInitOperator(Config)
	ServiceScheduleTask := service.GetAndInitTaskHub(Config)
	Service := service.NewService(Config, ServiceOperator, ServiceScheduleTask)
	newlog.SetLabel("通过配置文件配置access_token字段")
	Handler, err := handler.NewWebsocket("ws://"+Config.Handler.Websocket.Host+":"+Config.Handler.Websocket.Port+"?access_token="+Config.Handler.Websocket.Token, logger, Service)
	if err != nil {
		logger.LogInitFatal(err, "handler create fail")
	}
	Service.Run(Config)
	Handler.Run()

	fmt.Scanln() //回车触发优雅退出
	将配置文件写入持久化
}
