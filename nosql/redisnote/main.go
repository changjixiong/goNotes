package main

import (
	"fmt"

	redis "gopkg.in/redis.v4"
)

var rds *redis.Client

func init() {
	rds = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "123456",
		DB:       0,
	})

}

func stringSet(key, value string) {
	rds.Set(key, value, 0)
}

func stringGet(key string) (string, error) {
	return rds.Get(key).Result()
}

func stringDel(key ...string) (int64, error) {

	return rds.Del(key...).Result()
}
func hashSet(key, field, value string) {
	rds.HSet(key, field, value)
}

func hashGet(key, field string) string {
	result, _ := rds.HGet(key, field).Result()

	return result
}

func hashDel(key string, fields ...string) (int64, error) {
	return rds.HDel(key, fields...).Result()

}

func main() {
	// hashSet("hashA", "k1", "v1")
	// hashSet("hashA", "k2", "v2")
	// fmt.Println("hashGet", hashGet("hashA", "k1"))
	// fmt.Println("hashGet", hashGet("hashA", "k3")) //nil
	// hashDel("hashA", "k2")
	// stringSet("string1", "v1")
	// stringSet("string2", "v2")
	// fmt.Println(stringGet("string1"))
	// fmt.Println(stringGet("string2"))
	stringDel("string1")
	s, err := stringGet("string1")
	fmt.Println("s:", s, "err:", err)

	n, err := stringDel("string1")

	fmt.Println("n:", n, "err:", err)

}
