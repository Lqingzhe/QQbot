package agenttool

import (
	einotool "github.com/cloudwego/eino/components/tool"
)

var ReturnToolHub = make(map[string]einotool.InvokableTool)

func RegisterReturnTool(toolName string, tool einotool.InvokableTool) {
	ReturnToolHub[toolName] = tool
}
func RegisterOfficialReturnTool() {
}
