package blocks

import (
	"context"

	"github.com/iankencruz/webbuilder/internal/database/repository"
)

var RichText = BlockType{
	Collection: "rich_text",
	New: func(q *repository.Queries) Block {
		return NewRichTextBlock(q)
	},
}

type RichTextBlock struct {
	queries      *repository.Queries
	Params       repository.CreateRichTextBlockParams
	UpdateParams repository.UpdateRichTextBlockParams
}

func NewRichTextBlock(q *repository.Queries) *RichTextBlock {
	return &RichTextBlock{
		queries: q,
	}
}

func (b *RichTextBlock) Create(ctx context.Context) (int64, error) {
	result, err := b.queries.CreateRichTextBlock(ctx, b.Params)
	if err != nil {
		return 0, err
	}
	return result.ID, nil
}

func (b *RichTextBlock) Get(ctx context.Context, id int64) (any, error) {
	return b.queries.GetRichTextBlock(ctx, id)
}

func (b *RichTextBlock) Update(ctx context.Context) (any, error) {
	return b.queries.UpdateRichTextBlock(ctx, b.UpdateParams)
}

func (b *RichTextBlock) Delete(ctx context.Context, id int64) error {
	return b.queries.DeleteHeroBlock(ctx, id)
}
