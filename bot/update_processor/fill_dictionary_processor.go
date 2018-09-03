package update_processor

import (
	"github.com/go-redis/redis"
	"github.com/darkenery/gobot/api/model"
	"github.com/darkenery/gobot/bot/util"
	"strings"
	"fmt"
)

type FillDictionaryProcessor struct {
	redis *redis.Client
}

func NewFillDictionaryProcessor(redis *redis.Client) UpdateProcessorInterface {
	return &FillDictionaryProcessor{
		redis: redis,
	}
}

func (fdp *FillDictionaryProcessor) Process(incomingMessage *model.Message) (err error) {
	text := util.ExtractTextFromMessage(incomingMessage)
	words := strings.Split(text, " ")

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
		if len(words) - 1 <= i + 1 {
			break
		}
	}

	return nil
}
