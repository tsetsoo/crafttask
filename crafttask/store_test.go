package crafttask

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryStoreInsert(t *testing.T) {
	store := NewInMemoryStore()

	block1 := blockRequest{Content: "Block 1"}
	block2 := blockRequest{Content: "Block 2"}

	payload := insertPayload{
		InsertOperations: []insertOperation{
			{ParentBlockId: root, Index: 0, Block: block1},
			{ParentBlockId: root, Index: 1, Block: block2},
		},
	}

	err := store.insert(payload)
	assert.NoError(t, err)

	// Retrieve blocks and check if they exist in the store
	block1Retrieved, _, _, err := store.findBlockById(1)
	assert.NoError(t, err)
	assert.Equal(t, block1.Content, block1Retrieved.content)

	block2Retrieved, _, _, err := store.findBlockById(2)
	assert.NoError(t, err)
	assert.Equal(t, block2.Content, block2Retrieved.content)
}

func TestInMemoryStoreDelete(t *testing.T) {
	store := NewInMemoryStore()

	block1 := blockRequest{Content: "Block 1"}
	block2 := blockRequest{Content: "Block 2"}

	payload := insertPayload{
		InsertOperations: []insertOperation{
			{ParentBlockId: root, Index: 0, Block: block1},
			{ParentBlockId: root, Index: 1, Block: block2},
		},
	}

	err := store.insert(payload)
	assert.NoError(t, err)

	// Delete block with id 1
	store.delete([]uint64{1})

	// Verify that block1 is deleted
	_, _, _, err = store.findBlockById(1)
	assert.Error(t, err)
	assert.Equal(t, errBlockDoesNotExist, err)
}

func TestInMemoryStoreDuplicate(t *testing.T) {
	store := NewInMemoryStore()

	block1 := blockRequest{Content: "Block 1"}
	payload := insertPayload{
		InsertOperations: []insertOperation{
			{ParentBlockId: root, Index: 0, Block: block1},
		},
	}

	err := store.insert(payload)
	assert.NoError(t, err)

	childBlock := blockRequest{Content: "Child Block 1"}
	payload = insertPayload{
		InsertOperations: []insertOperation{
			{ParentBlockId: 1, Index: 0, Block: childBlock},
		},
	}

	err = store.insert(payload)
	assert.NoError(t, err)

	// Duplicate block with id 1
	duplicatedBlock, err := store.duplicate(1)
	assert.NoError(t, err)

	// Verify that the block is duplicated
	assert.NoError(t, err)
	assert.Equal(t, block1.Content, duplicatedBlock.content)
	require.Len(t, duplicatedBlock.subblocks.values, 1)
	assert.Equal(t, duplicatedBlock.subblocks.OrderedValues()[0].content, childBlock.Content)
}

// func TestInMemoryStoreMove(t *testing.T) {
// 	store := NewInMemoryStore()

// 	block1 := blockRequest{ Content: "Block 1"}
// 	block2 := blockRequest{ Content: "Block 2"}

// 	payload := insertPayload{
// 		InsertOperations: []insertOperation{
// 			{ParentBlockId: root, Index: 0, Block: block1},
// 			{ParentBlockId: root, Index: 1, Block: block2},
// 		},
// 	}

// 	err := store.insert(payload)
// 	assert.NoError(t, err)

// 	// Move block with id 1 to a new parent
// 	movePayload := movePayload{NewParentId: 2, Index: 0}
// 	err = store.move(1, movePayload)
// 	assert.NoError(t, err)

// 	// Verify that the block is moved
// 	_, _, newParentMap, err := store.findBlockById(1)
// 	assert.NoError(t, err)
// 	assert.Equal(t, block2, newParentMap.Get(1))
// }

func TestInMemoryStoreFetch(t *testing.T) {
	store := NewInMemoryStore()

	block1 := blockRequest{Content: "Block 1"}
	block2 := blockRequest{Content: "Block 2"}

	payload := insertPayload{
		InsertOperations: []insertOperation{
			{ParentBlockId: root, Index: 0, Block: block1},
			{ParentBlockId: root, Index: 1, Block: block2},
		},
	}

	err := store.insert(payload)
	assert.NoError(t, err)

	// Fetch blocks with ids 1 and 2
	fetchedBlocks := store.fetch([]uint64{1, 2})

	// Verify that the correct blocks are fetched
	assert.Equal(t, 2, len(fetchedBlocks))
	assert.Contains(t, fetchedBlocks, block1)
	assert.Contains(t, fetchedBlocks, block2)
}
