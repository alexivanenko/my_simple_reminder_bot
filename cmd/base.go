package cmd

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

//Command interface
type BaseCmd interface {
	Run(bot *tgbotapi.BotAPI, chatMessage *tgbotapi.Message) error
}

//Base abstract command
type Base struct {
}

func (base *Base) send(bot *tgbotapi.BotAPI, text string, chatId int64, customKeyboard bool) error {
	msg := base.getNewMessage(text, chatId, customKeyboard)

	_, err := bot.Send(msg)

	return err
}

func (base *Base) getNewMessage(text string, chatId int64, customKeyboard bool) tgbotapi.MessageConfig {

	msg := tgbotapi.NewMessage(chatId, text)

	if customKeyboard {
		var buttonsRow1 []tgbotapi.KeyboardButton

		if remindStep.parent.GetStep() == 2 {
			buttonsRow1 = []tgbotapi.KeyboardButton{
				{Text: "Send My Location", RequestLocation: true},
			}
		} else {
			buttonsRow1 = []tgbotapi.KeyboardButton{
				{Text: "/remind"},
				{Text: "/help"},
			}
		}

		msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(buttonsRow1)
	} else {
		msg.ReplyMarkup = tgbotapi.ReplyKeyboardHide{HideKeyboard: true}
	}

	return msg
}

// 'unknown' command
type UnknownCmd struct {
	base *Base
}

func (unknown *UnknownCmd) Run(bot *tgbotapi.BotAPI, chatMessage *tgbotapi.Message) error {
	text := "Unknown command."
	return unknown.base.send(bot, text, chatMessage.Chat.ID, true)
}

// not a command message
type NotCmd struct {
	base *Base
}

func (notCmd *NotCmd) Run(bot *tgbotapi.BotAPI, chatMessage *tgbotapi.Message) error {
	text := "Sorry I don't understand your messages. Please send me a command from list."
	return notCmd.base.send(bot, text, chatMessage.Chat.ID, true)
}

// '/start' command
type StartCmd struct {
	base *Base
}

func (start *StartCmd) Run(bot *tgbotapi.BotAPI, chatMessage *tgbotapi.Message) error {
	help := new(HelpCmd)
	return help.Run(bot, chatMessage)
}

// '/help' command
type HelpCmd struct {
	base *Base
}

func (help *HelpCmd) Run(bot *tgbotapi.BotAPI, chatMessage *tgbotapi.Message) error {
	text := "You can subscribe to notification  using this /remind command."
	return help.base.send(bot, text, chatMessage.Chat.ID, true)
}
