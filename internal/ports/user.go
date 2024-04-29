package ports

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/maximotejeda/new-tgb-bot-template/internal/application/domain"
)

type UserService interface {
	Get(int64) (*domain.User, error)
	Edit(*tgbotapi.User) (bool, error)
	Delete(int64) (bool, error)
	Create(*tgbotapi.User) (bool, error)
	AddBot(int64, string) (bool, error)
	GetBots(int64) ([]string, error)
	DeleteBot(int64, string) (bool, error)
	GetAllBotsUsers(string) ([]*domain.User, error)
}
