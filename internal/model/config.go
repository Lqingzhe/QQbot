package model

type Config struct {
	Global        GlobalConfig
	Api           ApiConfig
	Handler       HandlerConfig
	Service       ServiceConfig
	ResponserList []Responser
}
type Responser struct {
	QQid      int64
	IsPrivate bool
}
type ServiceConfig struct {
	SelfQQ        int64
	ChannelLength int64
	TimeoutSec    int64
}
type HandlerConfig struct {
	ChannelLength int64
}
type ApiConfig struct {
	ListenerPost int64
}
type GlobalConfig struct {
	DeviceID int
}
