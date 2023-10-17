package logic

import (
	"github.com/huanggengzhong/chatroom/global"
	"log"
)

// 广播器
type broadcaster struct {
	users map[string]*User //所有聊天用户
	//所有channel统一在这里
	enteringChannel chan *User
	leaveingChannel chan *User
	messageChannel  chan *Message

	checkUserChannel      chan string  //用户名
	checkUserCanInChannel chan bool    // 判断该昵称用户是否可进入聊天室（重复与否）：true 能，false 不能
	usersChannel          chan []*User //用户列表
}

var Broadcaster = &broadcaster{
	users:                 make(map[string]*User),
	enteringChannel:       make(chan *User),
	leaveingChannel:       make(chan *User),
	messageChannel:        make(chan *Message, global.MessageQueueLen),
	checkUserChannel:      make(chan string),
	checkUserCanInChannel: make(chan bool),
}

// Start 启动广播器
func (b *broadcaster) Start() {
	for {
		select {
		case user := <-b.enteringChannel:
			b.users[user.NickName] = user //加入users队列
			//todo离线发送
		case nickname := <-b.checkUserChannel:
			if _, ok := b.users[nickname]; ok {
				b.checkUserCanInChannel <- false
			} else {
				b.checkUserCanInChannel <- true
			}
		case user := <-b.leaveingChannel:
			delete(b.users, user.NickName)
			user.CloseMessageChannel()
		case msg := <-b.messageChannel:
			//给所有除本人在线用户发消息
			for _, user := range b.users {
				if user.UID == msg.User.UID {
					continue
				}
				user.MessageChannel <- msg
			}
			//保存离线信息
			OfflineProcessor.Save(msg)
		}
	}
}

func (b *broadcaster) CanEnterRoom(nickname string) bool {
	b.checkUserChannel <- nickname
	return <-b.checkUserCanInChannel
}

func (b *broadcaster) Broadcase(msg *Message) {
	//判断缓存通道溢出
	if len(b.messageChannel) >= global.MessageQueueLen {
		log.Println("通道满了")
	}
	b.messageChannel <- msg
}

func (b *broadcaster) UserEntering(user *User) {
	b.enteringChannel <- user
}

func (b *broadcaster) UserLeaving(user *User) {
	b.leaveingChannel <- user
}
