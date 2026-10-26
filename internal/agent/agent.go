package agent

import (
	"qqbot/internal/agent/chatllm"
	"qqbot/internal/agent/embedding"
)

type Agent struct {
	Embedding *embedding.EmbeddingModel
	ChatLLM  *chatllm.LLM

}
func NewAgent(chatLLM *chatllm.LLM,embedding *embedding.EmbeddingModel) *Agent {
	return &Agent{
		Embedding: embedding,
		ChatLLM:  chatLLM,
	}
}