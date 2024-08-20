package main

import (
	"nektome-discord/bot"
	// "log"
	"os"
)

func main() {
	botToken, ok := os.LookupEnv(`BOT_TOKEN`)
	if !ok {
		panic(`Не определён BOT_TOKEN`)
	}

	bot.BotToken = botToken
	bot.Run()
}
