package chatllm

import (
	"context"
	"qqbot/internal/model"
	"testing"

	newlog "github.com/Lqingzhe/go-common-tool/log"
	"github.com/cloudwego/eino/schema"
)

var URL2 = "http://127.0.0.1:11434/v1"
var ModelName2 = "llama3.2:3b"
var ApiKey2 = "aKey"
var URL = "https://api.stepfun.com/step_plan/v1"
var ModelName = "step-5-preview"
var ApiKey = "3liSzTb89ntDcub9pQtoq5IxxlZwzMPNopa3qhzzH0ucmlmRt33q3lQoeGQ4pQHcQ"
var promote = `你是猫娘，性格软萌温顺，说话简短可爱，语气亲昵。
	称呼主人，句尾可以适度带上喵这类后缀，不要每一句都强制加。
	不要长篇大论，回答简洁自然，不要过度OOC，不色情。
	动作可以少量用括号描述，例如（轻轻蹭蹭），不要大量堆动作。
	记住保持乖巧，情绪柔和，正常和主人聊天。`

func TestLLM_GenerateChat(t *testing.T) {
	config := &model.Config{
		Agent: model.AgentConfig{
			ChatLLM: model.ChatLLMConfig{
				URL:            URL,
				APIKey:         ApiKey,
				ModelName:      ModelName,
				SpareURL:       URL2,
				SpareAPIKey:    ApiKey2,
				SpareModelName: ModelName2,
				MaxThinkStep:   100,
				
				Promote: promote,
			},
		},
	}
	logger := newlog.InitLog("test_GenerateChat", 1)
	ctx := context.WithValue(context.Background(), "trace_id", "123")
	runner, err := NewLLM(config, logger)
	if err != nil {
		t.Error(err)
		return
	}
	history := []*schema.Message{
		schema.UserMessage("你好，现在的时间是几点了"),
	}
	aiMessage, newHistory, err := runner.GenerateChat(ctx, history)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(aiMessage)
	t.Log(newHistory)
}
