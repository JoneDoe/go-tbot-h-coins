package app

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/go-redis/redis"
)

// Connection connection
var Connection Redis

// Init redis
func (r *Redis) Init() {
	r.Client = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PWD"), // no password set
		DB:       0,                      // use default DB
	})
}

// Find Check if user exists
func (r *Redis) Find(chatID int64) bool {
	record, _ := r.Client.HExists("teUsers", strconv.FormatInt(chatID, 10)).Result()
	return record
}

// GetUser find user by id
func (r *Redis) GetUser(chatID int64) teUser {
	record, _ := r.Client.HGet("teUsers", strconv.FormatInt(chatID, 10)).Result()
	var user teUser
	json.Unmarshal([]byte(record), &user)

	return user
}

// CreateUser Создание пользователя
func (r *Redis) CreateUser(user teUser) error {
	e, err := json.Marshal(user)
	r.Client.HSet("teUsers", strconv.FormatInt(user.ChatID, 10), string(e))
	return err
}

// UpdateUser Обновление номера мобильного телефона
func (r *Redis) UpdateUser(user teUser) error {
	return r.CreateUser(user)
}
