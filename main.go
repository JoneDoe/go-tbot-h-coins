package main

import (
	app "go-tbot-h-coins/src"
	"log"
	"sync"

	"github.com/joho/godotenv"
)

var telegramBot app.TelegramBot

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cli := CLI{}
	cli.Run()

	app.Connection.Init()

	telegramBot.Init() // Инициализация бота
	telegramBot.Start()

	// 	if update.Message.IsCommand() {
	// 		CommandHandler(update.Message, bot)

	// 		pushes := make(chan string)
	// 		//wg.Add(1)

	// 		go OrderSubscriber(update.Message.Chat.ID, pushes, wg)

	// 		for push := range pushes {
	// 			fmt.Println("received message", push)

	// 			msg := tgbotapi.NewMessage(update.Message.Chat.ID, push)
	// 			bot.Send(msg)
	// 		}
	// 		//close(pushes)
	// 	}

	// 	/* b, _ := json.Marshal(update.Message)
	// 	fmt.Println(string(b)) */

	// 	//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	// 	/* b, _ := json.Marshal(update.Message)
	// 	fmt.Println(string(b)) */
	// 	//log.Printf("[%s]", update.Message.Contact)

	// 	//msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg)
	// 	//msg.ReplyToMessageID = update.Message.MessageID

	// 	//bot.Send(msg)
	// }
	// wg.Wait()
}

// OrderSubscriber for user
func OrderSubscriber(chatID int64, chanel chan string, wg *sync.WaitGroup) {
	//defer wg.Done()

	//chanel <- "You subscribed to push channel"
	//chanel <- "You subscribed to push channe2"

	Subscribe(chanel, chatID)
}
