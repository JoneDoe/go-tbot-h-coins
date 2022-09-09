package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/go-redis/redis"
)

var client *redis.Client

func redisCli() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PWD"), // no password set
		DB:       0,                      // use default DB
	})
	return client
}

// Subscribe to chanel
func Subscribe(chanel chan string, chatID int64) {

	client := redisCli()

	pubsub := client.Subscribe("111910124")
	defer pubsub.Close()

	//log.Printf("chatID %s", string(chatID))
	var wg sync.WaitGroup
	wg.Add(1)

	chanel <- "111910124"
	chanel <- "sdfsdfsdf"

	go func(chanel chan string) {
		defer wg.Done()
		for {
			msgi, err := pubsub.ReceiveMessage()
			if err != nil {
				panic(err)
			}

			switch msgi.Payload {
			case "stop":
				fmt.Println("stoped", msgi.Channel)
				return
			default:
				fmt.Println("received", msgi.Payload, "from", msgi.Channel)
				chanel <- msgi.Payload
			}
		}
	}(chanel)

	wg.Wait()
}
