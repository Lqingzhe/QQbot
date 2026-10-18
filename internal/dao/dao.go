package dao

import (
	"demo/internal/model"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDBContext(config *model.Config) (*model.DBContext, error) {

	pgvector,err:=initPgVector(config)
	if err != nil {
		return nil, err
	}
	return &model.DBContext{
		PgVector: pgvector,

	}, nil
}

func initPgVector(config *model.Config) (*model.PgVectorContext,error) {
	pgConfig := config.Dao.PgConfig
	newContext := &model.PgVectorContext{
		Limit: pgConfig.MessageList,
	}
	url := "postgres://" + pgConfig.Username + ":" + pgConfig.Password + "@" + pgConfig.Host + ":" + pgConfig.Port + "/" + "vec_db?sslmode=disable"
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Warn),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		return nil,fmt.Errorf("postgreSql connect error: %w", err)
	}
	if err=db.Exec("CREATE EXTENSION IF NOT EXISTS vector").Error;err!=nil{
		return nil,fmt.Errorf("postgreSQL can not support vector: %w", err)
	}
	newContext.Client=db
	return newContext,nil
}
