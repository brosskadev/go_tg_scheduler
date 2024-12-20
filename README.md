# Telegram Schedule Bot

Этот проект представляет собой Telegram-бота, написанного на языке Go, который позволяет пользователям получать расписание учебных занятий. Бот использует библиотеку `tgbotapi` для работы с Telegram API и библиотеку `chromedp` для парсинга расписания с веб-страницы.

## Оглавление
- [Функции](#Функции)
- [Технологии](#Технологии)
- [Установка](#Установка)
- [Использование](#Использование)
- [Команды](#Команды)
- [Примечания](#Примечания)

## Функции
- **Получение расписания**: Бот позволяет пользователям получить расписание занятий на текущую неделю по выбранному факультету, курсу, группе и дате.
- **Обработка пользовательских команд**: Бот реагирует на команды пользователей, например, для запроса расписания или помощи.
- **Парсинг данных с веб-страницы**: Бот использует библиотеку `chromedp` для автоматизации браузера и парсинга расписания с сайта учебного заведения.
  
## Технологии
- Язык: Go
- Библиотека для работы с Telegram API: [tgbotapi](https://github.com/go-telegram-bot-api/telegram-bot-api/v5)
- Парсинг с использованием: [chromedp](https://github.com/chromedp/chromedp)
- Система управления версиями: Git

## Установка

1. Клонируйте репозиторий:
    ```bash
    git clone https://github.com/brosskadev/go_tg_scheduler
    ```

2. Установите зависимости:
    - Убедитесь, что у вас установлен Go и `chromedp`:
    ```bash
    go get github.com/go-telegram-bot-api/telegram-bot-api/v5
    go get github.com/chromedp/chromedp
    ```

3. Настройте Telegram-бота:
    - Получите токен для вашего бота в Telegram.
    - Вставьте свой токен в код, в файле `config/config.go`.

4. Настройте конфигурацию для парсинга:
    - В файле `config/config.go` укажите правильный URL сайта вашего учебного заведения.

## Использование

Запустите бота, используя команду:
```bash
go run main.go
 ```

После этого бот будет готов к работе и начнёт обрабатывать команды пользователей.

## Команды
```bash
/help: Получить список доступных команд.
/writeinfo: Ввод информации (Факультет, форма обучения, курс, группа, дата).
/getschedule: Получить расписание на выбранную неделю.
 ```

## Примечания
Бот работает только с http://rsp.iseu.by/Raspisanie/TimeTable/umu.aspx
