package server

import (
	"fmt"
	"github.com/huanggengzhong/chatroom/logic"
	"log"
	"net/http"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
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

	// 1. 新用户进来，构建该用户的实例

	nickname := req.FormValue("nickname")
	if len := len(nickname); len < 2 || len > 20 {
		log.Println("非法昵称,昵称长度要求:2-20", nickname)
		wsjson.Write(req.Context(), conn, logic.NewErrorMessage("非法昵称,昵称长度要求:2-20"))
		conn.Close(websocket.StatusUnsupportedData, "非法昵称,昵称长度要求:2-20")
		return
	}
	if !logic.Broadcaster.CanEnterRoom(nickname) {
		log.Println("昵称已存在,请换个昵称", nickname)
		wsjson.Write(req.Context(), conn, logic.NewErrorMessage("昵称已存在,请换个昵称"))
		conn.Close(websocket.StatusUnsupportedData, "昵称已存在,请换个昵称")
		return

	}
	token := req.FormValue("token")
	userHasToken := logic.NewUser(conn, nickname, token, req.RemoteAddr)
	fmt.Printf("用户信息:%v", userHasToken)
	// 2. 开启给用户发送消息的 goroutine
	go userHasToken.SendMessage(req.Context())
	// 3. 给当前用户发送欢迎信息
	userHasToken.MessageChannel <- logic.NewWelcomeMessage(nickname)
	//防止token泄露
	tmpUser := *userHasToken
	user := &tmpUser
	user.Token = ""
	//其它用户告知到来
	msg := logic.NewUserEnterMessage(user)
	logic.Broadcaster.Broadcase(msg)
	// 4. 将该用户加入广播器的用列表中
	logic.Broadcaster.UserEntering(user)
	// 5. 接收用户消息
	err = user.ReceiveMessage(req.Context())
	// 6. 用户离开
	logic.Broadcaster.UserLeaving(user)
	msg = logic.NewUserLeaveMessage(user)
	logic.Broadcaster.Broadcase(msg)

	if err == nil {
		conn.Close(websocket.StatusNormalClosure, "")
	} else {
		log.Println("客户端读取失败", err)
		conn.Close(websocket.StatusInternalError, "客户端读取失败")
	}

}
