package update_handler

import (
	"github.com/darkenery/gobot/api/model"
	"github.com/darkenery/gobot/api/type"
	"github.com/darkenery/gobot/bot/update_processor"
	"github.com/go-kit/kit/log"
	"time"
)

type UpdateHandler struct {
	updaterCh           chan []*model.Update
	commandHandlerCh    chan *model.Message
	updateProcessorPool []update_processor.UpdateProcessorInterface
	logger              log.Logger
}

func NewUpdateHandler(updaterCh chan []*model.Update, commandHandlerCh chan *model.Message, logger log.Logger) *UpdateHandler {
	var processorsSlice []update_processor.UpdateProcessorInterface

	return &UpdateHandler{
		updaterCh:           updaterCh,
		commandHandlerCh:    commandHandlerCh,
		updateProcessorPool: processorsSlice,
		logger:              logger,
	}
}

func (uh *UpdateHandler) AddProcessor(newProcessor update_processor.UpdateProcessorInterface) {
	uh.updateProcessorPool = append(uh.updateProcessorPool, newProcessor)
}

func (uh *UpdateHandler) HandleUpdates() {
	var (
		isBotCommand bool
		err          error
	)

	now := time.Now().Unix()
	for updates := range uh.updaterCh {
		for _, update := range updates {
			if now - update.Message.Date > 300 {
				continue
			}
			message := update.Message
			for _, entity := range message.Entities {
				if entity.Type == _type.BotCommandEntityType {
					isBotCommand = true
				}
			}

			if isBotCommand {
				uh.commandHandlerCh <- message
			} else {
				for _, updateProcessor := range uh.updateProcessorPool {
					err = updateProcessor.Process(message)
					if err != nil {
						uh.logger.Log("err", err)
					}
				}
			}

			isBotCommand = false
		}
	}
}
