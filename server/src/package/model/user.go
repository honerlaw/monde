package model

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"index;unique;not null"`
	Hash     string `json:"hash" gorm:"not null"`
}
