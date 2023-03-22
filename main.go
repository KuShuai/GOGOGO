package main

import "fmt"

func main()  {
	err :=OpenSql(nil)
	if err != nil {
		fmt.Sprintln("MySQL Qpen is Error:",err)
	}else {
		fmt.Println("数据库连接成功")
	}

	server := NewServer("192.168.3.125",8888)
	server.Start()

}
