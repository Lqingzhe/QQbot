package agenttool

import (
	"errors"
	"net/http"
	"sync"

	newerror "github.com/Lqingzhe/go-common-tool/error"
	newlog "github.com/Lqingzhe/go-common-tool/log"
	einotool "github.com/cloudwego/eino/components/tool"
)

type Tool struct {
	ToolHub           map[string]bool // true表示为Retune类型，false表示为Timeliness类型
	TimelinessToolHub map[string]einotool.InvokableTool
	ReturnToolHub     map[string]einotool.InvokableTool
	rwMu              sync.RWMutex
}

func NewTool(logger newlog.Logger) *Tool {
	newTool := &Tool{
		TimelinessToolHub: make(map[string]einotool.InvokableTool),
		ReturnToolHub:     make(map[string]einotool.InvokableTool),
	}
	newTool.registerOfficialTimelinessTool(logger)
	newTool.registerOfficialReturnTool()
	return newTool
}
func (i *Tool) GetTool() ([]einotool.BaseTool, error) {
	i.rwMu.RLock()
	if i.ToolHub == nil {
		i.rwMu.RUnlock()
		i.rwMu.Lock()
		defer i.rwMu.Unlock()
		i.ToolHub = make(map[string]bool, len(i.ReturnToolHub)+len(i.TimelinessToolHub))
		for k := range i.ReturnToolHub {
			if _, ok := i.ToolHub[k]; ok {
				return nil, newerror.MakeError(newerror.CodeResourceDuplicate, errors.New("ToolName 重复: "+k), newerror.LevelFatal)
			}
			i.ToolHub[k] = true
		}
		for k := range i.TimelinessToolHub {
			if _, ok := i.ToolHub[k]; ok {
				return nil, newerror.MakeError(newerror.CodeResourceDuplicate, errors.New("ToolName 重复: "+k), newerror.LevelFatal)
			}
			i.ToolHub[k] = false
		}
	}
	list := make([]einotool.BaseTool, 0, len(i.ToolHub))
	for k, v := range i.ToolHub {
		if v {
			list = append(list, i.ReturnToolHub[k])
		} else {
			list = append(list, i.TimelinessToolHub[k])
		}
	}
	i.rwMu.RUnlock()
	return list, nil
}
func (i *Tool) GetUsedToolList() map[string]bool {
	i.rwMu.RLock()
	newHub := make(map[string]bool, len(i.ToolHub))
	for k, v := range i.ToolHub {
		newHub[k] = v
	}
	i.rwMu.RUnlock()
	return newHub
}
func (i *Tool) GetAllToolList() map[string]bool {
	i.rwMu.RLock()
	newHub := make(map[string]bool, len(i.ReturnToolHub)+len(i.TimelinessToolHub))
	for k := range i.ReturnToolHub {
		newHub[k] = true
	}
	for k, _ := range i.TimelinessToolHub {
		newHub[k] = false
	}
	i.rwMu.RUnlock()
	return newHub
}
func (i *Tool) AddTool(toolName string) error {
	i.rwMu.RLock()
	if _, ok := i.ToolHub[toolName]; ok {
		i.rwMu.RUnlock()
		return nil
	}
	if _, ok := i.ReturnToolHub[toolName]; ok {
		i.rwMu.RUnlock()
		i.rwMu.Lock()
		i.ToolHub[toolName] = true
		i.rwMu.Unlock()
		return nil
	}
	if _, ok := i.TimelinessToolHub[toolName]; ok {
		i.rwMu.RUnlock()
		i.rwMu.Lock()
		i.ToolHub[toolName] = false
		i.rwMu.Unlock()
		return nil
	}
	i.rwMu.RUnlock()
	return newerror.MakeError(newerror.CodeResourceNotFound, errors.New("Try To Add Unexist Tool"), newerror.LevelInfo).SetHttpError(http.StatusNotFound, "ToolName Is Not Exist: "+toolName)
}
func (i *Tool) RemoveTool(toolName string) error {
	i.rwMu.RLock()
	if _, ok := i.ToolHub[toolName]; !ok {
		i.rwMu.RUnlock()
		return nil
	}
	i.rwMu.Lock()
	delete(i.ToolHub, toolName)
	i.rwMu.Unlock()
	return nil
}
func (i *Tool) IsRetuneTool(toolName string) bool {
	i.rwMu.RLock()
	b := i.ToolHub[toolName]
	i.rwMu.RUnlock()
	return b
}
