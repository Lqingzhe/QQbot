package chatllm

import (
	"context"
	"errors"
	"net/http"
	"qqbot/internal/agent/chatllm/agenttool"
	"qqbot/internal/model"
	"time"

	newerror "github.com/Lqingzhe/go-common-tool/error"
	newlog "github.com/Lqingzhe/go-common-tool/log"
	"github.com/cloudwego/eino-ext/libs/acl/openai"
	"github.com/cloudwego/eino/adk"
	einomodel "github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func newModel(URl string, ApiKey string, ModelName string) (*openai.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	aiModel, err := openai.NewClient(ctx, &openai.Config{
		BaseURL: URl,
		Model:   ModelName,
		APIKey:  ApiKey,
	})
	if err != nil {
		return nil, newerror.MakeError(newerror.CodeDownstreamError, err, newerror.LevelError).SetHttpError(http.StatusBadRequest, "Invalid Model Config")
	}
	return aiModel, nil
}

type LLM struct {
	AgentTool *agenttool.Tool
	Runner    *adk.Runner
	config    *model.Config
}

func NewLLM(config *model.Config, logger newlog.Logger) (*LLM, error) {
	newLLM := &LLM{}
	tool, err := agenttool.NewTool()
	if err != nil {
		return nil, err
	}
	newLLM.AgentTool = tool
	err = newLLM.NewRunner(config, logger)
	if err != nil {
		return nil, err
	}
	newLLM.config = config
	return newLLM, nil
}
func (i *LLM) NewRunner(Config *model.Config, logger newlog.Logger) error {
	mainModel, err := newModel(Config.Agent.ChatLLM.URL, Config.Agent.ChatLLM.APIKey, Config.Agent.ChatLLM.ModelName)
	if err != nil {
		return err
	}
	var spareModel *openai.Client
	if Config.Agent.ChatLLM.SpareURL != "" {
		spareModel, err = newModel(Config.Agent.ChatLLM.SpareURL, Config.Agent.ChatLLM.SpareAPIKey, Config.Agent.ChatLLM.SpareModelName)
		if err != nil {
			return err
		}
	}

	agent, err := adk.NewChatModelAgent(context.Background(), &adk.ChatModelAgentConfig{
		Instruction: Config.Agent.ChatLLM.Promote,
		Model:       mainModel,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: i.AgentTool.GetTool(), // 注册工具
			},
		},
		MaxIterations: Config.Agent.ChatLLM.MaxThinkStep,
		ModelFailoverConfig: &adk.ModelFailoverConfig[*schema.Message]{
			MaxRetries: 1,
			ShouldFailover: func(ctx context.Context, msg *schema.Message, err error) bool {
				return !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded)
			},
			GetFailoverModel: func(ctx context.Context, fc *adk.FailoverContext[*schema.Message]) (einomodel.BaseChatModel, []*schema.Message, error) {
				traceID := ctx.Value("trace_id").(string)
				logger.AddOperation("更换模型").AddTraceID(traceID).AddLevel(newerror.LevelError).Log()
				return spareModel, nil, nil // nil 表示沿用原始输入
			},
		},
	})
	if err != nil {
		return newerror.MakeError(newerror.CodeInvalidParam, err, newerror.LevelError).SetHttpError(http.StatusBadRequest, "Agent Config Error")
	}
	runner := adk.NewRunner(context.Background(), adk.RunnerConfig{Agent: agent})
	i.Runner = runner
	return nil
}
func (i *LLM) GenerateChat(ctx context.Context, history []*schema.Message) (aiMessage string, newHistory []*schema.Message, err error) {
	iter := i.Runner.Run(ctx, history)
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			return "", nil, newerror.MakeError(newerror.CodeThirdPartyError, event.Err, newerror.LevelError)
		}
		message := event.Output.MessageOutput
		if message == nil || message.Message == nil {
			continue
		}
		switch message.Role {
		case schema.Assistant:
			if len(message.Message.ToolCalls) > 0 {
				for _, tc := range message.Message.ToolCalls {
					if i.AgentTool.IsRetuneTool(tc.Function.Name) {
						history = append(history, message.Message)
					}
				}
			} else {
				history = append(history, message.Message)
				aiMessage = message.Message.Content
			}
		case schema.Tool:
			if i.AgentTool.IsRetuneTool(message.ToolName) {
				history = append(history, message.Message)
			}
		}
	}
	return aiMessage, history, nil
}
func (i *LLM) CleanMessage(message []*schema.Message) []*schema.Message {
	if message == nil {
		return nil
	}
	return nil
}
