package command

import (
	"log/slog"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/maximotejeda/new-tgb-bot-template/internal/ports"
)

var commandPool *sync.Pool

type Command struct {
	bot    *tgbotapi.BotAPI
	update *tgbotapi.Update
	msg    *tgbotapi.MessageConfig
	log    *slog.Logger
	user   ports.UserService
}

// NewCommand
// Factory for Command Handler
func NewCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update, user ports.UserService) *Command {
	if commandPool == nil {
		commandPool = &sync.Pool{
			New: func() any { return &Command{} },
		}
		for i := 0; i < 20; i++ {
			commandPool.Put(commandPool.New())
		}
	}
	log := slog.Default()
	log = log.With("function", "command", "chat", update.Message.Chat.ID, "userid", update.Message.From.ID, "username", update.Message.From.UserName)
	commands := commandPool.Get().(*Command)
	commands.update = update
	commands.bot = bot
	commands.log = log
	commands.user = user
	return commands
}

// Empty
// Returns pointer to command pool
func (c *Command) Empty() {
	c.update = nil
	c.msg = nil
	c.log = nil
	commandPool.Put(c)
}

// Send
// Process command handlers
func (c *Command) Send() {
	defer c.Empty()
	c.bot.Send(*c.msg)
	del := tgbotapi.NewDeleteMessage(c.update.Message.From.ID, c.update.Message.MessageID)
	c.bot.Send(del)

}

// Handler
// Manage command message for chat
func (c *Command) Handler() {
	msg := tgbotapi.NewMessage(c.update.Message.Chat.ID, "command")
	c.msg = &msg
	//command := c.update.Message.Command()
	c.msg.Text = "message"
	c.Send()
}
