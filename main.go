package main

import (
	"log"

	"github.com/alexivanenko/my_simple_reminder_bot/cmd"
	"github.com/alexivanenko/my_simple_reminder_bot/config"
	"github.com/alexivanenko/my_simple_reminder_bot/model"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	defer model.GetSession().Close()

	bot, err := tgbotapi.NewBotAPI(config.String("bot", "token"))
	if err != nil {
		log.Panic(err)
	}

	if config.Is("bot", "debug") {
		bot.Debug = true
	} else {
		bot.Debug = false
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		cmd.Run(bot, update.Message)
	}

}
