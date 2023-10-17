package server

import (
	"fmt"
	"log"
	"net/http"
	"nhooyr.io/websocket"
	//"nhooyr.io/websocket/wsjson"
)

func WebSocketHandleFunc(w http.ResponseWriter, req *http.Request) {
	// Accept 从客户端接受 WebSocket 握手，并将连接升级到 WebSocket。
	// 如果 Origin 域与主机不同，Accept 将拒绝握手，除非设置了 InsecureSkipVerify 选项（通过第三个参数 AcceptOptions 设置）。
	// 换句话说，默认情况下，它不允许跨源请求。如果发生错误，Accept 将始终写入适当的响应
	conn, err := websocket.Accept(w, req, nil)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("执行WebSocketHandleFunc")
	if err == nil {
		conn.Close(websocket.StatusNormalClosure, "")
	} else {
		log.Println("客户端读取失败", err)
		conn.Close(websocket.StatusInternalError, "客户端读取失败")
	}

	// 1. 新用户进来，构建该用户的实例

	// 2. 开启给用户发送消息的 goroutine

	// 3. 给当前用户发送欢迎信息

	// 4. 将该用户加入广播器的用列表中

	// 5. 接收用户消息

	// 6. 用户离开
}
