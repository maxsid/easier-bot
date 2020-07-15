package easierbot

import (
	"github.com/maxsid/easier-bot/contentType"
	"reflect"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// MessagesHandlers contains all handlers of messages
type MessagesHandlers struct {
	commandHandlers map[string]MessageHandler // contains handlers which calls with an exact command match
	regexpHandlers  []*regexpMessageHandler
	contentHandlers map[contentType.ContentType]MessageHandler
	defaultHandler  MessageHandler
}

func NewMessagesHandlers() *MessagesHandlers {
	return &MessagesHandlers{
		commandHandlers: make(map[string]MessageHandler),
		regexpHandlers:  make([]*regexpMessageHandler, 0),
		contentHandlers: make(map[contentType.ContentType]MessageHandler),
	}
}

// runHandlerByMessage checks a message and execute an appropriate handler
func (h *MessagesHandlers) runHandlerByMessage(bot *Bot, msg *tgbotapi.Message) {
	if msg == nil {
		return
	}
	defer panicRecover()
	// check command
	if msg.IsCommand() {
		lowerCommand := strings.ToLower(msg.Command())
		handler, ok := h.commandHandlers[lowerCommand]
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
	// content handlers
	if ct := getContentType(h, msg); ct != contentType.NoContent {
		h.contentHandlers[ct](bot, msg)
		return
	}
	// other handlers
	if h.defaultHandler != nil {
		h.defaultHandler(bot, msg)
	}
}

// getContentType looks for not text content in a message and returns its type.
func getContentType(h *MessagesHandlers, msg *tgbotapi.Message) contentType.ContentType {
	elemValue := reflect.ValueOf(msg).Elem()
	typeOf := elemValue.Type()
	canBeAny := false
	for i := 0; i < elemValue.NumField(); i++ {
		ct := contentType.ContentType(typeOf.Field(i).Name)
		// check field in all possible types and the filed value is not nil
		if _, hasInAll := contentType.AllContentTypes[ct]; !hasInAll || elemValue.Field(i).IsNil() {
			continue
		}
		// if it went here, it can have any content
		canBeAny = true
		// check field in set content handlers and
		if _, hasInHandlers := h.contentHandlers[ct]; hasInHandlers {
			return ct
		}
	}
	if _, hasAnyInHandlers := h.contentHandlers[contentType.Any]; hasAnyInHandlers && canBeAny {
		return contentType.Any
	}
	return contentType.NoContent
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
	lowerCommand := strings.ToLower(command)
	h.commandHandlers[lowerCommand] = handler
}

// AddRegexpHandler sets a handler under a regular expression.
// The handler runs when the update listener got this text which has a regexp match.
func (h *MessagesHandlers) AddRegexpHandler(expr string, handler MessageHandler) {
	h.regexpHandlers = append(h.regexpHandlers, newRegexpMessageHandler(expr, handler))
}

// SetContentHandler sets handler which will be executed if a message has another type of content.
func (h *MessagesHandlers) SetContentHandler(contentType contentType.ContentType, handler MessageHandler) {
	h.contentHandlers[contentType] = handler
}

// SetDefaultHandler sets handler which will be executed in other cases.
func (h *MessagesHandlers) SetDefaultHandler(handler MessageHandler) {
	h.defaultHandler = handler
}

// MessageHandler is a function for handling a bot message
type MessageHandler func(bot *Bot, msg *tgbotapi.Message)

// regexpMessageHandler provides to contain regex and handler
type regexpMessageHandler struct {
	regexp  *regexp.Regexp
	handler MessageHandler
}

// newRegexpMessageHandler is constructor of regexpMessageHandler
func newRegexpMessageHandler(expr string, handler MessageHandler) *regexpMessageHandler {
	return &regexpMessageHandler{regexp.MustCompile(expr), handler}
}
