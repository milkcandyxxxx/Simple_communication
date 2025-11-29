package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
	"net/http"
	"time"
)

// App struct
type App struct {
	ctx    context.Context
	server *http.Server
}

var start_err uint8 = 0

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
func (a *App) Stop() uint8 {
	if a.server != nil {
		log.Println("🛑 正在关闭 WebSocket 服务...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := a.server.Shutdown(ctx)
		if err != nil {
			log.Printf("❌ 关闭服务失败：%v\n", err)
			return 1
		}
		a.server = nil
		log.Println("✅ WebSocket 服务已关闭")
		return 0
	}
	log.Println("⚠️ 服务未启动")
	return 0
}
func (a *App) Start() uint8 {
	// 防止重复启动
	if a.server != nil {
		log.Println("⚠️ 服务已经在运行中")
		return 0
	}

	ch_json := make(chan string, 10)

	// 通道接收协程：拿到数据后发射事件给前端
	go func(appCtx context.Context) {
		log.Println("📡 通道接收协程已启动，等待数据...")
		for data := range ch_json {
			log.Printf("📥 从通道收到：%s\n", data)
			// 核心：发射事件 "chat:update"，数据是 JSON 字符串
			runtime.EventsEmit(appCtx, "chat:update", data)
		}
	}(a.ctx)

	// 创建新的 ServeMux 避免路由冲突
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", processInfo(ch_json))

	// 启动 HTTP 服务（用协程启动，避免阻塞桌面主线程）
	go func() {
		log.Printf("🚀 WebSocket 服务启动：ws://localhost:8080/ws\n")
		a.server = &http.Server{
			Addr:    ":8080",
			Handler: mux,
		}
		err := a.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Printf("❌ WebSocket 服务启动失败：%v\n", err)
		}
	}()

	return 0
}

type name struct {
	Time string `json:"time"`
	Date string `json:"date"`
}

var chat_record []name
var upgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func processInfo(ch_json chan<- string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		connect, err := upgrade.Upgrade(w, r, nil)
		if err != nil {
			log.Println("❌连接失败", w, r, err)
			return
		}
		defer func(connect *websocket.Conn) {
			err := connect.Close()
			if err != nil {

			}
		}(connect)
		log.Println("✔️连接成功", w, r)

		for {
			_, message, err := connect.ReadMessage()
			if err != nil {
				log.Println("❗接收失败", err)
				break
			}
			log.Println("📩 收到客户端消息：", string(message))

			currentTime := time.Now().Format("2006-01-02 15:04:05.000")
			chat_record = append(chat_record, name{
				Time: currentTime,
				Date: string(message),
			})

			// 序列化 JSON
			json_chat, err := json.Marshal(chat_record)
			if err != nil {
				log.Println("❌ JSON 序列化失败：", err)
				continue
			}

			// 非阻塞发送到通道
			select {
			case ch_json <- string(json_chat):
			default:
				log.Println("⚠️ 通道繁忙，丢弃本次数据")
			}

			// 回复客户端
			_ = connect.WriteMessage(websocket.TextMessage, json_chat)
		}

	}
}
