package model

import "gorm.io/gorm"

type DBContext struct {
	PgVector *PgVectorContext
}

func (d *DBContext) GetClient() *DBContext {
	return d
}

type PgVectorContext struct {
	Client *gorm.DB
	Limit  int64
}
