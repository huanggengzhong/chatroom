package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

type Message struct {
	OwnerID int
	Content string
}

type User struct {
	ID             int
	Addr           string
	EnterAt        time.Time
	MessageChannel chan string
}

func (u *User) String() string {
	return u.Addr + ",UID:" + strconv.Itoa(u.ID) + ",进来时间:" + u.EnterAt.Format("2006-01-02 15:04:05+8000")
}

var (
	enteringChannel = make(chan *User)
	leavingChannel  = make(chan *User)
	messageChannel  = make(chan Message, 8)
	//messageChannel = make(chan string, 8)
)
var (
	globalID int
	idLocker sync.Mutex
)

func main() {
	listener, err := net.Listen("tcp", ":2020")
	if err != nil {
		panic(err)
	}
	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleConn(conn)
	}
}

// broadcaster 用于记录聊天室用户，并进行消息广播：
// 1. 新用户进来；2. 用户普通消息；3. 用户离开
func broadcaster() {
	users := make(map[*User]struct{})
	//值空结构体类型
	for {
		select {
		case user := <-enteringChannel:
			users[user] = struct{}{}
			//空结构体实例
		case user := <-leavingChannel:
			delete(users, user)
			close(user.MessageChannel)
		case msg := <-messageChannel:
			for user := range users {
				if user.ID == msg.OwnerID {
					continue
				}
				user.MessageChannel <- msg.Content
			}
		}

	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	user := &User{
		ID:             GenUserID(),
		Addr:           conn.RemoteAddr().String(),
		EnterAt:        time.Now(),
		MessageChannel: make(chan string, 8),
	}
	go sendMessage(conn, user.MessageChannel)

	user.MessageChannel <- "欢迎," + user.String()
	msg := Message{
		OwnerID: user.ID,
		Content: "用户:`" + strconv.Itoa(user.ID) + "`进来了",
	}
	messageChannel <- msg

	enteringChannel <- user

	//避免退出
	var userActive = make(chan struct{})
	go func() {
		d := 5 * time.Minute
		timer := time.NewTimer(d)
		for {
			select {
			case <-timer.C:
				conn.Close()
			case <-userActive:
				timer.Reset(d)
			}
		}
	}()

	input := bufio.NewScanner(conn)
	for input.Scan() {
		msg.Content = strconv.Itoa(user.ID) + ":" + input.Text()
		messageChannel <- msg

		//用户活跃
		userActive <- struct{}{}
	}
	if err := input.Err(); err != nil {
		log.Println("读取错误:", err)
	}

	leavingChannel <- user
	msg.Content = "用户:`" + strconv.Itoa(user.ID) + "`离开了"
	messageChannel <- msg
}

func GenUserID() int {
	idLocker.Lock()
	defer idLocker.Unlock()
	globalID++
	return globalID
}

func sendMessage(conn net.Conn, ch <-chan string) {
	//只写的通道
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}

}
