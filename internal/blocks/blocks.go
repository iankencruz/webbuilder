package blocks

import (
	"context"
	"fmt"
	"log/slog"

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

// --- Service ---
type Repository interface {
	AddBlockToPage(ctx context.Context, arg repository.AddBlockToPageParams) (repository.PagesBlock, error)
	GetPageBlocks(ctx context.Context, pageID int64) ([]repository.PagesBlock, error)
	UpdatePageBlock(ctx context.Context, arg repository.UpdatePageBlockParams) (repository.PagesBlock, error)
	DeletePageBlock(ctx context.Context, id int64) error
	ReorderPageBlocks(ctx context.Context, arg repository.ReorderPageBlocksParams) error
}

type BlockService struct {
	queries  *repository.Queries
	repo     Repository
	registry map[string]BlockType
}

func NewService(logger *slog.Logger, q *repository.Queries, types []BlockType) *BlockService {
	s := &BlockService{
		queries:  q,
		repo:     q,
		registry: make(map[string]BlockType),
	}

	for _, t := range types {
		s.registry[t.Collection] = t
	}

	return s
}

func (s *BlockService) Resolve(collection string) (Block, error) {
	bt, ok := s.registry[collection]
	if !ok {
		return nil, fmt.Errorf("block collection not found: %s", collection)
	}
	return bt.New(s.queries), nil
}

// ---- Junction methods ----

// AddBlockToPage adds a block to a page in the database using the provided parameters and returns the newly created PagesBlock. If an error occurs during
// the operation, it returns an error.
func (s *BlockService) AddBlockToPage(ctx context.Context, arg repository.AddBlockToPageParams) (repository.PagesBlock, error) {
	return s.repo.AddBlockToPage(ctx, arg)
}

// GetPageBlocks retrieves all blocks associated with a specific page from the
// database using the provided page ID. It returns a slice of PagesBlock if
// successful, or an error if there is an issue with the database query or if
// the page does not exist.
func (s *BlockService) GetPageBlocks(ctx context.Context, pageID int64) ([]repository.PagesBlock, error) {
	return s.repo.GetPageBlocks(ctx, pageID)
}

// UpdatePageBlock updates the details of a block associated with a page in the
// database using the provided update parameters. It returns the updated
// PagesBlock if successful, or an error if there is an issue with the database
// query or if the block does not exist.
func (s *BlockService) UpdatePageBlock(ctx context.Context, arg repository.UpdatePageBlockParams) (repository.PagesBlock, error) {
	return s.repo.UpdatePageBlock(ctx, arg)
}

// DeletePageBlock removes a block from a page in the database using the
// provided block ID. It returns an error if there is an issue with the database
// query or if the block does not exist.
func (s *BlockService) DeletePageBlock(ctx context.Context, id int64) error {
	return s.repo.DeletePageBlock(ctx, id)
}

// ReorderPageBlocks updates the order of blocks on a page in the database using the provided reorder parameters.
// It returns an error if there is an issue with the database query or if any of the blocks do not exist.
// This method allows for the dynamic reordering of blocks on a page, enabling users to customize the layout
// of their pages by changing the order of blocks as needed.
func (s *BlockService) ReorderPageBlocks(ctx context.Context, arg repository.ReorderPageBlocksParams) error {
	return s.repo.ReorderPageBlocks(ctx, arg)
}
