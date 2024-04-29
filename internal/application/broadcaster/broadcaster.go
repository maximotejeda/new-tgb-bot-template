package broadcaster

import (
	"context"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/maximotejeda/new-tgb-bot-template/internal/application/helpers"
	"github.com/maximotejeda/new-tgb-bot-template/internal/ports"
)

type Broadcaster struct {
	ctx     context.Context
	log     *slog.Logger
	user    ports.UserService
	data    []byte
	botName string
}

type Message struct {
	Message string `json:"message"`
	Data    string `json:"data"`
	Error   error  `json:"error"`
}

func NewBroadCast(ctx context.Context, user ports.UserService, data []byte) *Broadcaster {
	log := slog.Default()
	log = log.With("place", "bradcast")
	return &Broadcaster{
		ctx:  ctx,
		user: user,
		data: data,
		log:  log,
	}
}

func (b *Broadcaster) SendList() []tgbotapi.MessageConfig {
	// convert data to map
	//m := Message{}
	listMsg := []tgbotapi.MessageConfig{}

	/*		for _, userID := range userList {
				msg := tgbotapi.NewMessage(userID, text)
				msg.ReplyMarkup = keyboard
				listMsg = append(listMsg, msg)

			}
		}*/
	return listMsg
}

func (b Broadcaster) SendAllUsers(ctx context.Context, log *slog.Logger, data []byte, botname string) []tgbotapi.MessageConfig {
	userList, err := b.user.GetAllBotsUsers(botname)
	b.log.Info("broadcast", "user list", userList)
	msgs := []tgbotapi.MessageConfig{}
	if err != nil {
		return msgs
	}
	cancelBTN := map[string]string{}
	cancelBTN["Eliminar ❌"] = "cancelar=true"
	keyboard := helpers.CreateKeyboard(cancelBTN)
	for _, user := range userList {
		msg := tgbotapi.NewMessage(user.TguID, "")
		msg.Text = string(data)
		msg.ReplyMarkup = keyboard
		msgs = append(msgs, msg)
	}
	return msgs
}

func comparer(before, after float64) string {
	if before > after {
		return "⬇️"
	} else if before < after {
		return "⬆️"
	} else {
		return "🟰"
	}

}
