package update_processor

import (
	"github.com/darkenery/gobot/api/model"
	"github.com/darkenery/gobot/bot/util"
	"github.com/darkenery/gobot/service"
)

type FillOneOrderDictionaryProcessor struct {
	textGeneratorService service.TextGeneratorServiceInterface
	order                int
}

func NewFillOneOrderDictionaryProcessor(textGeneratorService service.TextGeneratorServiceInterface, order int) UpdateProcessorInterface {
	return &FillOneOrderDictionaryProcessor{
		textGeneratorService: textGeneratorService,
		order:                order,
	}
}

func (fdp *FillOneOrderDictionaryProcessor) Process(incomingMessage *model.Message) (err error) {
	err = fdp.textGeneratorService.Learn(util.ExtractTextFromMessage(incomingMessage), fdp.order)
	return
}