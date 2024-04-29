package user

import (
	"context"
	"log/slog"

	tgb "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/maximotejeda/msvc-proto/golang/tgbuser"
	"github.com/maximotejeda/new-tgb-bot-template/internal/application/domain"
	"google.golang.org/grpc"
)

type Adapter struct {
	user tgbuser.UserManagerClient
	conn *grpc.ClientConn
	log  *slog.Logger
}

func NewAdapter(conn *grpc.ClientConn) (*Adapter, error) {
	log := slog.Default()
	log = log.With("location", "user adapter")
	client := tgbuser.NewUserManagerClient(conn)
	return &Adapter{user: client, conn: conn, log: log}, nil
}

func (a *Adapter) Get(tgbid int64) (*domain.User, error) {
	hr, err := a.user.Get(context.Background(), &tgbuser.GetTGBUserRequest{TgbId: tgbid})
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		ID:      hr.User.Id,
		TguID:   hr.User.TgbId,
		Created: hr.User.Created,
		Edited:  hr.User.Edited,
	}
	return user, nil
}

func (a Adapter) Create(user *tgb.User) (b bool, err error) {
	_, err = a.user.Create(context.Background(), &tgbuser.CreateTGBUserRequest{
		User: &tgbuser.User{
			TgbId:     user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Username:  user.UserName,
		},
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (a Adapter) Edit(user *tgb.User) (b bool, err error) {

	_, err = a.user.Edit(context.Background(), &tgbuser.EditTGBUserRequest{
		User: &tgbuser.User{
			Username:  user.UserName,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		},
	})
	if err != nil {
		return false, err
	}

	return true, nil
}
func (a Adapter) Delete(tgbid int64) (b bool, err error) {
	_, err = a.user.Delete(context.Background(), &tgbuser.DeleteTGBUserRequest{
		TgbId: tgbid,
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

func (a Adapter) GetBots(tgbid int64) (s []string, err error) {
	hr, err := a.user.GetBots(context.Background(), &tgbuser.GetBotsTGBUserRequest{
		TgbId: tgbid,
	})
	if err != nil {
		return nil, err
	}
	s = []string{}

	if len(hr.Bots) <= 0 {
		return s, nil
	}

	for _, it := range hr.Bots {
		s = append(s, it.BotName)
	}
	return s, nil
}

func (a Adapter) AddBot(tgbid int64, botname string) (b bool, err error) {
	_, err = a.user.AddBot(context.Background(), &tgbuser.AddBotTGBUserRequest{
		TgbId:   tgbid,
		BotName: botname,
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (a Adapter) DeleteBot(tgbid int64, botname string) (b bool, err error) {
	_, err = a.user.DeleteBot(context.Background(), &tgbuser.DeleteBotTGBUserRequest{
		TgbId:   tgbid,
		BotName: botname,
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (a Adapter) GetAllBotsUsers(botname string) ([]*domain.User, error) {
	users, err := a.user.GetAllBotsUsers(context.Background(), &tgbuser.GetAllBotsUsersRequest{BotName: botname})
	if err != nil {
		a.log.Error("get all bots users", "error", err)
		return nil, err
	}
	a.log.Info("users", "result", users)
	list := []*domain.User{}
	for _, us := range users.Users {
		user := &domain.User{
			ID:        us.Id,
			TguID:     us.TgbId,
			Username:  us.Username,
			FirstName: us.FirstName,
			LastName:  us.LastName,
			Created:   us.Created,
			Edited:    us.Edited,
		}
		list = append(list, user)
	}
	return list, nil
}
