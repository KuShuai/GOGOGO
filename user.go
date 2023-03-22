package main

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C chan string
	conn net.Conn

	server *Server
}

//创建一个用户API
func  NewUser(conn net.Conn,server *Server) *User  {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
		server:server,
	}

	//启动监听当前user channel消息的goroutine
	go user.ListenMessage()
	return user
}

//用户上线功能
func (this *User )Online()  {
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()
	this.server.BroadCast(this,"已上线")
}

//用户下线功能
func (this *User)OffLine()  {
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	delete(this.server.OnlineMap,this.Name)
	this.server.mapLock.Unlock()
	this.server.BroadCast(this,"已下线")
}

//给当前user发送消息
func (this *User) SendMsg(msg string){
	this.conn.Write([]byte(msg))
}

//用户其他业务功能
func (this *User) DoMessage(msg string)  {
	if  len(msg)>9 && msg[:9] == "register|"{
		info := strings.Split(msg,"|")
		if len(info) == 3 {
			Insert(info[1],info[2],this.Addr)
			this.Name = info[1]
		}
	}else if msg == "who" {
		//查询当前在线用户有哪些
		this.server.mapLock.Lock()
		for _,user := range this.server.OnlineMap{
			onlineMsg :="["+user.Addr+"]"+user.Name+":在线..."
			this.SendMsg(onlineMsg)
		}
		this.server.mapLock.Unlock()
	}else if len(msg)>7 && msg[:7] == "rename|"{
		newName := strings.Split(msg,"|")[1]
		result := "修改成功"
		_,ok := this.server.OnlineMap[newName]
		if  ok{
			result = "用户名已存在！"
		}else{
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap,this.Name)
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()

			this.Name = newName
		}
		this.C <- result
	}else{
		this.server.BroadCast(this,msg)
	}
}

//监听当前User channel的方法，一旦有消息，直接发送给对端客户端
func (this *User)ListenMessage()  {
	for {
		msg:= <-this.C
		this.conn.Write([]byte(msg+"\n"))
	}
}
