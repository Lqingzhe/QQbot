package embedding

import (
	"context"
	"net/http"
	"qqbot/internal/model"

	newerror "github.com/Lqingzhe/go-common-tool/error"
	"github.com/cloudwego/eino-ext/libs/acl/openai"
)

type EmbeddingModel struct {
	Model *openai.EmbeddingClient
}

func NewEmbeddingModel(config *model.Config) (*EmbeddingModel, error) {
	embeddingConfig := config.Agent.Embedding
	embedder, err := openai.NewEmbeddingClient(context.Background(), &openai.EmbeddingConfig{
		BaseURL:    embeddingConfig.URL,
		Model:      embeddingConfig.ModelName,
		APIKey:     embeddingConfig.APIKey,
		HTTPClient: http.DefaultClient,
	})
	if err != nil {
		return nil, err
	}
	return &EmbeddingModel{
		Model: embedder,
	}, nil
}

func (i *EmbeddingModel) StringToVector(ctx context.Context, strings []string) (vectors [][]float64, err error) {
	vectors, err = i.Model.EmbedStrings(ctx, strings)
	if err != nil {
		return nil, newerror.MakeError(newerror.CodeThirdPartyError, err, newerror.LevelError)
	}
	return vectors, nil
}
