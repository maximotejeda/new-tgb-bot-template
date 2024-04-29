package message

import (
	"log/slog"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/maximotejeda/new-tgb-bot-template/internal/application/static"
	"github.com/maximotejeda/new-tgb-bot-template/internal/ports"
)

var ChatPool *sync.Pool

type Message struct {
	bot    *tgbotapi.BotAPI
	update *tgbotapi.Update
	msg    *tgbotapi.MessageConfig
	log    *slog.Logger
	user   ports.UserService
}

// NewMessage
// Factory for message handler
func NewMessage(bot *tgbotapi.BotAPI, update *tgbotapi.Update, user ports.UserService) *Message {
	if ChatPool == nil {
		ChatPool = &sync.Pool{
			New: func() any { return &Message{} },
		}
		for i := 0; i < 20; i++ {
			ChatPool.Put(ChatPool.New())
		}
	}
	log := slog.Default()
	log = log.With("function", "message", "chat", update.Message.Chat.ID, "userid", update.Message.From.ID, "username", update.Message.From.UserName)
	message := ChatPool.Get().(*Message)
	message.update = update
	message.bot = bot
	message.log = log
	message.user = user
	return message
}

// Empty
// Returns pointer to pool
func (m *Message) Empty() {
	m.update = nil
	m.msg = nil
	m.log = nil
	m.user = nil
	ChatPool.Put(m)
}

// Send
// Process message sending to bot
func (m *Message) Send() {
	defer m.Empty()
	m.bot.Send(m.msg)
}

// Handler
// Manage features for messages
func (m *Message) Handler() {
	msg := tgbotapi.NewMessage(m.update.Message.Chat.ID, "")
	m.msg = &msg
	msgtext := static.RemoveAccent(m.update.Message.Text)
	m.msg.Text = msgtext
	m.Send()
}
