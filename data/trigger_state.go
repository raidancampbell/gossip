package data

import (
	"gorm.io/gorm"
)

type TriggerMeta struct {
	gorm.Model
	Disabled bool
	Priority int // maybe a comparator? maybe don't need it at all?
	Name string `gorm:"primaryKey;uniqueIndex"`
}