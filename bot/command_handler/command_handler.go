package command_handler

import (
	"github.com/darkenery/gobot/api/model"
	"github.com/darkenery/gobot/bot/command"
	"github.com/go-kit/kit/log"
	"strings"
)

type CommandHandler struct {
	botInfo         *model.User
	commandPool     map[string]command.CommandInterface
	updateHandlerCh chan *model.Message
	logger          log.Logger
}

func NewCommandHandler(updateHandlerCh chan *model.Message, botInfo *model.User, logger log.Logger) *CommandHandler {
	return &CommandHandler{
		updateHandlerCh: updateHandlerCh,
		commandPool:     make(map[string]command.CommandInterface),
		botInfo:         botInfo,
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
		messageCommand := message.Text[entity.Offset : entity.Offset+entity.Length]
		messageCommandParts := strings.Split(messageCommand, "@")

		if len(messageCommandParts) == 1 || messageCommandParts[1] == ch.botInfo.Username {
			if botCommand, ok := ch.commandPool[messageCommandParts[0]]; ok {
				err := botCommand.Execute(message)
				if err != nil {
					ch.logger.Log("err", err)
				}
			}
		}
	}
}
