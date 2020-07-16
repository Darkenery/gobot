package model

import "github.com/jinzhu/gorm"

type ChainLink struct {
	gorm.Model
	Value string `gorm:"value"`
}
