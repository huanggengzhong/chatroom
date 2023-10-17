package server

import (
	"fmt"
	"github.com/huanggengzhong/chatroom/logic"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

var rootDir string

func RegisterHandle() {
	inferRootDir()
	//广播消息处理
	go logic.Broadcaster.Start()

	http.HandleFunc("/", homeHandleFunc)
	http.HandleFunc("/ws", WebSocketHandleFunc)

}

func inferRootDir() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	var infer func(d string) string
	infer = func(d string) string {
		// 这里要确保项目根目录下存在 template 目录
		if exists(d + "/template") {
			return d
		}
		//filepath.Dir 是 Go 语言标准库中的一个函数，它的作用是从一个文件路径中提取出该路径所在的目录部分。具体来说，它会返回一个字符串，表示给定路径的目录部分，不包括文件名。

		return infer(filepath.Dir(d))
	}
	rootDir = infer(cwd)
	//fmt.Println(rootDir, "rootDir目录")

}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	//Stat返回一个描述name指定的文件对象的FileInfo。如果指定的文件对象是一个符号链接，返回的FileInfo描述该符号链接指向的文件的信息，本函数会尝试跳转该链接。如果出错，返回的错误值为*PathError类型。
	//IsExist返回一个布尔值说明该错误是否表示一个文件或目录已经存在
	return err == nil || os.IsExist(err)
}

func homeHandleFunc(w http.ResponseWriter, req *http.Request) {

	tpl, err := template.ParseFiles(rootDir + "/template/home.html")
	if err != nil {
		fmt.Fprint(w, "html模版解析错误")
		return
	}
	err = tpl.Execute(w, nil)
	if err != nil {
		fmt.Fprint(w, "html模版执行错误")
		return
	}
}
