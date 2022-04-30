package tcp

import (
	"entry_task/myredis"
	"entry_task/mysqldb"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
	"sync"
)

var memoLock = &sync.RWMutex{}

type Node struct {
	UnimplementedTCPRPCServer
	server *grpc.Server
	db     *mysqldb.DB
	rds    *myredis.RDS
}

func loadUser(userMsg *UserMsg) *mysqldb.User {
	return &mysqldb.User{
		Acc:      userMsg.Acc,
		Pwd:      userMsg.Pwd,
		Nickname: userMsg.Nickname,
		Photo:    userMsg.Photo,
		Id:       userMsg.Userid,
	}
}

func loadUserMsg(user *mysqldb.User) *UserMsg {
	return &UserMsg{
		Acc:      user.Acc,
		Pwd:      user.Pwd,
		Nickname: user.Nickname,
		Photo:    user.Photo,
		Userid:   user.Id,
	}
}

func newServer(db *mysqldb.DB, rds *myredis.RDS) *Node {
	serverOptions := []grpc.ServerOption{}

	node := new(Node)
	node.server = grpc.NewServer(serverOptions...)
	node.db = db
	node.rds = rds

	return node
}

func Start(db *mysqldb.DB, rds *myredis.RDS, port int32) (*Node, error) {
	address := fmt.Sprintf("localhost:%d", port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("[TCP Server Start] Server cannot listen: ", err)
		return nil, err
	}

	node := newServer(db, rds)
	RegisterTCPRPCServer(node.server, node)
	go node.server.Serve(lis)

	return node, nil
}

func (node *Node) Stop() {
	node.server.Stop()
	node.db.Conn.Close()
	node.rds.Conn.Close()
}

func StartTCPServer(username string, password string, hostname string, dbname string, tablename string, port int32) (*Node, error) {
	dbInfo := mysqldb.GetInfo(username, password, hostname, dbname, tablename)
	db, err := mysqldb.DbConnection(dbInfo)
	if err != nil {
		fmt.Println("[TCP Server] Error when getting db connection: ", err)
		return nil, err
	}

	fmt.Println("[TCP Server] Successfully connected to database")

	err = mysqldb.CreateUserTable(db)
	if err != nil {
		fmt.Println("[TCP Server] Error when creating the table: ", err)
		return nil, err
	}

	rds := myredis.RdsConnection()
	err = myredis.Flush(rds)
	if err != nil {
		return nil, err
	}
	fmt.Println("[TCP Server] Successfully connected to redis")

	node, err := Start(db, rds, port)
	if err != nil {
		return nil, err
	}

	return node, nil
}

func (node *Node) HelloCaller(ctx context.Context, helloMsg *HelloMsg) (*OK, error) {
	fmt.Println("[TCP Server HelloCaller] ", helloMsg.Greeting)
	return &OK{Ok: true}, nil
}

func (node *Node) LoginCaller(ctx context.Context, userMsg *UserMsg) (*UserMsg, error) {
	user := loadUser(userMsg)

	// check redis first
	rds_key := user.Acc + "_" + user.Pwd
	var user_cache *mysqldb.User
	memoLock.Lock()
	err := myredis.Get(node.rds, rds_key, &user_cache)
	memoLock.Unlock()
	if user_cache == nil {
		if err != nil {
			return nil, err
		} else {
			// not in the redis, search db
			// fmt.Println("[TCP Server LoginCaller] Not in cache, search db")
			memoLock.Lock()
			user, err = mysqldb.Select(node.db, user)
			memoLock.Unlock()
			if err != nil {
				// fmt.Println("[TCP Server LoginCaller] Not in db")
				return nil, err
			}
			// update redis
			// fmt.Println("[TCP Server LoginCaller] Load user from db to cache")
			memoLock.Lock()
			err := myredis.Set(node.rds, rds_key, user)
			memoLock.Unlock()
			if err != nil {
				return nil, err
			}
		}
	} else {
		// fmt.Println("[TCP Server LoginCaller] Load from cache")
		user = user_cache
	}
	return loadUserMsg(user), nil
}

func (node *Node) NicknameCaller(ctx context.Context, userMsg *UserMsg) (*OK, error) {
	user := loadUser(userMsg)

	// update redis
	rds_key := user.Acc + "_" + user.Pwd
	err := myredis.Set(node.rds, rds_key, user)
	if err != nil {
		return nil, err
	}
	// update db
	err = mysqldb.UpdateNickname(node.db, user)
	if err != nil {
		return nil, err
	}

	return &OK{Ok: true}, nil
}

func (node *Node) PhotoCaller(ctx context.Context, userMsg *UserMsg) (*OK, error) {
	user := loadUser(userMsg)

	// update redis
	rds_key := user.Acc + "_" + user.Pwd
	err := myredis.Set(node.rds, rds_key, user)
	if err != nil {
		return nil, err
	}
	// update db
	err = mysqldb.UpdatePhoto(node.db, user)
	if err != nil {
		return nil, err
	}

	return &OK{Ok: true}, nil
}
