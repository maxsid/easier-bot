package easierbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/maxsid/easier-bot/contentType"
	"testing"
)

func Test_getContentType(t *testing.T) {
	emptyHandler := func(bot *Bot, msg *tgbotapi.Message) {}
	defH := &MessagesHandlers{contentHandlers: map[contentType.ContentType]MessageHandler{
		contentType.Location: emptyHandler,
		contentType.Document: emptyHandler,
		contentType.Any:      emptyHandler,
	}}
	type args struct {
		h   *MessagesHandlers
		msg *tgbotapi.Message
	}
	tests := []struct {
		name string
		args args
		want contentType.ContentType
	}{
		{
			name: "Empty",
			args: args{
				h:   defH,
				msg: &tgbotapi.Message{},
			},
			want: contentType.NoContent,
		},
		{
			name: "Found Location",
			args: args{
				h:   defH,
				msg: &tgbotapi.Message{Location: &tgbotapi.Location{}},
			},
			want: contentType.Location,
		},
		{
			name: "Found Document",
			args: args{
				h:   defH,
				msg: &tgbotapi.Message{Document: &tgbotapi.Document{}},
			},
			want: contentType.Document,
		},
		{
			name: "Found Any",
			args: args{
				h:   defH,
				msg: &tgbotapi.Message{Sticker: &tgbotapi.Sticker{}},
			},
			want: contentType.Any,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getContentType(tt.args.h, tt.args.msg); got != tt.want {
				t.Errorf("getContentType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessagesHandlers_checkCall(t *testing.T) {
	var lastResult = -1
	getHandler := func(result int) MessageHandler {
		return func(bot *Bot, msg *tgbotapi.Message) {
			lastResult = result
		}
	}

	defH := MessagesHandlers{
		commandHandlers: map[string]MessageHandler{
			"test1": getHandler(1),
			"test2": getHandler(2),
		},
		regexpHandlers: []*regexpMessageHandler{
			newRegexpMessageHandler(`^\d+$`, getHandler(3)),
			newRegexpMessageHandler("hello", getHandler(4)),
			newRegexpMessageHandler("[Bb]uy", getHandler(5)),
			newRegexpMessageHandler("/test1", getHandler(6)),
		},
		contentHandlers: map[contentType.ContentType]MessageHandler{
			contentType.Document: getHandler(7),
			contentType.Location: getHandler(8),
			contentType.Any:      getHandler(9),
		},
		defaultHandler: getHandler(10),
	}
	tests := []struct {
		name     string
		handlers MessagesHandlers
		msg      *tgbotapi.Message
		want     int
	}{
		{
			name:     "Nil Message",
			handlers: defH,
			msg:      nil,
			want:     -1,
		},
		{
			name:     "Default Handler",
			handlers: defH,
			msg: &tgbotapi.Message{Text: "/toto", Entities: &[]tgbotapi.MessageEntity{
				{
					Type:   "bot_command",
					Length: 5,
				},
			}},
			want: 10,
		},
		{
			name:     "Test1 Command",
			handlers: defH,
			msg: &tgbotapi.Message{Text: "/test1", Entities: &[]tgbotapi.MessageEntity{
				{
					Type:   "bot_command",
					Length: 6,
				},
			}},
			want: 1,
		},
		{
			name:     "Test2 Command",
			handlers: defH,
			msg: &tgbotapi.Message{Text: "/test2", Entities: &[]tgbotapi.MessageEntity{
				{
					Type:   "bot_command",
					Length: 6,
				},
			}},
			want: 2,
		},
		{
			name:     "Test1 Multicase Command",
			handlers: defH,
			msg: &tgbotapi.Message{Text: "/TeSt1", Entities: &[]tgbotapi.MessageEntity{
				{
					Type:   "bot_command",
					Length: 6,
				},
			}},
			want: 1,
		},
		{
			name:     "Regex Number",
			handlers: defH,
			msg:      &tgbotapi.Message{Text: "123"},
			want:     3,
		},
		{
			name:     "Regex 'hello'",
			handlers: defH,
			msg:      &tgbotapi.Message{Text: "So, hello man!"},
			want:     4,
		},
		{
			name:     "Regex multicase",
			handlers: defH,
			msg:      &tgbotapi.Message{Text: "So, BUY man!"},
			want:     5,
		},
		{
			name:     "Regex '/test1' in text",
			handlers: defH,
			msg:      &tgbotapi.Message{Text: "So, enter /test1 command, man!"},
			want:     6,
		},
		{
			name:     "Msg with Document",
			handlers: defH,
			msg:      &tgbotapi.Message{Document: &tgbotapi.Document{}},
			want:     7,
		},
		{
			name:     "Msg with Location",
			handlers: defH,
			msg:      &tgbotapi.Message{Location: &tgbotapi.Location{}},
			want:     8,
		},
		{
			name:     "Msg with any content",
			handlers: defH,
			msg:      &tgbotapi.Message{Audio: &tgbotapi.Audio{}},
			want:     9,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lastResult = -1
			tt.handlers.checkCall(nil, tt.msg)
			if tt.want != lastResult {
				t.Errorf("lastResult = %v, want %v", lastResult, tt.want)
			}
		})
	}
}
