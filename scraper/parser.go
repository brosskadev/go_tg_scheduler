package scraper

import (
	"strings"
	"time"
)

func ParseSchedule(schedule string) (map[string]string, error) {
	// Создаем словарь для хранения расписания по дням недели
	parsedData := make(map[string]string)

	// Указываем все дни недели
	daysOfWeek := []string{
		"Понедельник", "Вторник", "Среда", "Четверг", "Пятница", "Суббота", "Воскресенье",
	}

	lines := strings.Split(schedule, "\n")
	var currentDay string
	var currentSchedule string

	// Проходим по всем строкам
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Если строка начинается с дня недели, добавляем текущий день в карту
		for _, day := range daysOfWeek {
			if strings.HasPrefix(line, day) {
				// Если был найден предыдущий день, сохраняем его расписание
				if currentDay != "" {
					parsedData[currentDay] = currentSchedule
				}
				// Обновляем текущий день и очищаем расписание для нового дня
				currentDay = day
				currentSchedule = day + "\n"
				break
			}
		}

		if currentDay != "" && !strings.HasPrefix(line, currentDay) {
			// Если строка содержит дату, добавляем пустую строку после нее
			if _, err := time.Parse("02.01.2006", line); err == nil { // Проверка на формат даты
				currentSchedule += line + "\n\n" // Добавляем пустую строку после даты
			} else if strings.Contains(line, "ауд.") { // Проверка на наличие аудитории
				currentSchedule += line + "\n\n" // Пустая строка после занятия
			} else {
				currentSchedule += line + "\n" // Обычная строка
			}
		}
	}

	// Сохраняем последний день
	if currentDay != "" {
		parsedData[currentDay] = currentSchedule
	}

	return parsedData, nil
}
