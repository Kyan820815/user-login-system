package mysqldb

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DbInfo struct {
	username  string
	password  string
	hostname  string
	dbname    string
	tablename string
}

type User struct {
	Acc      string
	Pwd      string
	Nickname string
	Photo    string
	Id       int32
}

type DB struct {
	Conn   *sql.DB
	DbInfo *DbInfo
}

func GetInfo(username string, password string, hostname string, dbname string, tablename string) *DbInfo {
	return &DbInfo{
		username:  username,
		password:  password,
		hostname:  hostname,
		dbname:    dbname,
		tablename: tablename,
	}
}

func dsn(dbInfo *DbInfo) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", dbInfo.username, dbInfo.password, dbInfo.hostname, dbInfo.dbname)
}

func DbConnection(dbInfo *DbInfo) (*DB, error) {
	conn, err := sql.Open("mysql", dsn(dbInfo))
	if err != nil {
		fmt.Printf("[DB DbConnection] Error %s when opening DB\n", err)
		return nil, err
	}

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := conn.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+dbInfo.dbname)
	if err != nil {
		fmt.Printf("[DB DbConnection] Error %s when creating DB\n", err)
		return nil, err
	}

	_, err = res.RowsAffected()
	if err != nil {
		fmt.Printf("[DB DbConnection] Error %s when fetching rows\n", err)
		return nil, err
	}
	// fmt.Printf("[DB DbConnection] rows affected %d\n", no)
	conn.Close()

	conn, err = sql.Open("mysql", dsn(dbInfo))
	if err != nil {
		fmt.Printf("[DB DbConnection] Error %s when opening DB\n", err)
		return nil, err
	}

	conn.SetMaxOpenConns(20)
	conn.SetMaxIdleConns(20)
	conn.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	err = conn.PingContext(ctx)
	if err != nil {
		fmt.Printf("[DB DbConnection] Errors %s pinging DB\n", err)
		return nil, err
	}

	fmt.Printf("[DB DbConnection] Connected to DB %s successfully\n", dbInfo.dbname)
	return &DB{Conn: conn, DbInfo: dbInfo}, nil
}

func CreateUserTable(db *DB) error {
	statement := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(user_id int primary key auto_increment, 
                                                    user_name text, 
                                                    user_password text, 
                                                    user_nickname text, 
                                                    user_photo text)`, db.DbInfo.tablename)
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	res, err := db.Conn.ExecContext(ctx, statement)
	if err != nil {
		fmt.Printf("[DB CreateUserTable] Error %s when creating product table\n", err)
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		fmt.Printf("[DB CreateUserTable] Error %s when getting rows affected\n", err)
		return err
	}

	// fmt.Printf("[DB CreateUserTable] Rows affected when creating table: %d\n", rows)
	return nil
}

func DeleteTable(db *DB) error {
	statement := fmt.Sprintf("DROP TABLE %s", db.DbInfo.tablename)
	form, err := db.Conn.Prepare(statement)
	defer form.Close()
	_, err = form.Exec()
	return err
}

func Delete(db *DB, user *User) error {
	statement := fmt.Sprintf("DELETE FROM %s WHERE user_name=? AND user_password=?", db.DbInfo.tablename)
	delForm, err := db.Conn.Prepare(statement)
	defer delForm.Close()
	if err != nil {
		return err
	}

	_, err = delForm.Exec(user.Acc, user.Pwd)
	return err
}

func Select(db *DB, user *User) (*User, error) {
	statement := fmt.Sprintf("SELECT user_id, user_nickname, user_photo FROM %s where user_name=? AND user_password=?", db.DbInfo.tablename)
	selForm, err := db.Conn.Prepare(statement)
	defer selForm.Close()
	if err != nil {
		return nil, err
	}
	query_user := &User{}
	query_user.Acc = user.Acc
	query_user.Pwd = user.Pwd

	err = selForm.QueryRow(user.Acc, user.Pwd).Scan(&query_user.Id, &query_user.Nickname, &query_user.Photo)
	if err != nil {
		return nil, err
	}
	return query_user, nil
}

func Insert(db *DB, user *User) error {
	_, err := Select(db, user)
	// already have this user in the db
	if err == nil {
		return nil
	}

	statement := fmt.Sprintf("INSERT INTO %s(user_name, user_password, user_nickname, user_photo) VALUES(?,?,?,?)", db.DbInfo.tablename)
	insForm, err := db.Conn.Prepare(statement)
	defer insForm.Close()
	if err != nil {
		return err
	}
	_, err = insForm.Exec(user.Acc, user.Pwd, user.Nickname, user.Photo)
	return err
}

func UpdateNickname(db *DB, user *User) error {
	statement := fmt.Sprintf("UPDATE %s SET user_nickname=? WHERE user_id=?", db.DbInfo.tablename)
	updForm, err := db.Conn.Prepare(statement)
	defer updForm.Close()
	if err != nil {
		return err
	}
	_, err = updForm.Exec(user.Nickname, user.Id)

	return err
}

func UpdatePhoto(db *DB, user *User) error {
	statement := fmt.Sprintf("UPDATE %s SET user_photo=? WHERE user_id=?", db.DbInfo.tablename)
	updForm, err := db.Conn.Prepare(statement)
	defer updForm.Close()
	if err != nil {
		return err
	}
	_, err = updForm.Exec(user.Photo, user.Id)

	return err
}
