package command

import "github.com/darkenery/gobot/api/model"

type CommandInterface interface {
	Execute(incomingMessage *model.Message) error
}
