package command

import (
	"errors"
	"fmt"
	"github.com/darkenery/gobot/api"
	"github.com/darkenery/gobot/api/model"
	"github.com/darkenery/gobot/bot/util"
	"github.com/go-redis/redis"
	"math/rand"
	"strconv"
	"strings"
)

type GenerateRandomTextCommand struct {
	botApi    *api.BotApi
	redis     *redis.ClusterClient
	wordLimit int
}

var (
	WrongCommandParameterErr = errors.New("wrong amount of parameters")
)

func NewGenerateRandomTextCommand(botApi *api.BotApi, redis *redis.ClusterClient, wordLimit int) CommandInterface {
	return &GenerateRandomTextCommand{
		botApi:    botApi,
		redis:     redis,
		wordLimit: wordLimit,
	}
}

func (c *GenerateRandomTextCommand) Execute(incomingMessage *model.Message) error {
	text := util.ExtractTextFromMessage(incomingMessage)
	inputWords := strings.Split(text, " ")
	if len(inputWords) < 1 {
		return WrongCommandParameterErr
	}
	word := inputWords[0]

	var (
		relations        []string
		err              error
		result           []string
		currentWordCount int
		weightSum        int
	)

	result = append(result, word)

	for {
		wordKey := fmt.Sprintf(util.SetWordRedisTemplate, word)
		relations, err = c.redis.SMembers(wordKey).Result()
		if err != nil && err != redis.Nil {
			return err
		}

		if err == redis.Nil {
			break
		}

		if len(relations) == 0 {
			break
		}

		weightSum = 0
		relationHashes := make(map[int]map[string]string)
		for i, relation := range relations {
			relationHash, err := c.redis.HGetAll(relation).Result()
			if err != nil {
				continue
			}

			weight, _ := strconv.Atoi(relationHash["weight"])
			weightSum += weight
			relationHashes[i] = relationHash
		}

		if len(relationHashes) == 0 {
			break
		}

		word = relationHashes[rand.Intn(len(relationHashes))]["word"]
		result = append(result, word)
		currentWordCount++

		if currentWordCount >= c.wordLimit {
			break
		}
	}

	var message string
	if len(result) > 2 {
		message = strings.Join(result, " ")
	} else {
		message = "Чёт не получилось, сорян"
	}

	message = util.UcFirst(message)

	_, err = c.botApi.SendMessage(
		incomingMessage.Chat.Id,
		incomingMessage.MessageId,
		message,
	)

	return err
}
