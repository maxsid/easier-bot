package easierbot

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Bot  is the main object of the bot.
type Bot struct {
	API              *tgbotapi.BotAPI
	Handlers         *MessagesHandlers
	Data             interface{} // this attribute is for data storage
	MessageParseMode string
	updates          tgbotapi.UpdatesChannel
	isWebhook        bool
	listenAddr       string
	webhookSite      string
}

// NewBot is constructor of the Bot structure
func NewBot(token string, debug bool) *Bot {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panicf("Error of api starting: %v\n", err)
	}
	api.Debug = debug
	log.Printf("Authorized on account %s", api.Self.UserName)
	return &Bot{API: api, Handlers: NewMessagesHandlers(), MessageParseMode: "MarkdownV2"}
}

func NewBotViaWebhook(token, webhookSite, listenAddr string, debug bool) *Bot {
	bot := NewBot(token, debug)
	bot.isWebhook = true
	bot.listenAddr = listenAddr
	bot.webhookSite = webhookSite
	return bot
}

// RunBotServer runs listening webhook and updates. It needs running after set handlers.
func (bot *Bot) RunBotServer() {
	if bot.isWebhook {
		bot.setupWebhookMode()
	} else {
		bot.setupPushMode()
	}
	bot.listenUpdates()
}

// setupPushMode setups the bot in the mode when
// a bot server checking messages via repeated requests to telegram server.
func (bot *Bot) setupPushMode() {
	bot.deleteWebhook()
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.API.GetUpdatesChan(u)
	if err != nil {
		log.Panicf("SetupPushMode error: %v\n", err)
	}
	bot.updates = updates
}

// setupWebhookMode starts webserver for Telegram API webhook.
func (bot *Bot) setupWebhookMode() {
	webhookURL := fmt.Sprintf("%s/%s", bot.webhookSite, bot.API.Token)
	_, err := bot.API.SetWebhook(tgbotapi.NewWebhook(webhookURL))
	if err != nil {
		log.Panicf("SetupWebhook error: %v\n", err)
	}
	bot.updates = bot.API.ListenForWebhook("/" + bot.API.Token)
	log.Printf("Starting to listen %s for webhook requests.", bot.listenAddr)
	go func() {
		err := http.ListenAndServe(bot.listenAddr, nil)
		if err != nil {
			log.Panicf("Listening webhook error: %v\n", err)
		}
	}()
	bot.logWebhookInfo()
	log.Println("Webhook listening is started")
}

// deleteWebhook deletes webhook from the telegram api server.
func (bot *Bot) deleteWebhook() {
	resp, err := http.Get("https://api.telegram.org/bot" + bot.API.Token + "/deleteWebhook")
	if err != nil {
		log.Panicf("Delete webhook error: %v\n", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panicf("Error of reading body after a webhook delete: %v\n", err)
	}
	log.Println(string(body))
}

// logWebhookInfo reads the boot webhook info from telegram server and logs result.
func (bot *Bot) logWebhookInfo() {
	info, err := bot.API.GetWebhookInfo()
	if err != nil {
		log.Panic(err)
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}
}

// listenUpdates listens updates from Telegram API and runs handlers
func (bot *Bot) listenUpdates() {
	log.Println("Waiting messages...")
	for update := range bot.updates {
		go bot.Handlers.runHandlerByMessage(bot, update.Message)
	}
}
