package myredis

import (
	"testing"
)

func TestRDSConn(t *testing.T) {
	rds := RdsConnection()
	_, err := Ping(rds.Conn)
	if err != nil {
		t.Errorf("[TestRDSConn]: %s\n", err)
	}
}

func TestSet(t *testing.T) {
	rds := RdsConnection()

	// flush memory
	Flush(rds)

	key := "rds_key"
	val := "rds_val"
	err := Set(rds, key, val)
	if err != nil {
		t.Errorf("[TestSet]: %s\n", err)
	}

	var query_val string
	err = Get(rds, key, &query_val)
	if err != nil {
		t.Errorf("[TestSet]: %s\n", err)
	}
	if query_val != val {
		t.Errorf("[TestSet] handler returned wrong value : got %s want %s", query_val, val)
	}
}

func TestGet(t *testing.T) {
	rds := RdsConnection()

	// flush memory
	Flush(rds)

	// query non-existent key-value pair
	key := "rds_key"
	var query_val string
	err := Get(rds, key, &query_val)
	if err != nil {
		t.Errorf("[TestSet]: %s\n", err)
	}
	if err != nil || query_val != "" {
		t.Errorf("[TestGet]: %s\n", err)
	}
}
