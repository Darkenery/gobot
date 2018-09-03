package update_processor

import "github.com/darkenery/gobot/api/model"

type UpdateProcessorInterface interface {
	Process(incomingMessage *model.Message) error
}
