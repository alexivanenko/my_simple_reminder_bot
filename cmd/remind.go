package cmd

import (
	"errors"
	"time"

	"github.com/alexivanenko/my_simple_reminder_bot/model"
	"github.com/bradfitz/latlong"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/patrickmn/go-cache"
)

const DATE_FORMAT = "02/01/2006 15:04:05"

// '/remind' command
type RemindCmd struct {
	base  *Base
	cache *cache.Cache
}

func NewRemindCmd() *RemindCmd {
	return &RemindCmd{cache: cache.New(10*time.Minute, 30*time.Second)}
}

func (remind *RemindCmd) Run(bot *tgbotapi.BotAPI, chatMessage *tgbotapi.Message) error {
	text := "What is the name of the event?"
	err := remind.base.send(bot, text, chatMessage.Chat.ID, false)

	//Set Step 1
	if err == nil {
		remind.setStep(1)
	}

	return err
}

func (remind *RemindCmd) GetStep() int {
	result := 0
	step, found := remind.cache.Get("REMIND_STEP_" + ChatUserID)

	if found {
		result = step.(int)
	}

	return result
}

func (remind *RemindCmd) setStep(step int) {
	remind.cache.Set("REMIND_STEP_"+ChatUserID, step, cache.DefaultExpiration)
}

func (remind *RemindCmd) resetStep() {
	remind.cache.Delete("REMIND_STEP_" + ChatUserID)
}

//------------------------------------REMIND STEP-----------------------------------------------------

type RemindStep struct {
	parent *RemindCmd
}

func NewRemindStep(parent *RemindCmd) *RemindStep {
	return &RemindStep{parent: parent}
}

func (step *RemindStep) Run(bot *tgbotapi.BotAPI, chatMessage *tgbotapi.Message) error {

	currentStep := step.parent.GetStep()
	event := step.getTempEvent()

	param := chatMessage.Text
	var msgText string

	customKeyboard := true

	//Validate and store name
	if currentStep == 1 {
		if step.isValidName(param) {

			event.Name = param
			event.ChatID = chatMessage.Chat.ID
			step.storeEventTemporary(event)

			step.parent.setStep(2)
			msgText = "Please share your location so that we know your TimeZone for notifications."

		} else {
			customKeyboard = false
			msgText = "The event name should not be empty. Please enter again."

		}
		//Store location timezone
	} else if currentStep == 2 {
		customKeyboard = false

		tz := latlong.LookupZoneName(chatMessage.Location.Latitude, chatMessage.Location.Longitude)
		_, err := time.LoadLocation(tz)

		if err != nil {
			tz = time.Now().Location().String()
		}

		event.TimeZone = tz
		step.storeEventTemporary(event)

		step.parent.setStep(3)
		msgText = "When I should notify you about " + event.Name + "? (please use this format: dd/mm/yyyy H:m)"
		//Validate and store datetime
	} else if currentStep == 3 {

		date, err := step.parseDate(param, event.TimeZone)

		if err == nil {

			event.Date = date
			event.Save()

			step.parent.resetStep()
			msgText = "You will be notify about " + event.Name + " at " + date.Format(DATE_FORMAT)

		} else {

			customKeyboard = false
			msgText = "Please enter valid date time in dd/mm/yyyy H:m format. And the date should be at least 3 minutes after now."

		}
	}

	return step.parent.base.send(bot, msgText, chatMessage.Chat.ID, customKeyboard)
}

func (step *RemindStep) isValidName(name string) bool {
	if name != "" {
		return true
	} else {
		return false
	}
}

func (step *RemindStep) parseDate(dateStr string, timezone string) (time.Time, error) {
	var err error
	var location *time.Location
	var date time.Time

	location, err = time.LoadLocation(timezone)

	if err != nil {
		location = time.Now().Location()
	}

	date, err = time.ParseInLocation(DATE_FORMAT, dateStr+":00", location)

	if err == nil {
		now := time.Now().In(location).Add(2 * time.Minute)

		if date.Before(now) {
			err = errors.New("Past Date")
		}
	}

	return date, err
}

func (step *RemindStep) storeEventTemporary(event *model.Event) {
	step.parent.cache.Set("REMIND_EVENT_"+ChatUserID, event, cache.DefaultExpiration)
}

func (step *RemindStep) getTempEvent() *model.Event {
	event, found := step.parent.cache.Get("REMIND_EVENT_" + ChatUserID)

	if found {
		return event.(*model.Event)
	} else {
		return new(model.Event)
	}
}
