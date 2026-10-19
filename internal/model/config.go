package model

type Config struct {
	Global        GlobalConfig
	Api           ApiConfig
	Handler       HandlerConfig
	Service       ServiceConfig
	Dao           DBConfig
	Agent         AgentConfig
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
	Websocket     WebsocketConfig
}
type WebsocketConfig struct {
	Host  string
	Port  string
	Token string
}
type ApiConfig struct {
	ListenerPost int64
}
type GlobalConfig struct {
	DeviceID int
}
type DBConfig struct {
	PgConfig PgVectorConfig
}
type PgVectorConfig struct {
	Username    string
	Password    string
	Host        string
	Port        string
	MessageList int64
}

type AgentConfig struct {
	ChatLLM   ChatLLMConfig
	Embedding EmbeddingConfig
}
type ChatLLMConfig struct {
	URL       string
	APIKey    string
	ModelName string
}
type EmbeddingConfig struct {
	URL       string
	APIKey    string
	ModelName string
}
