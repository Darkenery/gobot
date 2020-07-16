package main

import (
	"flag"
	"fmt"
	"github.com/darkenery/gobot/api"
	"github.com/darkenery/gobot/api/model"
	"github.com/darkenery/gobot/bot/command"
	"github.com/darkenery/gobot/bot/command_handler"
	"github.com/darkenery/gobot/bot/update_getter"
	"github.com/darkenery/gobot/bot/update_handler"
	"github.com/darkenery/gobot/bot/update_processor"
	"github.com/darkenery/gobot/config"
	"github.com/darkenery/gobot/database"
	dbmodel "github.com/darkenery/gobot/database/model"
	"github.com/darkenery/gobot/database/repository"
	"github.com/darkenery/gobot/service"
	"github.com/go-kit/kit/log"
	"github.com/go-redis/redis"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
	"os/signal"
	"syscall"
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

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisClient := redis.NewClusterClient(&redis.ClusterOptions{Addrs: []string{redisHost + ":" + redisPort}})

	pgsql := os.Getenv("POSTGRES_HOST")
	db := database.NewDB("postgres", pgsql, 0, 0, 0, false)

	db.AutoMigrate(&dbmodel.ChainLink{})
	db.AutoMigrate(&dbmodel.ChainLinkRelation{})

	updatesCh := make(chan []*model.Update)
	messageChan := make(chan *model.Message)

	botKey := os.Getenv("TG_BOT_API_KEY")

	botApi := api.NewBotApi(
		cfg.Bot.ApiConfig.Url,
		botKey,
		cfg.Bot.ApiConfig.Timeout,
		cfg.Bot.ApiConfig.KeepAlive,
		cfg.Bot.ApiConfig.HandshakeTimeout,
	)

	me, err := botApi.GetMe()
	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}

	updateGetter := update_getter.NewUpdateGetter(
		botApi,
		updatesCh,
		redisClient,
		cfg.Bot.UpdateGetter.Limit,
		cfg.Bot.UpdateGetter.Timeout,
		cfg.Bot.UpdateGetter.AllowedUpdates,
		logger,
	)

	chainLinkRep := repository.NewChainLinkRepository(db)
	chainLinkRelationRep := repository.NewChainLinkRelationRepository(db)

	textGeneratorSvc := service.NewTextGeneratorService(chainLinkRep, chainLinkRelationRep)

	fillDictionaryProcessorOrder1 := update_processor.NewFillOneOrderDictionaryProcessor(textGeneratorSvc, 1)
	fillDictionaryProcessorOrder2 := update_processor.NewFillOneOrderDictionaryProcessor(textGeneratorSvc, 2)
	fillDictionaryProcessorOrder3 := update_processor.NewFillOneOrderDictionaryProcessor(textGeneratorSvc, 3)

	updateHandler := update_handler.NewUpdateHandler(
		updatesCh,
		messageChan,
		logger,
	)

	loggerProcessor := update_processor.NewLoggerProcessor(logger)
	updateHandler.AddProcessor(loggerProcessor)
	updateHandler.AddProcessor(fillDictionaryProcessorOrder1)
	updateHandler.AddProcessor(fillDictionaryProcessorOrder2)
	updateHandler.AddProcessor(fillDictionaryProcessorOrder3)

	generateRandomTextCommandOrder1 := command.NewGenerateRandomTextCommand(botApi, textGeneratorSvc, cfg.Bot.GenerateRandomTextCommandConfig.WordLimit, 1)
	generateRandomTextCommandOrder2 := command.NewGenerateRandomTextCommand(botApi, textGeneratorSvc, cfg.Bot.GenerateRandomTextCommandConfig.WordLimit, 2)
	generateRandomTextCommandOrder3 := command.NewGenerateRandomTextCommand(botApi, textGeneratorSvc, cfg.Bot.GenerateRandomTextCommandConfig.WordLimit, 3)

	commandHandler := command_handler.NewCommandHandler(messageChan, me, logger)
	commandHandler.AddCommand("/gen", &generateRandomTextCommandOrder1)
	commandHandler.AddCommand("/gen_2", &generateRandomTextCommandOrder2)
	commandHandler.AddCommand("/gen_3", &generateRandomTextCommandOrder3)

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
