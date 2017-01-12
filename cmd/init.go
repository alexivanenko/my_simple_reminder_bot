package cmd

import (
	"strconv"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var runList map[string]BaseCmd
var ChatUserID string
var remindStep *RemindStep

func init() {
	runList = make(map[string]BaseCmd)
	var command interface{}

	//System Commands
	command = new(UnknownCmd)
	runList["unknown"] = command.(BaseCmd)
	command = new(NotCmd)
	runList["not_cmd"] = command.(BaseCmd)

	//Base Commands
	command = new(StartCmd)
	runList["start"] = command.(BaseCmd)
	command = new(HelpCmd)
	runList["help"] = command.(BaseCmd)

	//Remind Commands
	//command = new(RemindCmd)
	remind := NewRemindCmd()
	remindStep = NewRemindStep(remind)

	command = remind
	runList["remind"] = command.(BaseCmd)

	command = remindStep
	runList["remind_step"] = command.(BaseCmd)
}

func Run(bot *tgbotapi.BotAPI, chatMessage *tgbotapi.Message) error {
	command := "unknown"

	ChatUserID = strconv.Itoa(chatMessage.From.ID)

	if chatMessage.IsCommand() {

		//Check is command in list
		if _, ok := runList[chatMessage.Command()]; ok {
			command = chatMessage.Command()
		}

	} else {

		//Not a Telegram bot command
		command = "not_cmd"

		//Remind Steps
		if remindStep.parent.GetStep() > 0 {
			command = "remind_step"
		}

	}

	return runList[command].Run(bot, chatMessage)
}
