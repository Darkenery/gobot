package util

import (
	"github.com/darkenery/gobot/api/model"
	"github.com/darkenery/gobot/api/type"
)

const SetWordRedisTemplate = "set.word.%s"
const HashRelationRedisTemplate = "hash.relation.%s"

func ExtractTextFromMessage(message *model.Message) string {
	var lastEntity *model.MessageEntity

	for _, entity := range message.Entities {
		if entity.Type == _type.BotCommandEntityType {
			lastEntity = entity
		}
	}

	if lastEntity != nil {
		if len(message.Text) == lastEntity.Length {
			return ""
		}
		return message.Text[lastEntity.Offset+lastEntity.Length+1:]
	}

	return message.Text
}
