package main

import (
	schedulebot "edushedule2/bot"
	"edushedule2/config"
)

func main() {
	token := config.GetToken()
	schedulebot.InitBot(token)

}
