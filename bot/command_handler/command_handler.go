package command_handler

import (
	"github.com/darkenery/gobot/api/model"
	"github.com/darkenery/gobot/bot/command"
)

type CommandHandler struct {
	commandPool     map[string]command.CommandInterface
	updateHandlerCh chan *model.Message
}

func NewCommandHandler(updateHandlerCh chan *model.Message) *CommandHandler {
	return &CommandHandler{
		updateHandlerCh: updateHandlerCh,
		commandPool: make(map[string]command.CommandInterface),
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
			botCommand.Execute(message)
		}
	}
}
