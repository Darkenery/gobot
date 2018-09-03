package command

import (
	"github.com/darkenery/gobot/api"
	"github.com/darkenery/gobot/api/model"
	"github.com/darkenery/gobot/bot/util"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"math/rand"
	"strings"
)

type GenerateRandomTextCommand struct {
	botApi    *api.BotApi
	redis     *redis.Client
	wordLimit int
}

var (
	WrongCommandParameterErr = errors.New("wrong amount of parameters")
)

func NewGenerateRandomTextCommand(botApi *api.BotApi, redis *redis.Client, wordLimit int) CommandInterface {
	return &GenerateRandomTextCommand{
		botApi:    botApi,
		redis:     redis,
		wordLimit: wordLimit,
	}
}

func (c *GenerateRandomTextCommand) Execute(incomingMessage *model.Message) error {
	text := util.ExtractTextFromMessage(incomingMessage)
	words := strings.Split(text, " ")
	if len(words) < 1 {
		return WrongCommandParameterErr
	}
	word := words[0]

	var (
		relations        []string
		err              error
		result           string
		currentWordCount int
	)

	for {
		wordKey := fmt.Sprintf(util.SetWordRedisTemplate, word)
		relations, err = c.redis.SMembers(wordKey).Result()
		if err != nil && err != redis.Nil {
			return err
		}

		if err == redis.Nil {
			break
		}

		result += word + " "
		currentWordCount++

		if len(relations) == 0 {
			break
		}

		relationHashes := make(map[int]map[string]string)
		for i, relation := range relations {
			relationHash, err := c.redis.HGetAll(relation).Result()
			if err != nil {
				continue
			}

			relationHashes[i] = relationHash
		}

		if len(relationHashes) == 0 {
			break
		}

		if currentWordCount == c.wordLimit {
			break
		}

		word = relationHashes[rand.Intn(len(relationHashes))]["word"]
	}

	var message string
	if len(result) != 0 {
		message = result
	} else {
		message = "Чёт не получилось, сорян"
	}

	_, err = c.botApi.SendMessage(
		incomingMessage.Chat.Id,
		incomingMessage.MessageId,
		message,
	)

	return err
}
