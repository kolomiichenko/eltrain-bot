package main

import (
	"encoding/json"
	"fmt"
	"logger"
	"os"
	"strconv"
	"strings"
	"time"

	railway "github.com/kolomiichenko/swrailway-api"

	"gopkg.in/telegram-bot-api.v4"
)

var (
	bot, _ = tgbotapi.NewBotAPI(config.TelegramBotToken)
	u      = tgbotapi.NewUpdate(0)

	userCache = loadCache()
	config    = loadConfig()
)

type (
	obj map[string]interface{}
	arr []interface{}

	configStruct struct {
		TelegramBotToken   string
		DefaultLanguage    string
		DefaultDisplayMode string
	}
	userCacheStruct struct {
		Language   string
		From       string
		FromCode   string
		To         string
		ToCode     string
		When       string
		WaitAnswer string
	}
)

func loadCache() map[int]userCacheStruct {
	file, err := os.Open("cache.json")
	if err != nil {
		logger.Error.Panic(err.Error())
	}
	decoder := json.NewDecoder(file)
	c := map[int]userCacheStruct{}
	if err := decoder.Decode(&c); err != nil {
		logger.Error.Panic(err.Error())
	}
	return c
}

func saveCache(c map[int]userCacheStruct) {
	file, err := os.OpenFile("cache.json", os.O_RDWR, 0644)
	if err != nil {
		logger.Error.Panic(err.Error())
	}
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(&c); err != nil {
		logger.Error.Panic(err.Error())
	}
}

func loadConfig() configStruct {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	c := configStruct{}
	if err := decoder.Decode(&c); err != nil {
		logger.Error.Panic(err.Error())
	}
	return c
}

func init() {
	if config.DefaultLanguage == "" {
		config.DefaultLanguage = "ua"
	}
	if config.DefaultDisplayMode == "" {
		config.DefaultDisplayMode = "1"
	}

	bot.Debug = false
	logger.Info.Printf("Authorized on account t.me/%s", bot.Self.UserName)
	u.Timeout = 60

	go func() {
		for range time.Tick(1 * time.Minute) {
			saveCache(userCache)
		}
	}()
}

func main() {

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		logger.Error.Panic(err.Error())
	}

	for update := range updates {
		if update.Message == nil && update.InlineQuery != nil { // inline mode
			// query := update.InlineQuery.Query
			inlineConfig := tgbotapi.InlineConfig{
				InlineQueryID: update.InlineQuery.ID,
				IsPersonal:    true,
				CacheTime:     0,
				Results:       arr{},
			}
			if _, err := bot.AnswerInlineQuery(inlineConfig); err != nil {
				logger.Error.Println(err.Error())
			}
		} else { // simple mode

			var command = ""
			if update.Message != nil {

				_, _, fromID := getSpecData(&update)
				c := getUserCache(fromID)

				// set default language
				if c.Language == "" {
					c.Language = config.DefaultLanguage
					userCache[fromID] = c
				}

				// set default display mode
				if c.When == "" {
					c.When = config.DefaultDisplayMode
					userCache[fromID] = c
				}

				command = update.Message.Command()

				if command == "" {

					switch c.WaitAnswer {
					case "typeFrom":
						typeFrom(&update)
					case "typeTo":
						typeTo(&update)
					default:
						sendMarkupMessage(&update, mainMenu(&update), i18n[c.Language]["mainMenuText"])
					}

				} else {

					switch command {

					case "start", "menu":
						sendMarkupMessage(&update, mainMenu(&update), i18n[c.Language]["mainMenuText"])
					case "author":
						sendMarkupMessage(&update, nil, "Author of this bot:\n\n"+
							"Andrii Kolomiichenko\n"+
							"Telegram: @kolomiichenko\n"+
							"Email: bboywilld@gmail.com\n"+
							"Github: <a href=\"https://github.com/kolomiichenko/eltrain-bot\">github.com/kolomiichenko/eltrain-bot</a>")

					case "total":
						sendMarkupMessage(&update, nil, "total users: "+strconv.Itoa(len(userCache)))
					}
				}

			} else {

				if update.CallbackQuery != nil {

					data := strings.Split(update.CallbackQuery.Data, "=")
					key := data[0]
					val := data[1]

					switch key {
					case "mainMenu":
						switch val {
						case "selectLanguage":
							selectLanguage(&update, key, val)
						case "selectFrom":
							selectFrom(&update, key, val)
						case "selectTo":
							selectTo(&update, key, val)
						case "reverse":
							reverse(&update, key, val)
						case "selectWhen":
							selectWhen(&update, key, val)
						case "search":
							search(&update)
						}
					case "setLanguage":
						setLanguage(&update, key, val)
					case "setFrom":
						setFrom(&update, key, val)
					case "setTo":
						setTo(&update, key, val)
					case "setWhen":
						setWhen(&update, key, val)
					}
				}
			}

		}

	}
}

// chatID, messageID, fromID
func getSpecData(update *tgbotapi.Update) (chatID int64, messageID, fromID int) {
	if update.Message != nil {
		chatID = update.Message.Chat.ID
		messageID = update.Message.MessageID
		fromID = update.Message.From.ID
	} else {
		chatID = update.CallbackQuery.Message.Chat.ID
		messageID = update.CallbackQuery.Message.MessageID
		fromID = update.CallbackQuery.From.ID
	}
	return
}

func getUserCache(fromID int) userCacheStruct {
	c := userCacheStruct{}
	c, _ = userCache[fromID]
	return c
}

func send(msg tgbotapi.Chattable) {
	if _, err := bot.Send(msg); err != nil {
		logger.Error.Println(err.Error())
	}
}

func sendMarkupMessage(update *tgbotapi.Update, markup *tgbotapi.InlineKeyboardMarkup, text string) {
	chatID, _, _ := getSpecData(update)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	if markup != nil {
		msg.ReplyMarkup = &markup
	}
	send(msg)
}

func editMarkupMessage(update *tgbotapi.Update, markup *tgbotapi.InlineKeyboardMarkup, text string) {
	chatID, messageID, _ := getSpecData(update)
	msg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	msg.ParseMode = "HTML"
	if markup != nil {
		msg.ReplyMarkup = markup
	}
	send(msg)
}

func createInlineButton(label, key, val string) tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(label, key+"="+val)
}

func mainMenu(update *tgbotapi.Update) *tgbotapi.InlineKeyboardMarkup {
	_, _, fromID := getSpecData(update)
	c := getUserCache(fromID)

	if c.From == "" {
		c.From = "-"
	}
	if c.To == "" {
		c.To = "-"
	}

	keyboard := tgbotapi.InlineKeyboardMarkup{}
	keyboard.InlineKeyboard = [][]tgbotapi.InlineKeyboardButton{
		[]tgbotapi.InlineKeyboardButton{
			createInlineButton(i18n[c.Language]["mainMenuLanguage"]+i18n[c.Language][c.Language], "mainMenu", "selectLanguage"),
		},
		[]tgbotapi.InlineKeyboardButton{
			createInlineButton(i18n[c.Language]["mainMenuFrom"]+c.From, "mainMenu", "selectFrom"),
		},
		[]tgbotapi.InlineKeyboardButton{
			createInlineButton(i18n[c.Language]["mainMenuTo"]+c.To, "mainMenu", "selectTo"),
		},
		[]tgbotapi.InlineKeyboardButton{
			createInlineButton(i18n[c.Language]["mainMenuReverse"], "mainMenu", "reverse"),
		},
		[]tgbotapi.InlineKeyboardButton{
			createInlineButton(i18n[c.Language]["mainMenuWhen"]+i18n[c.Language]["whenMenu"+c.When], "mainMenu", "selectWhen"),
		},
		[]tgbotapi.InlineKeyboardButton{
			createInlineButton(i18n[c.Language]["mainMenuSearch"], "mainMenu", "search"),
		},
	}

	return &keyboard
}

func langMenu(update *tgbotapi.Update) *tgbotapi.InlineKeyboardMarkup {
	_, _, fromID := getSpecData(update)
	c := getUserCache(fromID)

	markup := tgbotapi.InlineKeyboardMarkup{}
	markup.InlineKeyboard = [][]tgbotapi.InlineKeyboardButton{
		[]tgbotapi.InlineKeyboardButton{
			createInlineButton(i18n[c.Language]["ua"], "setLanguage", "ua"),
			createInlineButton(i18n[c.Language]["ru"], "setLanguage", "ru"),
			createInlineButton(i18n[c.Language]["en"], "setLanguage", "en"),
		},
	}
	return &markup
}

func whenMenu(update *tgbotapi.Update) *tgbotapi.InlineKeyboardMarkup {
	_, _, fromID := getSpecData(update)
	c := getUserCache(fromID)

	markup := tgbotapi.InlineKeyboardMarkup{}
	markup.InlineKeyboard = [][]tgbotapi.InlineKeyboardButton{
		[]tgbotapi.InlineKeyboardButton{
			createInlineButton(i18n[c.Language]["whenMenu1"], "setWhen", "1"),
			createInlineButton(i18n[c.Language]["whenMenu0"], "setWhen", "0"),
		},
	}
	return &markup
}

func selectLanguage(update *tgbotapi.Update, key, val string) {
	_, _, fromID := getSpecData(update)
	c := getUserCache(fromID)
	editMarkupMessage(update, langMenu(update), i18n[c.Language]["selectLanguage"])
}

func selectWhen(update *tgbotapi.Update, key, val string) {
	_, _, fromID := getSpecData(update)
	c := getUserCache(fromID)
	editMarkupMessage(update, whenMenu(update), i18n[c.Language]["selectWhen"])
}

func selectFrom(update *tgbotapi.Update, key, val string) {
	_, _, fromID := getSpecData(update)
	c := getUserCache(fromID)

	c.WaitAnswer = "typeFrom"
	userCache[fromID] = c

	editMarkupMessage(update, nil, i18n[c.Language]["typeFrom"])
}

func selectTo(update *tgbotapi.Update, key, val string) {
	_, _, fromID := getSpecData(update)
	c := getUserCache(fromID)

	c.WaitAnswer = "typeTo"
	userCache[fromID] = c

	editMarkupMessage(update, nil, i18n[c.Language]["typeTo"])
}

func reverse(update *tgbotapi.Update, key, val string) {
	_, _, fromID := getSpecData(update)
	c := getUserCache(fromID)

	c.From, c.To = c.To, c.From
	c.FromCode, c.ToCode = c.ToCode, c.FromCode

	userCache[fromID] = c

	editMarkupMessage(update, mainMenu(update), i18n[c.Language]["mainMenuText"])
}

func search(update *tgbotapi.Update) {
	_, _, fromID := getSpecData(update)
	c := getUserCache(fromID)

	if c.Language == "" || c.From == "" || c.To == "" || c.When == "" {
		editMarkupMessage(update, nil, i18n[c.Language]["searchError"])
	} else {

		loc, err := time.LoadLocation("Europe/Kiev")
		if err != nil {
			logger.Error.Panic(err.Error())
		}
		timeNow := time.Now().In(loc).Format("2006-01-02")

		var onlyRemaining bool = c.When != "0"
		shedule := railway.GetShedule(timeNow, "_"+c.Language, c.FromCode, c.ToCode, onlyRemaining)
		results := fmt.Sprintf(i18n[c.Language]["searchTotal"], len(shedule))

		for _, s := range shedule {

			period := ""
			if c.When == "0" {
				period = " [" + s.Period + "]"
			}

			results += "\n🚆 " + s.ID + " " + s.Route + ". " + s.DepartureFrom + "-" + s.ArrivalTo + " (" + s.TimeInTrip + ")" + period
		}

		results += i18n[c.Language]["searchFooter"]

		editMarkupMessage(update, nil, results)
	}
}

func setLanguage(update *tgbotapi.Update, key, val string) {
	_, _, fromID := getSpecData(update)
	c := getUserCache(fromID)

	c.Language = val
	userCache[fromID] = c

	editMarkupMessage(update, mainMenu(update), i18n[c.Language]["mainMenuText"])
}

func setFrom(update *tgbotapi.Update, key, val string) {
	_, _, fromID := getSpecData(update)
	c := getUserCache(fromID)

	station := railway.GetStation(val, "_"+c.Language)

	c.From = station.Label
	c.FromCode = val
	c.WaitAnswer = ""
	userCache[fromID] = c

	editMarkupMessage(update, mainMenu(update), i18n[c.Language]["mainMenuText"])
}

func setTo(update *tgbotapi.Update, key, val string) {
	_, _, fromID := getSpecData(update)
	c := getUserCache(fromID)

	station := railway.GetStation(val, "_"+c.Language)

	c.To = station.Label
	c.ToCode = val
	c.WaitAnswer = ""
	userCache[fromID] = c

	editMarkupMessage(update, mainMenu(update), i18n[c.Language]["mainMenuText"])
}

func setWhen(update *tgbotapi.Update, key, val string) {
	_, _, fromID := getSpecData(update)
	c := getUserCache(fromID)

	c.When = val
	userCache[fromID] = c

	editMarkupMessage(update, mainMenu(update), i18n[c.Language]["mainMenuText"])
}

func typeFrom(update *tgbotapi.Update) {
	_, _, fromID := getSpecData(update)
	c := getUserCache(fromID)

	markup := tgbotapi.InlineKeyboardMarkup{}
	row := []tgbotapi.InlineKeyboardButton{}
	for _, s := range railway.GetStations(update.Message.Text, "_"+c.Language) {
		row = append(row, createInlineButton("🚉 "+s.Label, "setFrom", s.ID))
		markup.InlineKeyboard = append(markup.InlineKeyboard, row)
	}

	sendMarkupMessage(update, &markup, i18n[c.Language]["selectFrom"])
}

func typeTo(update *tgbotapi.Update) {
	_, _, fromID := getSpecData(update)
	c := getUserCache(fromID)

	markup := tgbotapi.InlineKeyboardMarkup{}
	row := []tgbotapi.InlineKeyboardButton{}
	for _, s := range railway.GetStations(update.Message.Text, "_"+c.Language) {
		row = append(row, createInlineButton("🚉 "+s.Label, "setTo", s.ID))
		markup.InlineKeyboard = append(markup.InlineKeyboard, row)
	}

	sendMarkupMessage(update, &markup, i18n[c.Language]["selectTo"])
}

func dumpInterface(in interface{}) {
	b, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		logger.Error.Println(err.Error())
	}
	logger.Info.Print(string(b))
}
