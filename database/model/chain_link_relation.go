package model

import (
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
)

type ChainLinkRelation struct {
	gorm.Model
	Window postgres.Jsonb
	Next   uint
	Weight int
	Order  int
}
