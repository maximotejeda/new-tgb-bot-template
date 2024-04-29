package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	tb "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/maximotejeda/new-tgb-bot-template/config"
	"github.com/maximotejeda/new-tgb-bot-template/internal/adapters/user"
	"github.com/maximotejeda/new-tgb-bot-template/internal/application/api"
	"github.com/maximotejeda/new-tgb-bot-template/internal/application/broadcaster"
	"github.com/maximotejeda/new-tgb-bot-template/internal/ports"
	"github.com/nats-io/nats.go"
	"golang.org/x/sync/semaphore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	maxWorkers = runtime.GOMAXPROCS(0)
	sem        = semaphore.NewWeighted(int64(maxWorkers) * 2)
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	log = log.With("location", "main")
	nc, _ := nats.Connect(config.GetNatsURI())
	ctx := context.Background()

	bot, err := tb.NewBotAPI(config.GetToken())
	if err != nil {
		log.Error("token not found", "error", err)
		panic(err)
	}

	botName := bot.Self.UserName

	bot.Debug = config.GetEnvironment() == "development"
	log.Info("Bot Authorized", "username", botName)
	log.Info("Initiated with a concurrency limit", "max concurrency", maxWorkers*2)
	u := tb.NewUpdate(0)
	u.Timeout = 60

	// bot user update channel
	updtChan := bot.GetUpdatesChan(u)
	// subs chann
	changeChan := make(chan *nats.Msg, 64)
	broadcastChan := make(chan *nats.Msg, 64)
	defer close(changeChan)
	defer close(broadcastChan)
	sub, err := nc.ChanSubscribe("dolar-crawler", changeChan)
	if err != nil {
		log.Error("subscribing", "error", err.Error())
	}
	info, err := nc.ChanSubscribe("dolar-bot", broadcastChan)
	if err != nil {
		log.Error("subscribing", "error", err.Error())
	}
	defer sub.Drain()
	defer info.Drain()
	defer nc.Close()

	// exit channel
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)
	defer close(sign)
	app := api.NewApi(bot)
	for {
		select {
		case update := <-updtChan:
			if err = sem.Acquire(ctx, 1); err != nil {
				bot.Send(tb.NewMessage(update.FromChat().ID, "error adquiring update"))
				continue
			}
			go func() {
				defer sem.Release(1)
				user, userConn := CreateAdaptersGRPC()
				app.Run(&update, user)
				userConn.Close()
			}()
		case message := <-changeChan:
			user, userConn := CreateAdaptersGRPC()

			bcast := broadcaster.NewBroadCast(ctx, user, message.Data)
			userList := bcast.SendList()

			for _, msg := range userList {
				go bot.Send(msg)
			}

			userConn.Close()

		case message := <-broadcastChan:
			user, userConn := CreateAdaptersGRPC()

			bcast := broadcaster.NewBroadCast(ctx, user, message.Data)
			msgs := bcast.SendAllUsers(ctx, log, message.Data, bot.Self.UserName)
			log.Info("broadcast", "data", string(message.Data), "msg", msgs)
			for _, msg := range msgs {
				go bot.Send(msg)
			}
			userConn.Close()
		case <-sign:
			log.Error("killing app due to syscall ")
			os.Exit(1)
		}
	}
}
func CreateAdaptersGRPC() (ports.UserService, *grpc.ClientConn) {
	log := slog.Default()
	// we are outside update so we will be querying db to
	// get users interested in specific updates ex bpd, brd, apa
	// userID inst=> comma separated string
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	userConn, err := grpc.Dial(config.GetUserServiceURL(), opts...)
	if err != nil {
		log.Error("creating gerpc conn", "error", err)
		panic(err)
	}

	user, err := user.NewAdapter(userConn)
	if err != nil {
		log.Error("creating service adapter", "error", err)
		panic(err)
	}
	return user, userConn
}
