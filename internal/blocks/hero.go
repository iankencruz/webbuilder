package blocks

import (
	"context"

	"github.com/iankencruz/webbuilder/internal/database/repository"
)

// Hero is a BlockType that represents a hero block in the application. It
// defines the collection name for the hero block and provides a function to
// create a new instance of HeroBlock using the provided repository. This
// BlockType can be registered in the BlockService to allow for the creation,
// retrieval, updating, and deletion of hero blocks in the application.
var Hero = BlockType{
	Collection: "hero",
	New: func(q *repository.Queries) Block {
		return NewHeroBlock(q)
	},
}

// HeroBlock represents a hero block in the application. It contains a reference
// to the repository for database interactions, as well as parameters for
// creating and updating hero blocks. This struct provides methods for creating,
// retrieving, updating, and deleting hero blocks in the database.
type HeroBlock struct {
	queries      *repository.Queries
	Params       repository.CreateHeroBlockParams
	UpdateParams repository.UpdateHeroBlockParams
}

// NewHeroBlock creates a new instance of HeroBlock with the provided
// repository. This function initializes the HeroBlock struct, allowing it to
// interact with the database through the repository for creating, retrieving,
// updating, and deleting hero blocks. The returned HeroBlock instance can then
// be used to perform operations related to hero blocks in the application.
func NewHeroBlock(q *repository.Queries) *HeroBlock {
	return &HeroBlock{
		queries: q,
	}
}

// Create creates a new hero block in the database using the provided parameters
// and returns the ID of the newly created block. If an error occurs during
// creation, it returns 0 and the error.
func (b *HeroBlock) Create(ctx context.Context) (int64, error) {
	result, err := b.queries.CreateHeroBlock(ctx, b.Params)
	if err != nil {
		return 0, err
	}
	return result.ID, nil
}

// Get retrieves a hero block from the database using the provided ID. It
// returns the hero block data if found, or an error if the block does not exist
// or if there is an issue with the database query.
func (b *HeroBlock) Get(ctx context.Context, id int64) (any, error) {
	return b.queries.GetHeroBlock(ctx, id)
}

// Update modifies an existing hero block in the database using the provided
// update parameters. It returns the updated hero block data if the update is
// successful, or an error if the block does not exist or if there is an issue
// with the database query.
func (b *HeroBlock) Update(ctx context.Context) (any, error) {
	return b.queries.UpdateHeroBlock(ctx, b.UpdateParams)
}

// Delete removes a hero block from the database using the provided ID. It
// returns nil if the deletion is successful, or an error if the block does not
// exist or if there is an issue with the database query.
func (b *HeroBlock) Delete(ctx context.Context, id int64) error {
	return b.queries.DeleteHeroBlock(ctx, id)
}
