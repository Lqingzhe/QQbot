package id

import (
	newlog "system/pkg/log"

	"github.com/bwmarrin/snowflake"
)

func InitSnowNode(equipID int, logger newlog.Logger) *snowflake.Node {
	newNode, err := snowflake.NewNode(int64(equipID))
	if err != nil {
		logger.LogInitFatal(err, "Init SnowNode Failed")
	}
	return newNode
}
func MakeID(Node *snowflake.Node) int64 {
	return Node.Generate().Int64()
}
