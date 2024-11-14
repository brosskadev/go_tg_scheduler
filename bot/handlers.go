package schedulebot

import (
	"edushedule2/scraper"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var KeyboardFac = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Факультет мониторинга окружающей среды", "2"),
	),
)

var KeyboardForm = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Очное", "2"),
	),
)

var KeyboardCourse = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("1 курс", "1"),
		tgbotapi.NewInlineKeyboardButtonData("2 курс", "2"),
		tgbotapi.NewInlineKeyboardButtonData("3 курс", "3"),
		tgbotapi.NewInlineKeyboardButtonData("4 курс", "4"),
		tgbotapi.NewInlineKeyboardButtonData("5 курс", "5"),
	),
)

var KeyboarGroup1 = tgbotapi.NewInlineKeyboardMarkup(
	//1курс
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("А41ИТ1", "298"),
		tgbotapi.NewInlineKeyboardButtonData("А41ИТ2", "299"),
		tgbotapi.NewInlineKeyboardButtonData("А41ИТ3", "300"),
		tgbotapi.NewInlineKeyboardButtonData("А41ИТ4", "301"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("А41МФ1", "295"),
		tgbotapi.NewInlineKeyboardButtonData("А41ПД1", "297"),
		tgbotapi.NewInlineKeyboardButtonData("А41ТТ1", "296"),
		tgbotapi.NewInlineKeyboardButtonData("А41ЯР1", "294"),
	),
)

var KeyboarGroup2 = tgbotapi.NewInlineKeyboardMarkup(
	//2курс
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("А31ИТ1", "255"),
		tgbotapi.NewInlineKeyboardButtonData("А31ИТ2", "256"),
		tgbotapi.NewInlineKeyboardButtonData("А31ИТ3", "290"),
		tgbotapi.NewInlineKeyboardButtonData("А31МФ", "252"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("А31ПД1", "254"),
		tgbotapi.NewInlineKeyboardButtonData("А31ТТ1", "253"),
		tgbotapi.NewInlineKeyboardButtonData("А31ЯР1", "251"),
	),
)

var KeyboarGroup3 = tgbotapi.NewInlineKeyboardMarkup(
	//3курс
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("А21ИСТ1", "227"),
		tgbotapi.NewInlineKeyboardButtonData("А21ИСТ2", "228"),
		tgbotapi.NewInlineKeyboardButtonData("А21МЕФ1", "224"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("А21ПОД1", "226"),
		tgbotapi.NewInlineKeyboardButtonData("А21ЭТЭ", "225"),
		tgbotapi.NewInlineKeyboardButtonData("А21ЯРБ1", "223"),
	),
)

var KeyboarGroup4 = tgbotapi.NewInlineKeyboardMarkup(
	//4курс
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("А11ИСТ1", "191"),
		tgbotapi.NewInlineKeyboardButtonData("А11ИСТ2", "192"),
		tgbotapi.NewInlineKeyboardButtonData("А11МЕФ1", "195"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("А11ПОД1", "196"),
		tgbotapi.NewInlineKeyboardButtonData("А11ЭТЭ1", "193"),
		tgbotapi.NewInlineKeyboardButtonData("А11ЯРБ1", "194"),
	),
)

var KeyboarGroup5 = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		//5курс
		tgbotapi.NewInlineKeyboardButtonData("А01ЯРБ1", "160"),
	),
)

var CourseGroup = map[string]tgbotapi.InlineKeyboardMarkup{
	"1": KeyboarGroup1,
	"2": KeyboarGroup2,
	"3": KeyboarGroup3,
	"4": KeyboarGroup4,
	"5": KeyboarGroup5,
}

func getMondayOfWeek(t time.Time) time.Time {
	offset := int(time.Monday - t.Weekday())
	if offset > 0 {
		offset -= 7
	}
	return t.AddDate(0, 0, offset).Truncate(24 * time.Hour)
}

func CreateMondayKeyboard() tgbotapi.InlineKeyboardMarkup {
	now := time.Now()

	mondayCurrent := getMondayOfWeek(now)

	mondayNext := getMondayOfWeek(now.AddDate(0, 0, 7))

	mondayCurrentFormatted := mondayCurrent.Format("02.01.2006")
	mondayNextFormatted := mondayNext.Format("02.01.2006")

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(mondayCurrentFormatted, mondayCurrentFormatted),
			tgbotapi.NewInlineKeyboardButtonData(mondayNextFormatted, mondayNextFormatted),
		),
	)

	return keyboard
}

type userSelect struct {
	UserID         int64     `json:"user_id"`
	UserSelections [5]string `json:"user_selections"`
	Flag           int
}

var userSelections = make(map[int64]userSelect)

const jsonFilePath = "storage/user_sessions.json"

func SaveToJSON(data map[int64]userSelect, filePath string) error {
	existingData, err := LoadFromJSON(filePath)
	if err != nil {
		return fmt.Errorf("ошибка при чтении файла: %w", err)
	}

	for key, value := range data {
		existingData[key] = value
	}

	jsonData, err := json.MarshalIndent(existingData, "", "  ")
	if err != nil {
		return fmt.Errorf("ошибка при маршалинге данных: %w", err)
	}

	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("ошибка при записи в файл: %w", err)
	}

	return nil
}

func LoadFromJSON(filePath string) (map[int64]userSelect, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении файла: %w", err)
	}

	var data map[int64]userSelect

	err = json.Unmarshal(fileData, &data)
	if err != nil {
		return nil, fmt.Errorf("ошибка при декодировании JSON: %w", err)
	}

	return data, nil

}

func UserExistsAndComplete(userID int64, filePath string) (bool, error) {
	data, err := LoadFromJSON(filePath)
	if err != nil {
		return false, err
	}

	user, exists := data[userID]
	if !exists {
		return false, nil
	}

	for i := 0; i < 5; i++ {
		if user.UserSelections[i] == "" {
			return false, nil
		}
	}

	return true, nil
}

func WriteInfo(bot *tgbotapi.BotAPI, chatID int64) {
	usersel := userSelections[chatID]

	switch usersel.Flag {
	case 0:
		msg := tgbotapi.NewMessage(chatID, "Выберите факультет:")
		msg.ReplyMarkup = KeyboardFac
		bot.Send(msg)
	case 1:
		msg := tgbotapi.NewMessage(chatID, "Выберите форму обучения:")
		msg.ReplyMarkup = KeyboardForm
		bot.Send(msg)
	case 2:
		msg := tgbotapi.NewMessage(chatID, "Выберите курс:")
		msg.ReplyMarkup = KeyboardCourse
		bot.Send(msg)
	case 3:
		msg := tgbotapi.NewMessage(chatID, "Выберите группу:")
		msg.ReplyMarkup = CourseGroup[usersel.UserSelections[2]]
		bot.Send(msg)
	case 4:
		msg := tgbotapi.NewMessage(chatID, "Выберите неделю:")
		msg.ReplyMarkup = CreateMondayKeyboard()
		bot.Send(msg)
	default:
		msg := tgbotapi.NewMessage(chatID, "Вы сделали все необходимые выборы.")
		bot.Send(msg)
	}
}

func HandleCallBack(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.CallbackQuery.Message.Chat.ID
	callbackData := update.CallbackQuery.Data
	messageID := update.CallbackQuery.Message.MessageID

	if _, ok := userSelections[chatID]; !ok {
		userSelections[chatID] = userSelect{UserID: chatID}
	}
	usersel := userSelections[chatID]

	switch usersel.Flag {
	case 0:
		usersel.UserID = chatID
		usersel.UserSelections[0] = callbackData
		usersel.Flag = 1
	case 1:
		usersel.UserID = chatID
		usersel.UserSelections[1] = callbackData
		usersel.Flag = 2
	case 2:
		usersel.UserID = chatID
		usersel.UserSelections[2] = callbackData
		usersel.Flag = 3
	case 3:
		usersel.UserID = chatID
		usersel.UserSelections[3] = callbackData
		usersel.Flag = 4
	case 4:
		usersel.UserID = chatID
		usersel.UserSelections[4] = callbackData
		usersel.Flag = 5

		lastmes := tgbotapi.NewMessage(chatID, "Данные записаны")
		bot.Send(lastmes)
	}

	userSelections[chatID] = usersel

	err := SaveToJSON(userSelections, jsonFilePath)
	if err != nil {
		mes := tgbotapi.NewMessage(chatID, "Ошибка при сохранении данных")
		bot.Send(mes)
		log.Printf("Ошибка при сохранении данных в JSON: %v", err)
	}

	log.Printf("########## +Пользователь %d сделал выбор: %s + ##########", chatID, callbackData)

	deleteMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
	_, err = bot.Send(deleteMsg)
	if err != nil {
		log.Printf("Ошибка при удалении сообщения: %v", err)
	}

	if usersel.Flag != 5 {
		WriteInfo(bot, chatID)
	} else {
		usersel.Flag = 0
		userSelections[chatID] = usersel
		err = SaveToJSON(userSelections, jsonFilePath)
		if err != nil {
			mes := tgbotapi.NewMessage(chatID, "Ошибка при сбросе данных")
			bot.Send(mes)
			log.Printf("Ошибка при сбросе данных в JSON: %v", err)
		}
	}
}

func GetHelp(update tgbotapi.Update, bot tgbotapi.BotAPI) {
	helpMessage := "Доступные команды:\n" +
		"/writeinfo - Введите информацию\n" +
		"/getschedule - Получить расписание\n" +
		"/help - Получить список доступных команд"
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, helpMessage)
	bot.Send(msg)
}

func GetSchedule(update tgbotapi.Update, bot tgbotapi.BotAPI) {

	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	exists, err := UserExistsAndComplete(userID, "storage/user_sessions.json")
	if err != nil {
		log.Printf("Ошибка при проверке пользователя: %v", err)
		msg := tgbotapi.NewMessage(chatID, "Произошла ошибка при проверке пользователя.")
		bot.Send(msg)
		return
	}

	if !exists {
		msg := tgbotapi.NewMessage(chatID, "Введите информацию!")
		bot.Send(msg)
		return
	}

	data, err := LoadFromJSON(jsonFilePath)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Произошла ошибка при загрузке данных пользователя")
		bot.Send(msg)
		return
	}
	faculty := data[userID].UserSelections[0]
	form := data[userID].UserSelections[1]
	course := data[userID].UserSelections[2]
	group := data[userID].UserSelections[3]
	date := data[userID].UserSelections[4] + " 0:00:00"

	schedule, err := scraper.ScrapeSchedule(faculty, form, course, group, date)
	if err != nil {
		log.Printf("Ошибка при получении расписания: %v", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось получить расписание.")
		bot.Send(msg)
		return
	}

	parsedSchedule, err := scraper.ParseSchedule(schedule)
	if err != nil {
		log.Printf("Ошибка при парсинге расписания: %v", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось обработать расписание.")
		bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, parsedSchedule["Понедельник"])
	bot.Send(msg)
	msg = tgbotapi.NewMessage(update.Message.Chat.ID, parsedSchedule["Вторник"])
	bot.Send(msg)
	msg = tgbotapi.NewMessage(update.Message.Chat.ID, parsedSchedule["Среда"])
	bot.Send(msg)
	msg = tgbotapi.NewMessage(update.Message.Chat.ID, parsedSchedule["Четверг"])
	bot.Send(msg)
	msg = tgbotapi.NewMessage(update.Message.Chat.ID, parsedSchedule["Пятница"])
	bot.Send(msg)
	msg = tgbotapi.NewMessage(update.Message.Chat.ID, parsedSchedule["Суббота"])
	bot.Send(msg)
	msg = tgbotapi.NewMessage(update.Message.Chat.ID, parsedSchedule["Воскресенье"])
	bot.Send(msg)
}

func HandleCommands(update tgbotapi.Update, bot tgbotapi.BotAPI) {
	if update.Message != nil {
		switch update.Message.Text {
		case "/help":
			GetHelp(update, bot)
		case "/writeinfo":
			WriteInfo(&bot, update.Message.Chat.ID)
		case "/getschedule":
			GetSchedule(update, bot)
		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда")
			bot.Send(msg)
		}
	}
}
