package blocks

import (
	"context"

	"github.com/iankencruz/webbuilder/internal/database/repository"
)

// Block is an interface that defines the methods for creating, retrieving,
// updating, and deleting blocks in the application. Each block type (e.g.,
// HeroBlock, RichTextBlock) implements this interface, allowing for a
// consistent way to interact with different block types in the application.
type Block interface {
	Create(ctx context.Context) (int64, error)
	Get(ctx context.Context, id int64) (any, error)
	Update(ctx context.Context) (any, error)
	Delete(ctx context.Context, id int64) error
}

// BlockType is a struct that represents a block type in the application.
// It contains the collection name for the block type and a function to create a new
// instance of the block using the provided repository. This struct is used to
// register block types in the BlockService, allowing for the creation,
// retrieval, updating, and deletion of blocks of that type in the application.
type BlockType struct {
	Collection string
	New        func(q *repository.Queries) Block
}
