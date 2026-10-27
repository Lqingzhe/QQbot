package agenttool

import (
	"context"
	"time"

	newerror "github.com/Lqingzhe/go-common-tool/error"
	newlog "github.com/Lqingzhe/go-common-tool/log"
	einotool "github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

var TimelinessToolHub = make(map[string]einotool.InvokableTool)

func RegisterTimelinessTool(toolName string, tool einotool.InvokableTool) {
	TimelinessToolHub[toolName] = tool
}
func RegisterOfficialTimelinessTool(logger newlog.Logger) {
	RegisterTimelinessTool("get-time-now", NewGetTimeNow(logger))
}

type GetTimeNow struct {
	logger newlog.Logger
}

func NewGetTimeNow(logger newlog.Logger) *GetTimeNow {
	return &GetTimeNow{
		logger: logger,
	}
}
func (*GetTimeNow) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "get-time-now",
		Desc: `获取当前系统时间。当用户询问"现在几点"、"当前时间"、"现在是什么时候"等问题时调用此工具。`,
	}, nil
}

func (g *GetTimeNow) InvokableRun(ctx context.Context, _ string, _ ...einotool.Option) (result string, err error) {
	traceID := ctx.Value("trace_id").(string)
	g.logger.AddTraceID(traceID).AddLevel(newerror.LevelInfo).AddOperation("tool:GetTimeNow").Log()
	return time.Now().String(), nil
}
