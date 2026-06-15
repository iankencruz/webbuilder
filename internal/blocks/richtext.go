package blocks

import (
	"context"
	"encoding/json"

	"github.com/iankencruz/webbuilder/internal/database/repository"
)

var RichText = BlockType{
	Collection: "richtext",
	New: func(q *repository.Queries) Block {
		return NewRichTextBlock(q)
	},
}

type RichTextBlock struct {
	queries *repository.Queries
	Params  repository.CreateRichTextBlockParams
	id      int64
}

func NewRichTextBlock(q *repository.Queries) *RichTextBlock {
	return &RichTextBlock{
		queries: q,
	}
}

func (b *RichTextBlock) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &b.Params)
}

func (b *RichTextBlock) SetID(id int64) {
	b.id = id
}

func (b *RichTextBlock) Create(ctx context.Context) (int64, error) {
	result, err := b.queries.CreateRichTextBlock(ctx, repository.CreateRichTextBlockParams{
		Content: b.Params.Content,
	})
	if err != nil {
		return 0, err
	}
	return result.ID, nil
}

func (b *RichTextBlock) Get(ctx context.Context, id int64) (any, error) {
	return b.queries.GetRichTextBlock(ctx, id)
}

func (b *RichTextBlock) Update(ctx context.Context) (any, error) {
	return b.queries.UpdateRichTextBlock(ctx, repository.UpdateRichTextBlockParams{
		ID:      b.id,
		Content: b.Params.Content,
		Format:  b.Params.Format,
	})
}

func (b *RichTextBlock) Delete(ctx context.Context, id int64) error {
	return b.queries.DeleteRichTextBlock(ctx, id)
}
