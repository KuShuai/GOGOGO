package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct{
	Ip string
	Port int

	//在线用户列表
	OnlineMap map[string]*User
	mapLock sync.RWMutex

	//消息广播的channel
	Message chan string
}

//创建server的接口
func NewServer(ip string, port int) *Server{
	server := &Server{
		Ip : ip,
		Port: port,
		OnlineMap:make(map[string]*User),
		Message:make(chan string),
	}

	return server
}
//监听Message广播消息的channel的goroutine
func (this *Server)ListenMessager()  {
	for  {
		msg := <-this.Message

		this.mapLock.Lock()
		for _,cli := range this.OnlineMap{
			cli.C <-msg
		}
		this.mapLock.Unlock()
	}
}

func (this *Server)BroadCast(user *User,msg string)  {
	sendMsg := "["+user.Addr+"]"+user.Name+":"+msg
	this.Message<-sendMsg
}

func (this *Server)Handler(conn net.Conn)  {
	//当前链接业务
	fmt.Println("链接建立成功")
	//当前用户上线了
	//将用户加入到onlineMap中，并广播当前用户上线消息
	user:=NewUser(conn,this)
	//上线
	user.Online()

	//接收客户端发送的消息
	go func() {
		buf := make([]byte,4096)
		for  {
			n,err := conn.Read(buf)
			if n == 0 {
				user.OffLine();
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("conn Read err:",err)
				return
			}
			//提取用户的消息
			msg := string(buf[:n])
			user.DoMessage(msg)
		}
	}()

	//当前handler阻塞
	select {

	}
}


//启动服务器
func (this *Server)Start()  {
	//Listen
	listener,err := net.Listen("tcp",fmt.Sprintf(("%s:%d"),this.Ip,this.Port))
	if err != nil {
		fmt.Println("net.Listen err",err)
		return
	}

	//启动监听message的goroutine
	go this.ListenMessager()

	for {
		//assept
		conn,err := listener.Accept()
		fmt.Println("Accept")
		if err != nil {
			fmt.Println("listener accept err",err)
			continue
		}
		//do handler
		go this.Handler(conn)
	}
	//close listen socket
	defer listener.Close()
}