package main

import (
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
	"net/http"
	"time"
)

// App 总结构体，用去前端调用
// App struct
type App struct {
	ctx    context.Context // 上下文，用于控制线程
	server *http.Server    // 服务器连接存储，用于判断是否连接
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// 返还值是否开启
var start_err uint8 = 0

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

}

// 存储聊天记录的格式
type name struct {
	Time string `json:"time"`
	Date string `json:"date"`
}

// 存储用户名称
var chat_record []name

// （暂存所有连接，但是这样子是错的，会使得所有连接过的都在里面）
var connect_all []*websocket.Conn

// 项目的示例代码，这里不删掉了
//  Greet returns a greeting for the given name
// func (a *App) Greet(name string) string {
// 	return fmt.Sprintf("Hello %s, It's show time!", name)
// }

// 将http升级为websocket
var upgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
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

	if a.server != nil {
		log.Println("🔞服务已经在运行中了！不要再启动了啦！！！！（生气）")
		return 0
	}
	// 通道，用来获取客户端传入的数据，
	ch_json := make(chan string, 10)

	// 通道接收协程：拿到数据后发射事件给前端
	go func(appCtx context.Context) {
		log.Println("🔶通道接收协程已启动，等待数据...🔶")
		for data := range ch_json {
			log.Printf("📩~从通道收到：%s\n", data)
			// 核心：发射事件 "chat:update"，数据是 JSON 字符串
			runtime.EventsEmit(appCtx, "chat:update", data)
		}
	}(a.ctx) // 绑定上下文

	// 主动设置路由，避免重启后路由占用
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", processInfo(ch_json)) // 注册路由并启动

	// 启动服务
	go func() {
		log.Printf("🚀 WebSocket 服务启动：ws://localhost:8080/ws\n")
		a.server = &http.Server{
			Addr:    ":8080",
			Handler: mux, // 设置路由
		}
		err := a.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Printf("❌ WebSocket 服务启动失败：%v\n", err)
		}
	}()

	return 0
}

// 主程序入口
func processInfo(ch_json chan<- string) func(w http.ResponseWriter, r *http.Request) {
	// 因为mux.HandleFunc("/ws", processInfo(ch_json))注册路由只支持两个参数，但是我需要给通道进来，采用闭包的形式
	return func(w http.ResponseWriter, r *http.Request) {
		connect, err := upgrade.Upgrade(w, r, nil)
		// 依旧是史，只是先修复着
		connect_all = append(connect_all, connect)
		if err != nil {
			log.Println("＞﹏＜连接失败", w, r, err)
			return
		}
		defer func(connect *websocket.Conn) {
			err := connect.Close()
			if err != nil {

			}
		}(connect)
		log.Println("ヾ(≧▽≦*)o连接成功", w, r)

		for {
			_, message, err := connect.ReadMessage()
			if err != nil {
				log.Println("（⊙ｏ⊙）接收失败", err)
				break
			}
			log.Println("📩收到消息：", string(message))

			currentTime := time.Now().Format("2006-01-02 15:04:05.000")
			chat_record = append(chat_record, name{
				Time: currentTime,
				Date: string(message),
			})

			// 最后传到前端的数据格式为json
			json_chat, err := json.Marshal(chat_record)
			if err != nil {
				log.Println("❌ JSON 序列化失败：", err)
				continue
			}
			ch_json <- string(json_chat)
			for data := range connect_all {
				_ = connect_all[data].WriteMessage(websocket.TextMessage, json_chat)

			}
		}
	}
}
