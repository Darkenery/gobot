package repository

import (
	"github.com/darkenery/gobot/database/model"
	"github.com/jinzhu/gorm"
)

type ChainLinkRepository struct {
	db *gorm.DB
}

func NewChainLinkRepository(db *gorm.DB) *ChainLinkRepository {
	return &ChainLinkRepository{
		db: db,
	}
}

func (rep *ChainLinkRepository) FindByValue(value string) (*model.ChainLink, error) {
	chainLink := model.ChainLink{}
	err := rep.db.Find(&chainLink, model.ChainLink{Value: value}).Error
	return &chainLink, err
}

func (rep *ChainLinkRepository) CreateIfNotExists(chainLink *model.ChainLink) (*model.ChainLink, error) {
	newChainLink := model.ChainLink{}
	err := rep.db.FirstOrCreate(&newChainLink, chainLink).Error
	return &newChainLink, err
}

func (rep *ChainLinkRepository) GetById (id uint) (*model.ChainLink, error) {
	var chainLink model.ChainLink
	err := rep.db.First(&chainLink, id).Error
	return &chainLink, err
}
