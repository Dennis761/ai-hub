package promptports

import (
	"context"

	"ai_hub.com/app/core/domain/promptdomain"
)

type PromptWriteRepository interface {
	Create(ctx context.Context, p *promptdomain.Prompt) (*promptdomain.Prompt, error)

	Update(ctx context.Context, p *promptdomain.Prompt) (*promptdomain.Prompt, error)

	UpdateMany(ctx context.Context, prompts []*promptdomain.Prompt) error

	Delete(ctx context.Context, id promptdomain.PromptID) error
}
