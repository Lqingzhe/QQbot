package agenttool

import (
	"errors"
	"net/http"
	"sync"

	newerror "github.com/Lqingzhe/go-common-tool/error"
	einotool "github.com/cloudwego/eino/components/tool"
)

type Tool struct {
	ToolHub map[string]bool // true表示为Retune类型，false表示为Timeliness类型
	rwMu    sync.RWMutex
}

func NewTool() (*Tool, error) {
	newTool := &Tool{
		ToolHub: make(map[string]bool, len(ReturnToolHub)+len(TimelinessToolHub)),
	}
	for k := range ReturnToolHub {
		if _, ok := newTool.ToolHub[k]; ok {
			return nil, newerror.MakeError(newerror.CodeInternalError, errors.New("ToolName 重复: "+k), newerror.LevelError)
		}
		newTool.ToolHub[k] = true
	}
	for k := range TimelinessToolHub {
		if _, ok := newTool.ToolHub[k]; ok {
			return nil, newerror.MakeError(newerror.CodeInternalError, errors.New("ToolName 重复: "+k), newerror.LevelError)
		}
		newTool.ToolHub[k] = false
	}
	return newTool, nil
}
func (i *Tool) GetTool() []einotool.BaseTool {
	i.rwMu.RLock()
	list := make([]einotool.BaseTool, 0, len(i.ToolHub))
	for k, v := range i.ToolHub {
		if v {
			list = append(list, ReturnToolHub[k])
		} else {
			list = append(list, TimelinessToolHub[k])
		}
	}
	i.rwMu.RUnlock()
	return list
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
	newHub := make(map[string]bool, len(ReturnToolHub)+len(TimelinessToolHub))
	for k := range ReturnToolHub {
		newHub[k] = true
	}
	for k, _ := range TimelinessToolHub {
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
	if _, ok := ReturnToolHub[toolName]; ok {
		i.rwMu.RUnlock()
		i.rwMu.Lock()
		i.ToolHub[toolName] = true
		i.rwMu.Unlock()
		return nil
	}
	if _, ok := TimelinessToolHub[toolName]; ok {
		i.rwMu.RUnlock()
		i.rwMu.Lock()
		i.ToolHub[toolName] = false
		i.rwMu.Unlock()
		return nil
	}
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
