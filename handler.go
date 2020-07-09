package easierbot

import (
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// MessagesHandlers contains all handlers of messages
type MessagesHandlers struct {
	commandHandlers map[string]MessageHandler // contains handlers which calls with an exact command match
	regexpHandlers  []regexpMessageHandler
}

func NewMessagesHandlers() *MessagesHandlers {
	return &MessagesHandlers{
		commandHandlers: make(map[string]MessageHandler),
		regexpHandlers:  make([]regexpMessageHandler, 0),
	}
}

// checkCall checks a message and execute an appropriate handler
func (h *MessagesHandlers) checkCall(bot *Bot, msg *tgbotapi.Message) {
	if msg == nil {
		return
	}
	// check command
	if msg.IsCommand() {
		handler, ok := h.commandHandlers[msg.Command()]
		if ok {
			handler(bot, msg)
			return
		}
	}
	// check regular expressions
	lowerText := strings.ToLower(msg.Text)
	for _, v := range h.regexpHandlers {
		if v.regexp.MatchString(lowerText) {
			v.handler(bot, msg)
			return
		}
	}
}

// AddSeveralCommandsHandler sets a handler under several commands.
// The handler runs when the update listener got this commands
func (h *MessagesHandlers) AddSeveralCommandsHandler(commands []string, handler MessageHandler) {
	for _, command := range commands {
		h.AddCommandHandler(command, handler)
	}
}

// AddCommandHandler sets a handler under a command.
// The handler runs when the update listener got this command.
func (h *MessagesHandlers) AddCommandHandler(command string, handler MessageHandler) {
	h.commandHandlers[command] = handler
}

// AddRegexpHandler sets a handler under a regular expression.
// The handler runs when the update listener got this text which has a regexp match.
func (h *MessagesHandlers) AddRegexpHandler(expr string, handler MessageHandler) {
	expr = strings.ToLower(expr)
	reMsgHandler := regexpMessageHandler{regexp.MustCompile(expr), handler}
	h.regexpHandlers = append(h.regexpHandlers, reMsgHandler)
}

// MessageHandler is a function for handling a bot message
type MessageHandler func(bot *Bot, msg *tgbotapi.Message)

// regexpMessageHandler provides to contain regex and handler
type regexpMessageHandler struct {
	regexp  *regexp.Regexp
	handler MessageHandler
}
