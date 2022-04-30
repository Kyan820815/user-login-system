package mysqldb

import (
	"fmt"
	"testing"
)

const (
	username  = "root"
	password  = "12345678"
	hostname  = "localhost:3306"
	dbname    = "entry_task_user_db"
	tablename = "entry_task_db_test_table"
)

func CreateConnAndTable() (*DB, error) {
	dbInfo := GetInfo(username, password, hostname, dbname, tablename)
	db, err := DbConnection(dbInfo)
	if err != nil {
		return nil, err
	}

	err = CreateUserTable(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func TestDBConnection(t *testing.T) {
	dbInfo := GetInfo(username, password, hostname, dbname, tablename)
	db, err := DbConnection(dbInfo)
	if err != nil {
		t.Errorf("[TestDBConnection] Error %s when getting db connection", err)
	}
	defer db.Conn.Close()
}

func TestDBCreateTable(t *testing.T) {
	db, err := CreateConnAndTable()
	if err != nil {
		t.Errorf("[TestDBCreateTable] Error %s when getting db connection and creating a table", err)
	}
	defer db.Conn.Close()

	statement := fmt.Sprintf("SELECT * FROM %s", db.DbInfo.tablename)
	form, err := db.Conn.Prepare(statement)
	defer form.Close()
	if err != nil {
		t.Errorf("[TestDBCreateTable] Error %s when preparing the statement", err)
	}
	_, err = form.Exec()
	if err != nil {
		t.Errorf("[TestDBCreateTable] Error %s when searching table that exists", err)
	}
}

func TestDeleteTable(t *testing.T) {
	db, err := CreateConnAndTable()
	if err != nil {
		t.Errorf("[TestDeleteTable] Error %s when getting db connection and creating a table", err)
	}
	defer db.Conn.Close()

	err = DeleteTable(db)
	if err != nil {
		t.Errorf("[TestDeleteTable] Error %s when deleting tables", err)
	}
}

func TestRowOperation(t *testing.T) {
	db, err := CreateConnAndTable()
	if err != nil {
		t.Errorf("[TestRowOperation] Error %s when getting db connection and creating a table", err)
	}
	defer db.Conn.Close()

	// insert a user and select it
	user := &User{
		Acc:      "TestRowOperation_user",
		Pwd:      "TestRowOperation_pwd",
		Nickname: "TestRowOperation_nickname",
		Photo:    "TestRowOperation_photo",
		Id:       -1,
	}

	err = Insert(db, user)
	if err != nil {
		t.Errorf("[TestRowOperation] Error %s when inserting the user", err)
	}

	user, err = Select(db, user)
	if err != nil || user == nil {
		t.Errorf("[TestRowOperation] Error %s when querying the user", err)
	}

	// Delete this user
	err = Delete(db, user)
	if err != nil {
		t.Errorf("[TestRowOperation] Error %s when deleting the user", err)
	}

	// select a non-existent row
	user, err = Select(db, user)
	if user != nil {
		t.Errorf("[TestRowOperation] Error %s when querying the user", err)
	}

	// delete table
	err = DeleteTable(db)
	if err != nil {
		t.Errorf("[TestRowOperation] Error %s when deleting tables", err)
	}
}

func TestUpdateNickname(t *testing.T) {
	db, err := CreateConnAndTable()
	if err != nil {
		t.Errorf("[TestUpdateNickname] Error %s when getting db connection and creating a table", err)
	}
	defer db.Conn.Close()

	// insert a user and select it
	user := &User{
		Acc:      "TestUpdateNickname_user",
		Pwd:      "TestUpdateNickname_pwd",
		Nickname: "TestUpdateNickname_nickname",
		Photo:    "TestUpdateNickname_photo",
		Id:       -1,
	}

	err = Insert(db, user)
	if err != nil {
		t.Errorf("[TestUpdateNickname] Error %s when inserting the user", err)
	}

	user, err = Select(db, user)
	if err != nil || user == nil {
		t.Errorf("[TestUpdateNickname] Error %s when querying the user", err)
	}

	user.Nickname = "TestUpdateNickname_new_nickname"
	err = UpdateNickname(db, user)
	if err != nil {
		t.Errorf("[TestUpdateNickname] Error %s when revising the user in db", err)
	}

	// check if we really updated
	query_user, err := Select(db, user)
	if query_user.Nickname != user.Nickname {
		t.Errorf("[TestUpdateNickname] Error %s when updating the user", err)
	}

	// delete table
	err = DeleteTable(db)
	if err != nil {
		t.Errorf("[TestUpdateNickname] Error %s when deleting tables", err)
	}
}

func TestUpdatePhoto(t *testing.T) {
	db, err := CreateConnAndTable()
	if err != nil {
		t.Errorf("[TestUpdatePhoto] Error %s when getting db connection and creating a table", err)
	}
	defer db.Conn.Close()

	// insert a user and select it
	user := &User{
		Acc:      "TestUpdatePhoto_user",
		Pwd:      "TestUpdatePhoto_pwd",
		Nickname: "TestUpdatePhoto_nickname",
		Photo:    "TestUpdatePhoto_photo",
		Id:       -1,
	}

	err = Insert(db, user)
	if err != nil {
		t.Errorf("[TestUpdatePhoto] Error %s when inserting the user", err)
	}

	user, err = Select(db, user)
	if err != nil || user == nil {
		t.Errorf("[TestUpdatePhoto] Error %s when querying the user", err)
	}

	user.Photo = "TestUpdatePhoto_new_photo"
	err = UpdatePhoto(db, user)
	if err != nil {
		t.Errorf("[TestUpdatePhoto] Error %s when revising the user in db", err)
	}

	// check if we really updated
	query_user, err := Select(db, user)
	if query_user.Photo != user.Photo {
		t.Errorf("[TestUpdatePhoto] Error %s when updating the user", err)
	}

	// delete table
	err = DeleteTable(db)
	if err != nil {
		t.Errorf("[TestUpdatePhoto] Error %s when deleting tables", err)
	}
}
