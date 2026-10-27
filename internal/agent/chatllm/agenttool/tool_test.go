package agenttool

import (
	"testing"

	newlog "github.com/Lqingzhe/go-common-tool/log"
)

func Test_NewTool(t *testing.T) {
	tool := NewTool(newlog.InitLog("task-NewTool",1))

	t.Log(tool)
}
