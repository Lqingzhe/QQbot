package chatcache

import (
	"context"
	"errors"
	"qqbot/internal/model"
	"strconv"

	"github.com/Lqingzhe/go-common-tool/commondao"
	newerror "github.com/Lqingzhe/go-common-tool/error"
)

func idToKey(sessionID int64, isPrivate bool) string {
	stringID := strconv.FormatInt(sessionID, 10)
	if isPrivate {
		return "1" + stringID
	} else {
		return "2" + stringID
	}
}
func keyToID(key string) (sessionID int64, isPrivate bool, err error) {
	if key == "" {
		return 0, false, newerror.MakeError(newerror.CodeInternalError, errors.New("Redis Key Is ' '"), newerror.LevelError)
	}
	b := []byte(key)
	if b[0] == '1' {
		isPrivate = true
	} else if b[0] == '2' {
		isPrivate = false
	} else {
		return 0, false, newerror.MakeError(newerror.CodeInternalError, errors.New("Redis Key Is Invalid"), newerror.LevelError)
	}
	sessionString := b[1:]
	sessionID, err = strconv.ParseInt(string(sessionString), 10, 64)
	if err != nil {
		return 0, false, newerror.MakeError(newerror.CodeInternalError, err, newerror.LevelError)
	}
	return sessionID, isPrivate, nil
}

func (i *ContextInfo) GetDB(ctx context.Context, dbContext commondao.DBContextInterface[model.DBContext]) (gotInfo *ContextInfo, exist bool, err error) {

}
func (i *ContextInfo) SetDB(ctx context.Context, dbContext commondao.DBContextInterface[model.DBContext]) (err error) {

}
func (i *ContextInfo) DeleteDB(ctx context.Context, dbContext commondao.DBContextInterface[model.DBContext]) (err error) {

}
func (i *ContextInfo) UpdateDB(ctx context.Context, dbContext commondao.DBContextInterface[model.DBContext]) (exist bool, err error) {

}
