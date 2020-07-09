package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	easierbot "github.com/maxsid/easier-bot"
	"os"
)

func main() {
	config := NewConfig(true)
	bot := easierbot.NewBot(config)
	bot.Handlers.AddSeveralCommandsHandler([]string{"hello", "start", "hi"}, func(bot *easierbot.Bot, msg *tgbotapi.Message) {
		bot.SendMessage(msg.Chat.ID, "Hello World!")
	})
	bot.Handlers.AddCommandHandler("bye", func(bot *easierbot.Bot, msg *tgbotapi.Message) {
		bot.SendMessage(msg.Chat.ID, "Ok. Bye.")
	})
	bot.Handlers.AddRegexpHandler(".*ping.*", func(bot *easierbot.Bot, msg *tgbotapi.Message) {
		bot.SendMessage(msg.Chat.ID, "pong")
	})
	bot.RunBotServer()
}

type Config struct {
	token         string
	isDebug       bool
	isWebhookMode bool
	webhookSite   string
	listenAddress string
}

func NewConfig(isDebug bool) *Config {
	return &Config{token: os.Getenv("TOKEN"), isDebug: isDebug, isWebhookMode: true,
		webhookSite: os.Getenv("WEBHOOK_SITE"), listenAddress: os.Getenv("LISTEN_ADDRESS")}
}

func (c Config) GetToken() string {
	return c.token
}

func (c Config) IsDebug() bool {
	return c.isDebug
}

func (c Config) IsWebhookMode() bool {
	return c.isWebhookMode
}

func (c Config) GetWebhookSite() string {
	return c.webhookSite
}

func (c Config) GetListenAddress() string {
	return c.listenAddress
}
