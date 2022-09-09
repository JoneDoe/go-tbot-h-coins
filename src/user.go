package app

// User teUser
var User teUser

func (user *teUser) Find(chatID int64) bool {
	find := Connection.Find(chatID)
	return find
}

func (user *teUser) Get(chatID int64) teUser {
	return Connection.GetUser(chatID)
}

func (user *teUser) Create(chatID int64) error {
	err := Connection.CreateUser(teUser{chatID, ""})
	return err
	// if err != nil {
	// 	msg := tgbotapi.NewMessage(user.ChatID, "Произошла ошибка! Бот может работать неправильно!")
	// 	telegramBot.API.Send(msg)
	// }
}

// Обновление номера мобильного телефона пользователя
func (user *teUser) Update(u teUser) error {
	err := Connection.UpdateUser(u)
	return err
}
