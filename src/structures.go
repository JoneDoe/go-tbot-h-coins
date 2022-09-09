package app

import (
	"github.com/go-redis/redis"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

// TelegramBot struct
type TelegramBot struct {
	API                   *tgbotapi.BotAPI        // API телеграмма
	Updates               tgbotapi.UpdatesChannel // Канал обновлений
	ActiveContactRequests []int64                 // ID чатов, от которых мы ожидаем номер
}

// Redis client
type Redis struct {
	Client *redis.Client
}

// User identity
type teUser struct {
	ChatID      int64
	PhoneNumber string
}

// Request struct
type Request struct {
	ID       int64
	Msg      tgbotapi.Update
	Callback *tgbotapi.CallbackQuery
	Cmd      string
}

// CoinsTransaction struct
type CoinsTransaction struct {
	Amount            int
	Sender, Recipient string
}

// Amounts for transaction
var Amounts = map[string]int{
	"amount.5":   5,
	"amount.10":  10,
	"amount.50":  50,
	"amount.100": 100,
	"amount.500": 500,
}
