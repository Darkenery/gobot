package command_handler

import (
	"github.com/darkenery/gobot/api/model"
	"github.com/darkenery/gobot/bot/command"
	"github.com/go-kit/kit/log"
)

type CommandHandler struct {
	commandPool     map[string]command.CommandInterface
	updateHandlerCh chan *model.Message
	logger          log.Logger
}

func NewCommandHandler(updateHandlerCh chan *model.Message, logger log.Logger) *CommandHandler {
	return &CommandHandler{
		updateHandlerCh: updateHandlerCh,
		commandPool:     make(map[string]command.CommandInterface),
		logger:          logger,
	}
}

func (ch *CommandHandler) AddCommand(commandName string, newCommand *command.CommandInterface) {
	ch.commandPool[commandName] = *newCommand
}

func (ch *CommandHandler) WaitUpdate() {
	for message := range ch.updateHandlerCh {
		go ch.processUpdate(message)
	}
}

func (ch *CommandHandler) processUpdate(message *model.Message) {
	for _, entity := range message.Entities {
		if botCommand, ok := ch.commandPool[message.Text[entity.Offset:entity.Offset+entity.Length]]; ok {
			err := botCommand.Execute(message)
			if err != nil {
				ch.logger.Log("err", err)
			}
		}
	}
}
