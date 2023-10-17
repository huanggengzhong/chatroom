package logic

import "time"

// 消息类型常量
const (
	MsgTypeNormal  = iota //普通用户消息
	MsgTypeWelcome        //欢迎消息
	MsgTypeEnter          //进入消息
	MsgTypeLeave          //离开消息
	MsgTypeError          //错误消息
)

type Message struct {
	// 哪个用户发送的消息
	User    *User     `json:"user"`
	Type    int       `json:"type"`
	Content string    `json:"content"`
	MsgTime time.Time `json:"msg_time"`
	// 发送时间
	ClientSendTime time.Time `json:"client_send_time"`
	// 消息 @ 了谁
	Ats []string `json:"ats"`
}

func NewErrorMessage(msg string) *Message {
	return &Message{
		User:    System,
		Type:    MsgTypeError,
		Content: msg,
		MsgTime: time.Now(),
	}
}

func NewWelcomeMessage(nickname string) *Message {
	return &Message{
		User:    System,
		Type:    MsgTypeWelcome,
		Content: nickname + ",您好，欢迎加入聊天室！",
		MsgTime: time.Now(),
	}
}

func NewUserEnterMessage(user *User) *Message {
	return &Message{
		User:    user,
		Type:    MsgTypeEnter,
		Content: user.NickName + "加入了聊天室",
		MsgTime: time.Now(),
	}
}
func NewUserLeaveMessage(user *User) *Message {
	return &Message{
		User:    user,
		Type:    MsgTypeLeave,
		Content: user.NickName + "离开了聊天室",
		MsgTime: time.Now(),
	}
}

func NewMessage(user *User, content, clientTime string) *Message {
	message := &Message{
		User:    user,
		Type:    MsgTypeNormal,
		Content: content,
		MsgTime: time.Now(),
	}
	if clientTime != "" {
		//todo
	}
	return message
}
