package main

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
	"time"
)

// var name string
// var url string
var start_err uint8 = 0

// App struct
type App struct {
	ctx context.Context
}

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

var connect *websocket.Conn

// GreetMany 连接客户端
func (a *App) GreetMany(name string, url string) uint8 {
	// 如果已经存在连接，先关闭它
	if connect != nil {
		connect.Close()
		connect = nil
	}

	cha_Login := make(chan string, 10)
	var err error
	connect, err = connect_server(url, name, cha_Login)

	if err != nil {
		log.Printf("连接服务器失败: %v\n", err)
		start_err = 1
		return start_err
	}

	if connect != nil {
		go get_message(connect, cha_Login)
		go func(appCtx context.Context) {
			log.Println("📡 通道接收协程已启动，等待数据...")
			for data := range cha_Login {
				log.Printf("📥 从通道收到：%s\n", data)
				// 核心：发射事件 "chat:update"，数据是 JSON 字符串
				runtime.EventsEmit(appCtx, "chat:update", data)
			}
		}(a.ctx)
	}

	start_err = 0
	return start_err
}

// SeedMany
func (a *App) SeedMany(name string, mes string) {
	if connect == nil {
		log.Println("未建立连接，无法发送消息")
		return
	}

	err := connect.WriteMessage(websocket.TextMessage, []byte(name+":"+mes))
	if err != nil {
		log.Printf("发送消息失败: %v\n", err)
	}
}

func connect_server(url string, name string, cha_Login chan<- string) (*websocket.Conn, error) {
	// 设置超时时间，避免无限重试
	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return nil, fmt.Errorf("连接超时")
		case <-ticker.C:
			connect, _, err := websocket.DefaultDialer.Dial((fmt.Sprintf("ws://%s:8080/ws", url)), nil)
			if err == nil {
				log.Println("连接成功")
				err = connect.WriteMessage(websocket.TextMessage, []byte(name))
				if err != nil {
					log.Printf("发送初始消息失败: %v\n", err)
					connect.Close()
					continue
				}
				return connect, nil
			} else {
				log.Printf("连接失败: %v\n", err)
				log.Println("将在3秒后重试")
			}
		}
	}
}

func get_message(connect *websocket.Conn, cha_mes chan<- string) {
	defer close(cha_mes)
	for {
		_, mes, err := connect.ReadMessage()
		if err != nil {
			log.Printf("接收消息失败: %v\n", err)
			// 连接可能已断开，退出循环
			break
		}
		cha_mes <- string(mes)
	}
}
