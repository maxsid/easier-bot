package easierbot

// Configurator is an interface with necessary data for the bot working.
type Configurator interface {
	GetToken() string
	IsDebug() bool
	IsWebhookMode() bool
	GetWebhookSite() string
	GetListenAddress() string
}
