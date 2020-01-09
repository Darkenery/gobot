package update_processor

import (
	"fmt"
	"github.com/darkenery/gobot/api/model"
	"github.com/darkenery/gobot/bot/util"
	"github.com/go-redis/redis"
	"strings"
)

type FillOneOrderDictionaryProcessor struct {
	redis *redis.ClusterClient
}

func NewFillOneOrderDictionaryProcessor(redis *redis.ClusterClient) UpdateProcessorInterface {
	return &FillOneOrderDictionaryProcessor{
		redis: redis,
	}
}

func (fdp *FillOneOrderDictionaryProcessor) Process(incomingMessage *model.Message) (err error) {
	text := util.ExtractTextFromMessage(incomingMessage)
	text = fdp.preProcessText(text)

	words := strings.Split(text, " ")
	words = fdp.filterWords(words)
	if len(words) < 2 {
		return nil
	}

	for i, word := range words {
		currentWordSetKey := fmt.Sprintf(util.SetWordRedisTemplate, word)
		nextWordSetKey := fmt.Sprintf(util.SetWordRedisTemplate, words[i+1])
		relationToTheNextWordKey := fmt.Sprintf(util.HashRelationRedisTemplate, words[i+1])
		//add word and key for relation to the next word in words array
		err = fdp.redis.SAdd(currentWordSetKey, relationToTheNextWordKey).Err()
		if err != nil {
			return err
		}

		//fill relation to the next word in words array

		//fill word in relation
		kv := make(map[string]interface{})
		kv["set"] = nextWordSetKey
		kv["word"] = words[i+1]
		err = fdp.redis.HMSet(relationToTheNextWordKey, kv).Err()
		if err != nil {
			return
		}

		//fill or update weight
		err = fdp.redis.HIncrBy(relationToTheNextWordKey, "weight", 1).Err()
		if err != nil {
			return
		}

		//if the next word is last - stop filling
		if len(words)-1 <= i+1 {
			break
		}
	}

	return nil
}

func (fdp *FillOneOrderDictionaryProcessor) preProcessText(text string) string {
	text = util.ToLowerCase(text)
	text = util.RemoveWhitespace(text)
	text = util.RemoveNonWordSymbols(text)

	return text
}

func (fdp *FillOneOrderDictionaryProcessor) filterWords(words []string) []string {
	var result []string

	for _, word := range words {
		if word != " " {
			result = append(result, word)
		}
	}

	return result
}
