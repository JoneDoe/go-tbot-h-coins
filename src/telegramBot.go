package app

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	api "go-tbot-h-coins/src/api"
	blockchain "go-tbot-h-coins/src/blockchain"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

var request Request
var transaction CoinsTransaction
var transactionPending map[int]string

// Init bot
func (telegramBot *TelegramBot) Init() {
	botAPI, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_API_KEY")) // Инициализация API
	if err != nil {
		log.Fatal(err)
	}

	//botAPI.Debug = true

	log.Printf("Authorized on account %s", botAPI.Self.UserName)

	telegramBot.API = botAPI
	botUpdate := tgbotapi.NewUpdate(0) // Инициализация канала обновлений
	botUpdate.Timeout = 64
	botUpdates, err := telegramBot.API.GetUpdatesChan(botUpdate)
	if err != nil {
		log.Fatal(err)
	}
	telegramBot.Updates = botUpdates

	Handler.bc = blockchain.NewBlockchain()
}

// Start bot
func (telegramBot *TelegramBot) Start() {

	transactionPending = make(map[int]string)

	for update := range telegramBot.Updates {
		//log.Println(update.Message.Contact)
		if update.CallbackQuery != nil {
			request.Callback = update.CallbackQuery
			telegramBot.callbackHandler()
		} else if update.Message != nil {
			request.ID = update.Message.Chat.ID

			if update.Message.IsCommand() {
				request.Cmd = update.Message.Command()
				telegramBot.commandHandler()
			} else {
				request.Msg = update
				// Если сообщение есть  -> начинаем обработку
				telegramBot.updateHandler()
			}
		}
	}
}

func (telegramBot *TelegramBot) btnClickHandler() {
	command := request.Msg.Message.Text

	fmt.Println(fmt.Sprintf("Bot btnClickHandler %s", command))

	switch command {
	case BTN_COINS:
		telegramBot.drawCoinsMenu(tgbotapi.NewMessage(request.ID, "Что вы хотите сделать"))
	case BTN_COINS_BALANCE:
		telegramBot.execCoinBalance()
	case BTN_COINS_SEND:
		telegramBot.drawCoinsChooseMenu()
	case BTN_COINS_BACK:
		telegramBot.drawUserMenu()
	case BTN_VOCATIONS:
		msg := tgbotapi.NewMessage(request.ID, "0 дней отпуска, прийдется по-пахать еще, дружок 😈")
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
		telegramBot.API.Send(msg)
		telegramBot.drawUserMenu()
	}
}

func (telegramBot *TelegramBot) callbackHandler() {
	callbackData := request.Callback.Data

	log.Println("Choose case", callbackData)

	callbackCommand := strings.Split(callbackData, ".")

	switch callbackCommand[0] {
	case "amount":
		if _, exists := transactionPending[request.Callback.From.ID]; exists {
			//log.Println(request.Callback.From.FirstName)
			transaction.Sender = request.Callback.From.FirstName
			transaction.Amount, _ = strconv.Atoi(callbackCommand[1])

			text := telegramBot.execTransferCoins(transaction)
			telegramBot.drawCoinsMenu(tgbotapi.NewMessage(request.ID, text))

			delete(transactionPending, request.Callback.From.ID)
		}
	default:
		transaction.Recipient = callbackData
		transactionPending[request.Callback.From.ID] = callbackData

		//telegramBot.API.Send(tgbotapi.NewMessage(request.Callback.Message.Chat.ID, "Получатель "+callbackData))
		telegramBot.drawCoinsChooseAmount()
	}
}

func (telegramBot *TelegramBot) execTransferCoins(transaction CoinsTransaction) (text string) {
	if err := Handler.CoinsTransfer(transaction); err != nil {
		text = fmt.Sprintf("Я бы перевел твои %d коинов %s, но у тебя недостаточно их 🤡, проверь свой баланс", transaction.Amount, transaction.Recipient)
	} else {
		text = fmt.Sprintf("Ты перечислил %d коинов %s", transaction.Amount, api.GetUserList()[transaction.Recipient])
	}

	return
}

func (telegramBot *TelegramBot) updateHandler() {
	chatID := request.ID

	if User.Find(chatID) { // Есть ли пользователь в БД?
		telegramBot.btnClickHandler()
		telegramBot.analyzeUser()
	} else {
		User.Create(chatID)          // Создаём пользователя
		telegramBot.requestContact() // Запрашиваем номер
	}
}

// func (telegramBot *TelegramBot) findUser(chatID int64) bool {
// 	find := Connection.Find(chatID)
// 	return find
// }

// func (telegramBot *TelegramBot) createUser(user teUser) {
// 	err := Connection.CreateUser(user)
// 	if err != nil {
// 		msg := tgbotapi.NewMessage(user.ChatID, "Произошла ошибка! Бот может работать неправильно!")
// 		telegramBot.API.Send(msg)
// 	}
// }

func (telegramBot *TelegramBot) requestContact() {
	// Создаём сообщение
	requestContactMessage := tgbotapi.NewMessage(request.ID, "Согласны ли вы предоставить ваш номер телефона для регистрации в системе?")
	// Создаём кнопку отправки контакта
	acceptButton := tgbotapi.NewKeyboardButtonContact("Да")
	declineButton := tgbotapi.NewKeyboardButton("Нет")
	// Создаём клавиатуру
	requestContactMessage.ReplyMarkup = tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{acceptButton, declineButton})
	telegramBot.API.Send(requestContactMessage) // Отправляем сообщение

	telegramBot.addContactRequestID() // Добавляем ChatID в лист ожидания
}

func (telegramBot *TelegramBot) addContactRequestID() {
	telegramBot.ActiveContactRequests = append(telegramBot.ActiveContactRequests, request.ID)
}

func (telegramBot *TelegramBot) analyzeUser() {
	chatID := request.ID
	user := User.Get(chatID) // Вытаскиваем данные из БД для проверки номера
	// if err != nil {
	// 	msg := tgbotapi.NewMessage(chatID, "Произошла ошибка! Бот может работать неправильно!")
	// 	telegramBot.API.Send(msg)
	// 	return
	// }
	if len(user.PhoneNumber) > 0 {
		//telegramBot.drawUserMenu()
		// msg := tgbotapi.NewMessage(chatID, "Ваш номер: "+user.PhoneNumber) // Если номер у нас уже есть, то пишем его
		// telegramBot.API.Send(msg)
		return
	} else if telegramBot.findContactRequestID() {
		// Если номера нет, то проверяем ждём ли мы контакт от этого ChatID
		telegramBot.checkRequestContactReply() // Если да -> проверяем
		return
	} else {
		telegramBot.requestContact() // Если нет -> запрашиваем его
		return
	}
}

func (telegramBot *TelegramBot) drawUserMenu() {
	msg := tgbotapi.NewMessage(request.ID, "Меню пожеланий активировано, чего пожелает твоя душа?")

	// Создаём клавиатуру
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButton(BTN_COINS),
		tgbotapi.NewKeyboardButton(BTN_VOCATIONS)})
	telegramBot.API.Send(msg)
}

func (telegramBot *TelegramBot) drawCoinsChooseMenu() {
	keyboard := tgbotapi.InlineKeyboardMarkup{}
	for key, recipient := range api.GetUserList() {
		var row []tgbotapi.InlineKeyboardButton
		btn := tgbotapi.NewInlineKeyboardButtonData(recipient, key)
		row = append(row, btn)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	}

	msg := tgbotapi.NewMessage(request.ID, "Выбери кому хочешь переслать свои H-coins")
	msg.ReplyMarkup = keyboard
	telegramBot.API.Send(msg)
}

func (telegramBot *TelegramBot) drawCoinsChooseAmount() {
	keyboard := tgbotapi.InlineKeyboardMarkup{}
	var row []tgbotapi.InlineKeyboardButton

	for key, amount := range Amounts {
		btn := tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(amount), key)
		row = append(row, btn)
	}

	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	msg := tgbotapi.NewMessage(request.ID, "Сколько H-coins будем пересылать?")
	msg.ReplyMarkup = keyboard
	telegramBot.API.Send(msg)
}

func (telegramBot *TelegramBot) drawCoinsMenu(msg tgbotapi.MessageConfig) {
	row1 := []tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButton(BTN_COINS_BALANCE),
		tgbotapi.NewKeyboardButton(BTN_COINS_SEND)}

	row2 := []tgbotapi.KeyboardButton{tgbotapi.NewKeyboardButton(BTN_COINS_BACK)}
	// Создаём клавиатуру
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(row1, row2)

	telegramBot.API.Send(msg)
}

// CommandHandler handle commands from message
func (telegramBot *TelegramBot) commandHandler() {
	command := request.Msg.Message.Text

	fmt.Println(fmt.Sprintf("Bot command %s", command))
}

func (telegramBot *TelegramBot) execCoinBalance() {
	account := api.GetCoins(User.Get(request.ID).PhoneNumber)
	msg := tgbotapi.NewMessage(request.ID, fmt.Sprintf("На вашем счету %d коинов", account.Balance))
	telegramBot.API.Send(msg)
}

// Проверка принятого контакта
func (telegramBot *TelegramBot) checkRequestContactReply() {
	chatID := request.ID
	message := request.Msg.Message

	if message.Contact != nil { // Проверяем, содержит ли сообщение контакт
		if message.Contact.UserID == message.From.ID { // Проверяем действительно ли это контакт отправителя
			err := User.Update(teUser{chatID, message.Contact.PhoneNumber}) // Обновляем номер
			if err != nil {
				msg := tgbotapi.NewMessage(chatID, "Произошла ошибка! Бот может работать неправильно!")
				telegramBot.API.Send(msg)
				return
			}

			telegramBot.deleteContactRequestID() // Удаляем ChatID из списка ожидания
			msg := tgbotapi.NewMessage(chatID, "Спасибо!")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false) // Убираем клавиатуру
			telegramBot.API.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(chatID, "Номер телефона, который вы предоставили, принадлежит не вам!")
			telegramBot.API.Send(msg)
			telegramBot.requestContact()
		}
	} else {
		msg := tgbotapi.NewMessage(chatID, "Если вы не предоставите ваш номер телефона, вы не сможете пользоваться системой!")
		telegramBot.API.Send(msg)
		telegramBot.requestContact()
	}
}

// Есть ChatID в листе ожидания?
func (telegramBot *TelegramBot) findContactRequestID() bool {
	for _, v := range telegramBot.ActiveContactRequests {
		if v == request.ID {
			return true
		}
	}
	return false
}

// Удаление ChatID из листа ожидания
func (telegramBot *TelegramBot) deleteContactRequestID() {
	for i, v := range telegramBot.ActiveContactRequests {
		if v == request.ID {
			copy(telegramBot.ActiveContactRequests[i:], telegramBot.ActiveContactRequests[i+1:])
			telegramBot.ActiveContactRequests[len(telegramBot.ActiveContactRequests)-1] = 0
			telegramBot.ActiveContactRequests = telegramBot.ActiveContactRequests[:len(telegramBot.ActiveContactRequests)-1]
		}
	}
}
