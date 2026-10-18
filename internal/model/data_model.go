package model

import (
	"time"

	"gorm.io/gorm"
)

type MemeryModel struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	ID      int64 `gorm:"primarykey;autoIncrement"`
	Content string
	Vector  []float32
	缺少置信度部分
}
