package model

import (
	"time"

	"gorm.io/gorm"
)

type MemeryModel struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	SessionID int64
	UserID    int64
	Content   string
	Vector    []float32
	//缺少置信度部分
}

func (_ *MemeryModel) Data() {

}

type MessageData struct {
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	SessionID  int64
	MessageID  int64
	IsPrivate  bool
	SendUserID int64
	Content    string
}

func (_ *MessageData) Data() {}

type ChatCacheModel struct {
	SessionID int64
	IsPrivate bool
	Data      []byte
}
