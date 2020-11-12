package data

import (
	"gorm.io/gorm"
)

type Karma struct {
	gorm.Model
	Object string // e.g. notepad
	Value int // e.g. 7
	Location string // e.g. #gorm-bot-test
}