package main

import (
	"entry_task/mysqldb"
	"entry_task/tcp"
	"entry_task/web"
	"log"
	// "time"
)

const (
	username  = "root"
	password  = "12345678"
	hostname  = "localhost:3306"
	dbname    = "entry_task_user_db"
	tablename = "user_table"
	tcp_port  = 10000
)

func main() {
	iniDBSettings()
	tcp.StartTCPServer(username, password, hostname, dbname, tablename, tcp_port)
	web.StartWEBServer()
}

func iniDBSettings() {
	// initial setting for database
	dbInfo := mysqldb.GetInfo(username, password, hostname, dbname, tablename)
	db, err := mysqldb.DbConnection(dbInfo)
	if err != nil {
		log.Printf("[main] Error %s when getting db connection", err)
		return
	}
	defer db.Conn.Close()

	mysqldb.CreateUserTable(db)

	user1 := &mysqldb.User{
		Acc:      "123",
		Pwd:      "123",
		Nickname: "",
		Photo:    "photos/init.jpeg",
		Id:       -1,
	}
	user2 := &mysqldb.User{
		Acc:      "789",
		Pwd:      "789",
		Nickname: "",
		Photo:    "photos/init.jpeg",
		Id:       -1,
	}
	mysqldb.Insert(db, user1)
	mysqldb.Insert(db, user2)
}
