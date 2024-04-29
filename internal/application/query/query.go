package query

import (
	"log/slog"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/maximotejeda/new-tgb-bot-template/internal/ports"
)

var chatPool *sync.Pool

type Query struct {
	bot    *tgbotapi.BotAPI
	update *tgbotapi.Update
	msg    *tgbotapi.MessageConfig
	log    *slog.Logger
	user   ports.UserService
}

// NewQuery
// Factory for query handlers
func NewQuery(bot *tgbotapi.BotAPI, update *tgbotapi.Update, user ports.UserService) *Query {
	if chatPool == nil {
		chatPool = &sync.Pool{
			New: func() any { return &Query{} },
		}
		for i := 0; i < 20; i++ {
			chatPool.Put(chatPool.New())
		}
	}
	log := slog.Default()
	log = log.With("function", "query", "chat", update.CallbackQuery.From.ID, "userid", update.CallbackQuery.From.ID, "username", update.CallbackQuery.From.UserName)
	query := chatPool.Get().(*Query)
	query.update = update
	query.bot = bot
	query.log = log

	query.user = user
	return query
}

// Empty
// Returns pointer to pool
func (q *Query) Empty() {
	q.update = nil
	q.msg = nil
	q.log = nil
	q.user = nil
	chatPool.Put(q)
}

// Send
// Process Query message
func (q *Query) Send() {
	defer q.Empty()
	q.bot.Send(q.msg)
	// Delete previous message
	del := tgbotapi.NewDeleteMessage(q.update.CallbackQuery.From.ID, q.update.CallbackQuery.Message.MessageID)
	q.bot.Send(del)

}

// Handler
// Manage query message
func (q *Query) Handler() {
	msg := tgbotapi.NewMessage(q.update.CallbackQuery.Message.Chat.ID, "")
	q.msg = &msg
	q.msg.Text = "message"
}
