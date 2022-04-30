package myredis

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type RDS struct {
	Conn *redis.Client
}

func RdsConnection() *RDS {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password
		DB:       0,  // default DB
	})

	return &RDS{Conn: redisClient}
}

func Flush(rds *RDS) error {
	return rds.Conn.FlushDB(ctx).Err()
}

func Ping(conn *redis.Client) (string, error) {
	result, err := conn.Ping(ctx).Result()
	if err != nil {
		return "", err
	} else {
		return result, nil
	}
}

func Set(rds *RDS, key string, val interface{}) error {
	enc_val, err := json.Marshal(val)
	if err != nil {
		return err
	}

	err = rds.Conn.Set(ctx, key, enc_val, 0).Err()
	return err
}

func Get(rds *RDS, key string, val interface{}) error {
	enc_val, err := rds.Conn.Get(ctx, key).Bytes()
	if err == redis.Nil {
		val = nil
		return nil
	} else if err != nil {
		return err
	} else {
		err = json.Unmarshal(enc_val, val)
		if err != nil {
			return err
		} else {
			return nil
		}
	}
}
