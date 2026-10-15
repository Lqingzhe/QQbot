package handlemessage

import "sync"

type OfficialPlug struct {

	//operator提供的功能
	IsPrivate func(int64) bool
	IsGroup   func(int64) bool
	GetSelfQQ func() int64
}

var officialPlug *OfficialPlug

func NewOfficialPlug(IsPrivate func(int64) bool, IsGroup func(int64) bool, GetSelfQQ func() int64) *OfficialPlug {
	sync.OnceFunc(
		func() {
			officialPlug = &OfficialPlug{
				IsPrivate: IsPrivate,
				IsGroup:   IsGroup,
				GetSelfQQ: GetSelfQQ,
			}
		})
	return officialPlug
}
