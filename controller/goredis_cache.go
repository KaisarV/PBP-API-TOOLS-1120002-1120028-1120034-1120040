package controllers

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

// ctx (global) for redis
var ctx = context.Background()

// GoRedis - Set & Get from Redis
func SetRedis(rdb *redis.Client, key string, value string, expiration int) {
	err := rdb.Set(ctx, key, value, 0).Err()
	if err != nil {
		log.Fatal(err)
		// fmt.Println(err)
	}
}

func GetRedis(rdb *redis.Client, key string) string {
	val, err := rdb.Get(ctx, key).Result()

	if err != nil {
		log.Fatal(err)
		// fmt.Println(err)
	}
	return val
}

// func RDB()  {
// 	// GoRedis
// 	rdb := redis.NewClient(&redis.Options{
// 		Addr:     "localhost:6379",
// 		Password: "", // no password set
// 		DB:       0,  // use default DB
// 	})
// 	return rdb
// }
