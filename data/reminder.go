package data

import (
	"gorm.io/gorm"
	"time"
)

type Reminder struct {
	gorm.Model
	Location string
	Text string
	At time.Time
}