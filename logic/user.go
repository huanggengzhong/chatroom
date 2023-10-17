package logic

import (
	"context"
	"errors"
	"io"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"regexp"

	"time"
)

type User struct {
	UID            int           `json:"uid"`
	NickName       string        `json:"nickname"`
	EnterAt        time.Time     `json:"enter_at"`
	Addr           string        `json:"addr"`
	MessageChannel chan *Message `json:"-"`
	Token          string        `json:"token"`
	conn           *websocket.Conn
	isNew          bool
}

// 系统用户，代表是系统主动发送的消息
var System = &User{}

func NewUser(conn *websocket.Conn, nickname, token, addr string) *User {
	user := &User{
		NickName:       nickname,
		Addr:           addr,
		EnterAt:        time.Now(),
		MessageChannel: make(chan *Message, 32),
		Token:          token,
		conn:           conn,
	}
	return user
}

func (u *User) SendMessage(ctx context.Context) {
	for msg := range u.MessageChannel {
		wsjson.Write(ctx, u.conn, msg)
	}
}

func (u *User) ReceiveMessage(ctx context.Context) error {

	var (
		receiveMsg map[string]string
		err        error
	)
	for {
		err = wsjson.Read(ctx, u.conn, &receiveMsg)
		if err != nil {
			// 判定连接是否关闭了，正常关闭，不认为是错误
			var closeErr websocket.CloseError
			if errors.As(err, &closeErr) {
				return nil
			} else if errors.Is(err, io.EOF) {
				return nil
			}
			//errors.As() //是否errors.As 用于检查 err 是否实现了 error 接口
			//errors.Is() //是否相同错误
			return err
		}

		sendMsg := NewMessage(u, receiveMsg["content"], receiveMsg["send_time"])
		//todo 内容敏感词过滤
		//sendMsg.Content=

		//@谁
		reg := regexp.MustCompile(`@[^\s@]{2,20}`)
		sendMsg.Ats = reg.FindAllString(sendMsg.Content, -1)
		//广播消息
		Broadcaster.Broadcase(sendMsg)
	}

}

func (u *User) CloseMessageChannel() {
	close(u.MessageChannel)
}
