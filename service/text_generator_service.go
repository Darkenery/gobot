package service

import (
	"encoding/json"
	"errors"
	"github.com/darkenery/gobot/bot/util"
	"github.com/darkenery/gobot/database/model"
	"github.com/darkenery/gobot/database/repository"
	"github.com/darkenery/gobot/service/consts"
	"github.com/jinzhu/gorm/dialects/postgres"
	"math/rand"
	"strings"
)

type Relation struct {
	From []string
	To   string
}

type TextGeneratorServiceInterface interface {
	Learn(text string, order int) error
	Generate(startingWord string, limit, order int) (string, error)
}

type textGeneratorService struct {
	chainLinkRep         *repository.ChainLinkRepository
	chainLinkRelationRep *repository.ChainLinkRelationRepository
}

func NewTextGeneratorService(chainLinkRep *repository.ChainLinkRepository, chainLinkRelationRep *repository.ChainLinkRelationRepository) TextGeneratorServiceInterface {
	return &textGeneratorService{
		chainLinkRep:         chainLinkRep,
		chainLinkRelationRep: chainLinkRelationRep,
	}
}

func (svc *textGeneratorService) Generate(startingWord string, limit, order int) (string, error) {
	if startingWord == "" {
		startingWord = consts.StartToken
	}

	FirstChainLink, err := svc.chainLinkRep.FindByValue(startingWord)
	if err != nil {
		return "", err
	}

	chainLinkRelations, err := svc.chainLinkRelationRep.FindAllRelationsWithChainLink(FirstChainLink.Model.ID, order)
	if chainLinkRelations == nil || len(*chainLinkRelations) == 0 {
		return "", errors.New("no result")
	}

	a := *chainLinkRelations
	startingRelation := a[rand.Intn(len(a))]
	text, err := svc.generate(&startingRelation, limit, order)
	if err != nil {
		return "", err
	}

	return svc.postProcessText(text), nil
}

func (svc *textGeneratorService) generate(startingRelation *model.ChainLinkRelation, limit int, order int) (result string, err error) {
	chainLinksIds := []uint{}
	windowInt := []uint{}
	err = json.Unmarshal(startingRelation.Window.RawMessage, &windowInt)
	if err != nil {
		return
	}

	chainLinksIds = append(chainLinksIds, windowInt...)
	for len(chainLinksIds) < limit {
		relations, err := svc.chainLinkRelationRep.FindAllRelationsByWindow(windowInt, order)
		if err != nil {
			return "", err
		}

		if len(*relations) == 0 {
			break
		}

		relation := svc.getRelation(relations)
		chainLinksIds = append(chainLinksIds, relation.Next)
		windowInt = chainLinksIds[len(chainLinksIds)-order:]
	}

	resultWords := []string{}
	for _, chainLinkId := range chainLinksIds {
		chainLink, err := svc.chainLinkRep.GetById(chainLinkId)
		if err != nil {
			return "", err
		}

		resultWords = append(resultWords, chainLink.Value)
	}

	return strings.Join(resultWords, " "), err
}

func (svc *textGeneratorService) getRelation(relations *[]model.ChainLinkRelation) *model.ChainLinkRelation {
	mapp := map[int]model.ChainLinkRelation{}

	weightSum := 0
	a := *relations
	for _, relation := range a {
		weightSum += relation.Weight
		mapp[weightSum] = relation
	}

	random := rand.Intn(weightSum)
	for k, relation := range mapp {
		if random > k-relation.Weight && random <= k {
			return &relation
		}
	}

	return &a[0]
}

func (svc *textGeneratorService) Learn(text string, order int) error {
	text = svc.preProcessText(text)

	words := strings.Split(text, " ")
	words = svc.filterWords(words)
	words = svc.fillWithTokens(words)

	if len(words) < order-1 {
		return nil
	}

	var chainLinks []model.ChainLink
	for _, word := range words {
		chainLink, err := svc.chainLinkRep.CreateIfNotExists(&model.ChainLink{Value: word})
		if err != nil {
			return err
		}
		chainLinks = append(chainLinks, *chainLink)
	}

	for i, _ := range chainLinks {
		var windowIds []uint
		var nextId uint

		j := i
		for len(windowIds) < order {
			windowIds = append(windowIds, chainLinks[j].Model.ID)
			if j >= len(chainLinks) {
				break
			}
			j++
		}

		if j >= len(chainLinks) {
			break
		}
		nextId = chainLinks[j].Model.ID

		jsonWindow, err := json.Marshal(windowIds)
		if err != nil {
			return err
		}

		err = svc.chainLinkRelationRep.AddRelation(&model.ChainLinkRelation{Window: postgres.Jsonb{RawMessage: jsonWindow}, Next: nextId, Order: order})
		if err != nil {
			return err
		}
	}

	return nil
}

func (svc *textGeneratorService) preProcessText(text string) string {
	text = util.ToLowerCase(text)
	text = util.RemoveWhitespace(text)
	text = util.RemoveNonWordSymbols(text)

	return text
}

func (svc *textGeneratorService) postProcessText(text string) string {
	text = util.ClearTokens(text)
	text = util.UcFirst(text)

	return text
}

func (svc *textGeneratorService) filterWords(words []string) []string {
	var result []string

	for _, word := range words {
		if word != " " {
			result = append(result, word)
		}
	}

	return result
}

func (svc *textGeneratorService) fillWithTokens(words []string) []string {
	return append([]string{consts.StartToken}, append(words, consts.FinishToken)...)
}
