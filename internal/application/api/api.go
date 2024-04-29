package api

import (
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/maximotejeda/new-tgb-bot-template/internal/application/command"
	"github.com/maximotejeda/new-tgb-bot-template/internal/application/message"
	"github.com/maximotejeda/new-tgb-bot-template/internal/application/query"
	"github.com/maximotejeda/new-tgb-bot-template/internal/ports"
)

type api struct {
	bot     *tgbotapi.BotAPI
	command ports.Tgb
	message ports.Tgb
	query   ports.Tgb
	user    ports.UserService
	log     *slog.Logger
}

func NewApi(bot *tgbotapi.BotAPI) *api {
	log := slog.Default()
	log = log.With("location", "root")
	return &api{bot: bot, log: log}
}

func (a *api) Run(update *tgbotapi.Update, user ports.UserService) {
	us, err := user.Get(update.SentFrom().ID)
	if err != nil {
		a.log.Error("geting user", "id", update.SentFrom().ID)
	}
	a.log.Info("geted user", "user", us)
	msg := update.Message
	if msg != nil { // message is not nil can be a command or a text message
		if msg.IsCommand() {
			a.command = command.NewCommand(a.bot, update, user)
			a.command.Handler()
			// is a command
		} else if msg.Text != "" {
			// is a text message
			a.message = message.NewMessage(a.bot, update, user)
			a.message.Handler()
		}
	} else if update.CallbackQuery != nil {
		// is a cal back query
		a.query = query.NewQuery(a.bot, update, user)
		a.query.Handler()
	}
}
