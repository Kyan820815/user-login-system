package tcp

import (
	"entry_task/mysqldb"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var cc *grpc.ClientConn

func makeClientConn(port int32) (*grpc.ClientConn, error) {
	port = 0
	dialOptions := []grpc.DialOption{grpc.WithInsecure()}
	// address := fmt.Sprintf("localhost:%d", port)
	return grpc.Dial("localhost:10000", dialOptions...)
}

func ClientConn() (TCPRPCClient, error) {
	if cc == nil {
		var err error
		cc, err = makeClientConn(0)
		if err != nil {
			fmt.Println("[TCP Client ClientConn] Error to connect to TCP server: ", err)
			return nil, err
		}
	}
	return NewTCPRPCClient(cc), nil
}

func SayhelloRPC() (bool, error) {
	cc, err := ClientConn()
	if err != nil {
		return false, err
	}

	helloMsg := &HelloMsg{
		Greeting: "hi from client",
	}

	rsp, err := cc.HelloCaller(context.Background(), helloMsg)
	if err != nil {
		return false, err
	}

	return rsp.Ok, nil
}

func LoginRPC(user *mysqldb.User) (*mysqldb.User, error) {
	cc, err := ClientConn()

	if err != nil {
		return nil, err
	}
	userMsg := loadUserMsg(user)

	rsp, err := cc.LoginCaller(context.Background(), userMsg)
	if err != nil {
		return nil, err
	}
	user = loadUser(rsp)
	return user, nil
}

func NicknameRPC(user *mysqldb.User) error {
	cc, err := ClientConn()
	if err != nil {
		return err
	}
	userMsg := loadUserMsg(user)

	_, err = cc.NicknameCaller(context.Background(), userMsg)

	return err
}

func PhotoRPC(user *mysqldb.User) error {
	cc, err := ClientConn()
	if err != nil {
		return err
	}

	userMsg := loadUserMsg(user)

	_, err = cc.PhotoCaller(context.Background(), userMsg)

	return err
}
