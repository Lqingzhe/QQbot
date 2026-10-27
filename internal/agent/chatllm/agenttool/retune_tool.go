package agenttool

import (
	"errors"

	newerror "github.com/Lqingzhe/go-common-tool/error"
	einotool "github.com/cloudwego/eino/components/tool"
)

func (i *Tool)RegisterReturnTool(toolName string, tool einotool.InvokableTool) error{
	if _,ok:=i.ReturnToolHub[toolName];!ok {
		i.ReturnToolHub[toolName] = tool
		return nil
	}else{
		return newerror.MakeError(newerror.CodeResourceDuplicate,errors.New("ReturnTool名字重复: "+toolName),newerror.LevelFatal)
	}
}
func (i *Tool)registerOfficialReturnTool() {

}
