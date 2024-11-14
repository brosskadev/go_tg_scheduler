package schedulebot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func InitBot(token string) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Установка команд для бота
	commands := []tgbotapi.BotCommand{
		{Command: "/help", Description: "Получить список доступных команд"},
		{Command: "/writeinfo", Description: "Введите информацию"},
		{Command: "/getschedule", Description: "Получить расписание на эту неделю"},
	}

	_, err = bot.Request(tgbotapi.NewSetMyCommands(commands...))
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.CallbackQuery != nil {
			// Обработка нажатия на инлайн-кнопку
			HandleCallBack(bot, update)
			continue
		}

		if update.Message != nil { // Проверка на наличие команды
			HandleCommands(update, *bot)
		}
	}
}
