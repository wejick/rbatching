package main

import (
	"fmt"
	"log"

	"github.com/gomodule/redigo/redis"
	"github.com/wejick/rbatching"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	var redisPool = &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6380")
		},
	}
	batch := rbatching.NewRBatcher("testcase", 4, redisPool)

	batch.Enqueue("1")
	batch.Enqueue("2")
	batch.Enqueue("3")
	batch.Enqueue("4")

	{
		val, err := batch.GetBatch()
		if err != nil {
			log.Println(err)
			return
		}
		for _, val := range val {
			string := string(val.([]byte))
			fmt.Println(string)
		}
		batch.CloseBatch()
	}
	{
		val, err := batch.GetBatch()
		if err != nil {
			log.Println(err)
			return
		}
		for _, val := range val {
			string := string(val.([]byte))
			fmt.Println(string)
		}
		batch.CloseBatch()
	}
}
