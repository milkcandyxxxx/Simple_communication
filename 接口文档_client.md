Stop()关闭连接
Start()启动服务器

// 通道接收协程：拿到数据后发射事件给前端
go func(appCtx context.Context) {
log.Println("🔶通道接收协程已启动，等待数据...🔶")
for data := range ch_json {
log.Printf("📩~从通道收到：%s\n", data)
// 核心：发射事件 "chat:update"，数据是 JSON 字符串
runtime.EventsEmit(appCtx, "chat:update", data)
}