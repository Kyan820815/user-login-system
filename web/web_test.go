package web

import (
	"entry_task/mysqldb"
	"entry_task/tcp"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"

	"bytes"
	"io"
	"mime/multipart"
	"os"
)

const (
	username    = "root"
	password    = "12345678"
	hostname    = "localhost:3306"
	dbname      = "entry_task_user_db"
	tablename   = "entry_task_web_test_table"
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	tcp_port    = 10000
)

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func createRandomUsers(db *mysqldb.DB, users_cnt int) ([]*mysqldb.User, error) {
	var users []*mysqldb.User
	for i := 0; i < users_cnt; i++ {
		username := randStringBytes(10)
		password := randStringBytes(10)
		user := &mysqldb.User{
			Acc:      username,
			Pwd:      password,
			Nickname: "",
			Photo:    "",
			Id:       -1,
		}
		users = append(users, user)
		err := mysqldb.Insert(db, user)
		if err != nil {
			return nil, err
		}
	}
	return users, nil
}

func loginRequest(user *mysqldb.User) (int, error) {
	err_code := -1
	userData := url.Values{
		"username": {user.Acc},
		"password": {user.Pwd},
	}
	rsp, err := http.PostForm("http://localhost:8080/login/", userData)
	if err != nil {
		return err_code, err
	}
	defer rsp.Body.Close()

	return rsp.StatusCode, err
}

func updatePhotoRequest(user *mysqldb.User, photo_tag string) (int, error) {
	err_code := -1
	b, w, err := createMultipartFormData(photo_tag, "../photos/init.jpeg")
	if err != nil {
		return err_code, err
	}

	photo_url := fmt.Sprintf("http://localhost:8080/photo/%s", user.Acc)
	req, err := http.NewRequest("POST", photo_url, &b)
	if err != nil {
		return err_code, err
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	client := &http.Client{}
	rsp, err := client.Do(req)
	if err != nil {
		return err_code, err
	}
	defer rsp.Body.Close()

	return rsp.StatusCode, err
}

func updateNicknameRequest(user *mysqldb.User) (int, error) {
	err_code := -1
	user.Nickname = "test_nickname"
	nickname_url := fmt.Sprintf("http://localhost:8080/nickname/%s", user.Acc)
	userData := url.Values{
		"nickname": {user.Nickname},
	}

	rsp, err := http.PostForm(nickname_url, userData)
	if err != nil {
		return err_code, err
	}
	defer rsp.Body.Close()

	return rsp.StatusCode, err
}

func mustOpen(f string) (*os.File, error) {
	r, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func createMultipartFormData(fieldName, fileName string) (bytes.Buffer, *multipart.Writer, error) {
	var b bytes.Buffer
	var err error
	w := multipart.NewWriter(&b)
	var fw io.Writer
	file, err := mustOpen(fileName)
	if err != nil {
		return bytes.Buffer{}, &multipart.Writer{}, err
	}
	if fw, err = w.CreateFormFile(fieldName, file.Name()); err != nil {
		return bytes.Buffer{}, &multipart.Writer{}, err
	}
	if _, err = io.Copy(fw, file); err != nil {
		return bytes.Buffer{}, &multipart.Writer{}, err
	}
	w.Close()
	return b, w, nil
}

func createConnAndTable() (*mysqldb.DB, error) {
	dbInfo := mysqldb.GetInfo(username, password, hostname, dbname, tablename)
	db, err := mysqldb.DbConnection(dbInfo)
	if err != nil {
		return nil, err
	}

	err = mysqldb.CreateUserTable(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func TestConcurrentLoginRequest(t *testing.T) {
	db, err := createConnAndTable()
	if err != nil {
		t.Errorf("[TestConcurrentLoginRequest] Error %s when getting db connection and creating a table", err)
	}
	defer db.Conn.Close()

	// create 1000 login requests with 200 users
	users_cnt := 200
	users, err := createRandomUsers(db, users_cnt)
	if err != nil {
		t.Errorf("[TestConcurrentLoginRequest] Error %s when inserting a user into the table", err)
	}

	time.Sleep(4 * time.Second)
	server, _ := tcp.StartTCPServer(username, password, hostname, dbname, tablename, tcp_port)
	time.Sleep(time.Second)
	go StartWEBServer()

	group_sz := 40
	times := users_cnt / group_sz

	var wg sync.WaitGroup
	start := time.Now()
	for i := 0; i < times; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, idx int) {
			defer wg.Done()
			for j := 0; j < group_sz; j++ {
				user := users[idx*group_sz+j]
				userData := url.Values{
					"username": {user.Acc},
					"password": {user.Pwd},
				}
				same_user_request := 1000 / users_cnt
				for k := 0; k < same_user_request; k++ {
					rsp, err := http.PostForm("http://localhost:8080/login/", userData)
					if err != nil {
						t.Errorf("[TestConcurrentLoginRequest] Error %s when send login request to the http server", err)
					}

					if rsp.StatusCode != http.StatusOK {
						t.Errorf("[TestConcurrentLoginRequest] handler returned wrong status code: got %v want %v", rsp.StatusCode, http.StatusOK)
					}
					rsp.Body.Close()
				}
			}

		}(&wg, i)
	}
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("[TestConcurrentLoginRequest] took %s\n", elapsed)

	// delete table
	err = mysqldb.DeleteTable(db)
	if err != nil {
		t.Errorf("[TestConcurrentLoginRequest] %s", err)
	}

	server.Stop()
}

func TestNickNameHandler(t *testing.T) {
	db, err := createConnAndTable()
	if err != nil {
		t.Errorf("[TestNickNameHandler] Error %s when getting db connection and creating a table", err)
	}
	defer db.Conn.Close()

	users_cnt := 2
	users, err := createRandomUsers(db, users_cnt)
	if err != nil {
		t.Errorf("[TestNickNameHandler] Error %s when inserting a user into the table", err)
	}

	server, _ := tcp.StartTCPServer(username, password, hostname, dbname, tablename, tcp_port)
	time.Sleep(time.Second)

	// login first
	code, err := loginRequest(users[0])
	if err != nil {
		t.Errorf("[TestNickNameHandler] Error %s when send login request to the http server", err)
	}
	if code != http.StatusOK {
		t.Errorf("[TestNickNameHandler] handler returned wrong status code: got %v want %v", code, http.StatusOK)
	}

	// update nick name
	code, err = updateNicknameRequest(users[0])
	if err != nil {
		t.Errorf("[TestNickNameHandler] Error %s when send update nickname request to the http server", err)
	}
	if code != http.StatusOK {
		t.Errorf("[TestNickNameHandler] handler returned wrong status code: got %v want %v", code, http.StatusOK)
	}

	// check if we really update into db
	query_user, err := mysqldb.Select(db, users[0])
	if err != nil || query_user == nil {
		t.Errorf("[TestNickNameHandler] Error %s when querying the user", err)
	}
	if query_user.Nickname != users[0].Nickname {
		t.Errorf("[TestNickNameHandler] handler returned wrong user nickname: got %s want %s", query_user.Nickname, users[0].Nickname)
	}

	// delete table
	err = mysqldb.DeleteTable(db)
	if err != nil {
		t.Errorf("[TestNickNameHandler] %s", err)
	}

	server.Stop()
}

func TestPhotoHandler(t *testing.T) {
	db, err := createConnAndTable()
	if err != nil {
		t.Errorf("[TestPhotoHandler] Error %s when getting db connection and creating a table", err)
	}
	defer db.Conn.Close()

	users_cnt := 2
	users, err := createRandomUsers(db, users_cnt)
	if err != nil {
		t.Errorf("[TestPhotoHandler] Error %s when inserting a user into the table", err)
	}

	server, _ := tcp.StartTCPServer(username, password, hostname, dbname, tablename, tcp_port)
	time.Sleep(time.Second)

	// login first
	code, err := loginRequest(users[0])
	if err != nil {
		t.Errorf("[TestPhotoHandler] Error %s when send login request to the http server", err)
	}
	if code != http.StatusOK {
		t.Errorf("[TestPhotoHandler] handler returned wrong status code: got %v want %v", code, http.StatusOK)
	}

	photo_tag := "photo"
	code, err = updatePhotoRequest(users[0], photo_tag)
	if err != nil {
		t.Errorf("[TestPhotoHandler] Error %s when send update photo request to the http server", err)
	}
	if code != http.StatusOK {
		t.Errorf("[TestPhotoHandler] handler returned wrong status code: got %v want %v", code, http.StatusOK)
	}

	// check if we really update into db
	query_user, err := mysqldb.Select(db, users[0])
	if err != nil || query_user == nil {
		t.Errorf("[TestPhotoHandler] Error %s when querying the user", err)
	}
	if query_user.Nickname != users[0].Nickname {
		t.Errorf("[TestPhotoHandler] handler returned wrong user photo: got %s want %s", query_user.Photo, users[0].Photo)
	}

	// delete table
	err = mysqldb.DeleteTable(db)
	if err != nil {
		t.Errorf("[TestPhotoHandler] %s", err)
	}

	server.Stop()
}

func TestUnauthUser(t *testing.T) {
	db, err := createConnAndTable()
	if err != nil {
		t.Errorf("[TestUnauthUser] Error %s when getting db connection and creating a table", err)
	}
	defer db.Conn.Close()

	users_cnt := 2
	users, err := createRandomUsers(db, users_cnt)
	if err != nil {
		t.Errorf("[TestUnauthUser] Error %s when inserting a user into the table", err)
	}

	// test update nick name when bypass login phase, should redirect to login page
	code, err := updateNicknameRequest(users[0])
	if err != nil {
		t.Errorf("[TestUnauthUser] Error %s when send update nickname request to the http server", err)
	}
	if code != http.StatusOK {
		t.Errorf("[TestUnauthUser] handler returned wrong status code: got %v want %v", code, http.StatusOK)
	}

	// test update photo when bypass login phase, should redirect to login page
	photo_tag := "photo"
	code, err = updatePhotoRequest(users[0], photo_tag)
	if err != nil {
		t.Errorf("[TestUnauthUser] Error %s when send update photo request to the http server", err)
	}
	if code != http.StatusOK {
		t.Errorf("[TestUnauthUser] handler returned wrong status code: got %v want %v", code, http.StatusOK)
	}

	// delete table
	err = mysqldb.DeleteTable(db)
	if err != nil {
		t.Errorf("[TestUnauthUser] %s", err)
	}
}

func TestTCPServerDown(t *testing.T) {
	db, err := createConnAndTable()
	if err != nil {
		t.Errorf("[TestTCPServerDown] Error %s when getting db connection and creating a table", err)
	}
	defer db.Conn.Close()

	server, _ := tcp.StartTCPServer(username, password, hostname, dbname, tablename, tcp_port)
	time.Sleep(time.Second)

	users_cnt := 2
	users, err := createRandomUsers(db, users_cnt)
	if err != nil {
		t.Errorf("[TestTCPServerDown] Error %s when inserting a user into the table", err)
	}

	// login first to get the certification of user by http server
	code, err := loginRequest(users[0])
	if err != nil {
		t.Errorf("[TestTCPServerDown] Error %s when send login request to the http server", err)
	}
	if code != http.StatusOK {
		t.Errorf("[TestTCPServerDown] handler returned wrong status code: got %v want %v", code, http.StatusOK)
	}

	// stop the server
	server.Stop()

	// test update nick name after tcp server down
	code, err = updateNicknameRequest(users[0])
	if err != nil {
		t.Errorf("[TestTCPServerDown] Error %s when send update nickname request to the http server", err)
	}
	if code != http.StatusInternalServerError {
		t.Errorf("[TestTCPServerDown] handler returned wrong status code: got %v want %v", code, http.StatusInternalServerError)
	}

	// test update photo after tcp server down
	photo_tag := "photo"
	code, err = updatePhotoRequest(users[0], photo_tag)
	if err != nil {
		t.Errorf("[TestTCPServerDown] Error %s when send update photo request to the http server", err)
	}
	if code != http.StatusInternalServerError {
		t.Errorf("[TestTCPServerDown] handler returned wrong status code: got %v want %v", code, http.StatusInternalServerError)
	}

	// test login after tcp server down
	code, err = loginRequest(users[0])
	if err != nil {
		t.Errorf("[TestTCPServerDown] Error %s when send login request to the http server", err)
	}
	if code != http.StatusOK {
		t.Errorf("[TestTCPServerDown] handler returned wrong status code: got %v want %v", code, http.StatusOK)
	}

	// delete table
	err = mysqldb.DeleteTable(db)
	if err != nil {
		t.Errorf("[TestPhotoHandler] %s", err)
	}
}

func TestWrongPhotoUploadDir(t *testing.T) {
	db, err := createConnAndTable()
	if err != nil {
		t.Errorf("[TestTCPServerDown] Error %s when getting db connection and creating a table", err)
	}
	defer db.Conn.Close()

	server, _ := tcp.StartTCPServer(username, password, hostname, dbname, tablename, tcp_port)
	time.Sleep(time.Second)

	users_cnt := 2
	users, err := createRandomUsers(db, users_cnt)
	if err != nil {
		t.Errorf("[TestTCPServerDown] Error %s when inserting a user into the table", err)
	}

	// login first to get the certification of user by http server
	code, err := loginRequest(users[0])
	if err != nil {
		t.Errorf("[TestTCPServerDown] Error %s when send login request to the http server", err)
	}
	if code != http.StatusOK {
		t.Errorf("[TestTCPServerDown] handler returned wrong status code: got %v want %v", code, http.StatusOK)
	}

	// give invalid photo address
	photo_tag := "image"
	code, err = updatePhotoRequest(users[0], photo_tag)
	if err != nil {
		t.Errorf("[TestWrongPhotoUploadDir] Should not pop out the error: %s", err)
	}
	if code != http.StatusInternalServerError {
		t.Errorf("[TestWrongPhotoUploadDir] handler returned wrong status code: got %v want %v", code, http.StatusInternalServerError)
	}

	// delete table
	err = mysqldb.DeleteTable(db)
	if err != nil {
		t.Errorf("[TestPhotoHandler] %s", err)
	}

	server.Stop()
}
