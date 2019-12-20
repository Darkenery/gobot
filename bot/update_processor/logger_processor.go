package update_processor

import (
	"github.com/darkenery/gobot/api/model"
	"github.com/go-kit/kit/log"
)

type LoggerProcessor struct {
	logger log.Logger
}

func NewLoggerProcessor (logger log.Logger) UpdateProcessorInterface {
	return &LoggerProcessor{
		logger: logger,
	}
}

func (lg *LoggerProcessor) Process(incomingMessage *model.Message) (err error) {
	return lg.logger.Log("Incoming message", incomingMessage.Text, "From", incomingMessage.From.Id, "Chat", incomingMessage.Chat.Id)
}
