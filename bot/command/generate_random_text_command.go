package command

import (
	"github.com/darkenery/gobot/api"
	"github.com/darkenery/gobot/api/model"
	"github.com/darkenery/gobot/bot/util"
	"github.com/darkenery/gobot/service"
	"strings"
)

type GenerateRandomTextCommand struct {
	botApi           *api.BotApi
	textGeneratorSvc service.TextGeneratorServiceInterface
	wordLimit        int
	order            int
}

func NewGenerateRandomTextCommand(botApi *api.BotApi, textGeneratorSvc service.TextGeneratorServiceInterface, wordLimit int, order int) CommandInterface {
	return &GenerateRandomTextCommand{
		botApi:           botApi,
		textGeneratorSvc: textGeneratorSvc,
		wordLimit:        wordLimit,
		order:            order,
	}
}

func (c *GenerateRandomTextCommand) Execute(incomingMessage *model.Message) error {
	text := util.ExtractTextFromMessage(incomingMessage)
	words := strings.Split(text, " ")

	var message string
	var err error

	if len(words) > 0 {
		message, err = c.textGeneratorSvc.Generate(words[0], c.wordLimit, c.order)
	} else {
		message, err = c.textGeneratorSvc.Generate("", c.wordLimit, c.order)
	}

	if err != nil {
		message = "Ошибонька"
	}

	_, _ = c.botApi.SendMessage(incomingMessage.Chat.Id, incomingMessage.MessageId, message)
	return err
}
