package util

import (
	"github.com/darkenery/gobot/api/model"
	"github.com/darkenery/gobot/api/type"
	"regexp"
	"strings"
	"unicode"
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

func ToLowerCase(text string) string {
	return strings.ToLower(text)
}

func RemoveWhitespace(text string) string {
	var b strings.Builder
	b.Grow(len(text))
	for _, ch := range text {
		if !unicode.IsControl(ch) {
			b.WriteRune(ch)
		} else {
			b.Write([]byte(" "))
		}
	}

	return b.String()
}

func RemoveNonWordSymbols(text string) string {
	reg, _ := regexp.Compile(`[^a-zа-яA-ZА-Я0-9\s]+`)
	return reg.ReplaceAllString(text, "")
}

func UcFirst(text string) string {
	words := strings.Split(text, " ")
	words[0] = strings.Title(words[0])

	return strings.Join(words, " ")
}