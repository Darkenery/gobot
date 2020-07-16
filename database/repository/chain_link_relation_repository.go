package repository

import (
	"encoding/json"
	"fmt"
	"github.com/darkenery/gobot/database/model"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
)

type ChainLinkRelationRepository struct {
	db *gorm.DB
}

func NewChainLinkRelationRepository(db *gorm.DB) *ChainLinkRelationRepository {
	return &ChainLinkRelationRepository{
		db: db,
	}
}

func (rep *ChainLinkRelationRepository) FindAllRelationsByWindow(window []uint, order int) (*[]model.ChainLinkRelation, error) {
	jsonWindow, err := json.Marshal(window)
	if err != nil {
		return nil, err
	}

	var relations []model.ChainLinkRelation
	err = rep.db.Find(&relations, model.ChainLinkRelation{Window: postgres.Jsonb{RawMessage: jsonWindow}, Order: order}).Error
	return &relations, err
}

func (rep *ChainLinkRelationRepository) FindAllRelationsWithChainLink(chainLinkId uint, order int) (*[]model.ChainLinkRelation, error) {
	var result []model.ChainLinkRelation
	query := `"order" = $2 AND ("window"->>0)::INTEGER = $1`
	err := rep.db.Where(query, chainLinkId, order).Find(&result).Error
	if err != nil {
		return nil, err
	}
	fmt.Printf("%+v\n", result)
	return &result, err
}

func (rep *ChainLinkRelationRepository) GetRandomRelationByOrder(order int) (result *model.ChainLinkRelation, err error) {
	query := `SELECT * FROM chain_link_relation clr WHERE clr.order = $1 ORDER BY random() DESC LIMIT 1`
	err = rep.db.Raw(query, order).Row().Scan(result)
	return
}

func (rep *ChainLinkRelationRepository) AddRelation(relation *model.ChainLinkRelation) (err error) {
	existingRelation := &model.ChainLinkRelation{}
	if err = rep.db.FirstOrCreate(&existingRelation, relation).Error; err != nil {
		return
	}

	if existingRelation != nil {
		existingRelation.Weight++
		if err = rep.db.Save(existingRelation).Error; err != nil {
			return err
		}
	}

	return nil
}
