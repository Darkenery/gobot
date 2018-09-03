package main

import (
	"github.com/darkenery/gobot/api"
	"github.com/darkenery/gobot/api/model"
	"github.com/darkenery/gobot/bot/command"
	"github.com/darkenery/gobot/bot/command_handler"
	"github.com/darkenery/gobot/bot/update_getter"
	"github.com/darkenery/gobot/bot/update_handler"
	"github.com/darkenery/gobot/config"
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-redis/redis"
	"os"
	"os/signal"
	"syscall"
	"github.com/darkenery/gobot/bot/update_processor"
)

var (
	GitHash    = "dev"
	Build      = ""
	configPath = flag.String("config", "./config.yaml", "Set config path")
)

func main() {
	var logger log.Logger

	logger = log.NewJSONLogger(os.Stdout)
	logger = log.With(logger, "@timestamp", log.DefaultTimestampUTC)
	logger = log.With(logger, "@message", "info")
	logger = log.With(logger, "caller", log.DefaultCaller)

	logger.Log("version", GitHash)
	logger.Log("builddate", Build)
	logger.Log("msg", "hello")
	defer logger.Log("msg", "goodbye")

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
	})
	updatesCh := make(chan []*model.Update)
	messageChan := make(chan *model.Message)

	botApi := api.NewBotApi(
		cfg.Bot.ApiConfig.Url,
		cfg.Bot.ApiConfig.Key,
		cfg.Bot.ApiConfig.Timeout,
		cfg.Bot.ApiConfig.KeepAlive,
		cfg.Bot.ApiConfig.HandshakeTimeout,
	)

	updateGetter := update_getter.NewUpdateGetter(
		botApi,
		updatesCh,
		redisClient,
		cfg.Bot.UpdateGetter.Limit,
		cfg.Bot.UpdateGetter.Timeout,
		cfg.Bot.UpdateGetter.AllowedUpdates,
		logger,
	)

	updateHandler := update_handler.NewUpdateHandler(
		updatesCh,
		messageChan,
		logger,
	)

	fillDictionaryProcessor := update_processor.NewFillDictionaryProcessor(redisClient)
	updateHandler.AddProcessor(fillDictionaryProcessor)

	commandHandler := command_handler.NewCommandHandler(messageChan)
	genRndTxtCmd := command.NewGenerateRandomTextCommand(botApi, redisClient, cfg.Bot.WordLimit)
	commandHandler.AddCommand("/gen", &genRndTxtCmd)

	errCh := make(chan error)
	// Interrupt handler.
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errCh <- fmt.Errorf("%s", <-c)
	}()

	go updateGetter.GetUpdates()
	go updateHandler.HandleUpdates()
	go commandHandler.WaitUpdate()

	logger.Log("exit", <-errCh)
}
