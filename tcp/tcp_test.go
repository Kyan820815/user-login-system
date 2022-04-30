package tcp

import (
	"entry_task/myredis"
	"entry_task/mysqldb"
	"testing"
	"time"
)

const (
	username  = "root"
	password  = "12345678"
	hostname  = "localhost:3306"
	dbname    = "entry_task_user_db"
	tablename = "entry_task_tcp_test_table"
	tcp_port  = 10000
)

func iniDBSettings() error {
	// initial setting for database
	dbInfo := mysqldb.GetInfo(username, password, hostname, dbname, tablename)
	db, err := mysqldb.DbConnection(dbInfo)
	if err != nil {
		return err
	}
	defer db.Conn.Close()

	err = mysqldb.CreateUserTable(db)
	if err != nil {
		return err
	}

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

	err = mysqldb.Insert(db, user1)
	if err != nil {
		return err
	}

	err = mysqldb.Insert(db, user2)
	if err != nil {
		return err
	}

	return nil
}

func startSeverAndClientConn() (*Node, error) {
	err := iniDBSettings()
	if err != nil {
		return nil, err
	}
	server, err := StartTCPServer(username, password, hostname, dbname, tablename, tcp_port)
	if err != nil {
		return nil, err
	}
	_, err = ClientConn()
	if err != nil {
		return nil, err
	}
	return server, nil
}

func TestClientConn(t *testing.T) {
	server, err := startSeverAndClientConn()
	if err != nil {
		t.Errorf("[TestClientConn]: %s\n", err)
	}

	rsp, err := SayhelloRPC()
	if err != nil || rsp != true {
		t.Errorf("[TestClientConn]: %s\n", err)
	}

	// stop the server and go to error handling
	server.Stop()

	rsp, err = SayhelloRPC()
	if err == nil || rsp == true {
		t.Errorf("[TestClientConn]: client should get error while server is off\n")
	}
}

func TestLoginRPC(t *testing.T) {
	server, err := startSeverAndClientConn()
	if err != nil {
		t.Errorf("[TestLoginRPC]: %s\n", err)
	}
	time.Sleep(time.Second)

	user := &mysqldb.User{
		Acc:      "123",
		Pwd:      "123",
		Nickname: "",
		Photo:    "photos/init.jpeg",
	}

	// test db access
	err = myredis.Flush(server.rds)
	if err != nil {
		t.Errorf("[TestLoginRPC]: %s\n", err)
	}
	_, err = LoginRPC(user)
	if err != nil {
		t.Errorf("[TestLoginRPC]: %s\n", err)
	}

	// test cache access
	_, err = LoginRPC(user)
	if err != nil {
		t.Errorf("[TestLoginRPC]: %s\n", err)
	}

	// delete the table
	err = mysqldb.DeleteTable(server.db)
	if err != nil {
		t.Errorf("[TestLoginRPC] Error %s when deleting tables", err)
	}

	server.Stop()
}

func TestNicknameRPC(t *testing.T) {
	server, err := startSeverAndClientConn()
	if err != nil {
		t.Errorf("[TestNicknameRPC]: %s\n", err)
	}
	time.Sleep(time.Second)

	user := &mysqldb.User{
		Acc:      "123",
		Pwd:      "123",
		Nickname: "123",
		Photo:    "photos/init.jpeg",
	}

	err = NicknameRPC(user)
	if err != nil {
		t.Errorf("[TestNicknameRPC]: %s\n", err)
	}

	// delete the table
	err = mysqldb.DeleteTable(server.db)
	if err != nil {
		t.Errorf("[TestNicknameRPC] Error %s when deleting tables", err)
	}

	server.Stop()
}

func TestPhotoRPC(t *testing.T) {
	server, err := startSeverAndClientConn()
	if err != nil {
		t.Errorf("[TestPhotoRPC]: %s\n", err)
	}
	time.Sleep(time.Second)

	user := &mysqldb.User{
		Acc:      "123",
		Pwd:      "123",
		Nickname: "",
		Photo:    "photos/new.jpeg",
	}

	err = PhotoRPC(user)
	if err != nil {
		t.Errorf("[TestPhotoRPC]: %s\n", err)
	}

	// delete the table
	err = mysqldb.DeleteTable(server.db)
	if err != nil {
		t.Errorf("[TestPhotoRPC] Error %s when deleting tables", err)
	}

	server.Stop()
}
