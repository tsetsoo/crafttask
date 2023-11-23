package crafttask

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryStore_Insert(t *testing.T) {
	store := NewInMemoryStore()

	blockRequest1 := blockRequest{Content: "Block 1"}
	blockRequest2 := blockRequest{Content: "Block 2"}

	payload := []insertOperation{
		{ParentBlockId: root, Index: 0, Block: blockRequest1},
		{ParentBlockId: root, Index: 1, Block: blockRequest2},
	}

	_, err := store.InsertBlocks(payload)
	require.NoError(t, err)

	block1Retrieved, _, _, err := store.findBlockById(1)
	require.NoError(t, err)
	assert.Equal(t, blockRequest1.Content, block1Retrieved.content)

	block2Retrieved, _, _, err := store.findBlockById(2)
	require.NoError(t, err)
	assert.Equal(t, blockRequest2.Content, block2Retrieved.content)
}

func TestInMemoryStore_Delete(t *testing.T) {
	store := NewInMemoryStore()

	blockRequest1 := blockRequest{Content: "Block 1"}
	blockRequest2 := blockRequest{Content: "Block 2"}

	payload := []insertOperation{
		{ParentBlockId: root, Index: 0, Block: blockRequest1},
		{ParentBlockId: root, Index: 1, Block: blockRequest2},
	}

	_, err := store.InsertBlocks(payload)
	require.NoError(t, err)

	store.DeleteBlocks([]id{1})

	_, _, _, err = store.findBlockById(1)
	assert.Equal(t, errBlockDoesNotExist, err)
}

func TestInMemoryStore_Duplicate(t *testing.T) {
	store := NewInMemoryStore()

	blockToDuplicateIndex := 0
	blockRequest1 := blockRequest{Content: "Block 1"}
	payload := []insertOperation{
		{ParentBlockId: root, Index: blockToDuplicateIndex, Block: blockRequest1},
	}

	blocks, err := store.InsertBlocks(payload)
	require.NoError(t, err)

	childBlock := blockRequest{Content: "Child Block 1"}
	payload = []insertOperation{
		{ParentBlockId: blocks[0].id, Index: 0, Block: childBlock},
	}

	_, err = store.InsertBlocks(payload)
	require.NoError(t, err)

	duplicatedBlock, err := store.DuplicateBlock(1)
	require.NoError(t, err)

	require.NoError(t, err)
	assert.Equal(t, blockRequest1.Content, duplicatedBlock.content)
	require.Len(t, duplicatedBlock.subblocks.values, 1)
	assert.Equal(t, childBlock.Content, duplicatedBlock.subblocks.OrderedValues()[0].content)

	require.Len(t, store.document.blocks.keys, 2)
	assert.Equal(t, id(1), store.document.blocks.keys[blockToDuplicateIndex])
	assert.Equal(t, duplicatedBlock.id, store.document.blocks.keys[blockToDuplicateIndex+1])
}

func TestInMemoryStore_Move(t *testing.T) {
	store := NewInMemoryStore()

	blockRequest1 := blockRequest{Content: "Block 1"}
	blockRequest2 := blockRequest{Content: "Block 2"}

	payload := []insertOperation{
		{ParentBlockId: root, Index: 0, Block: blockRequest1},
		{ParentBlockId: root, Index: 1, Block: blockRequest2},
	}

	blocks, err := store.InsertBlocks(payload)
	require.NoError(t, err)

	childBlockRequest1 := blockRequest{Content: "Child Block 1"}
	childBlockRequest2 := blockRequest{Content: "Child Block 2"}

	payload = []insertOperation{
		{ParentBlockId: blocks[0].id, Index: 0, Block: childBlockRequest1},
		{ParentBlockId: blocks[1].id, Index: 1, Block: childBlockRequest2},
	}

	childBlocks, err := store.InsertBlocks(payload)
	require.NoError(t, err)

	grandChildBlockRequest1 := blockRequest{Content: "Grand Child Block 1"}
	grandChildBlockRequest2 := blockRequest{Content: "Grand Child Block 2"}

	payload = []insertOperation{
		{ParentBlockId: childBlocks[0].id, Index: 0, Block: grandChildBlockRequest1},
		{ParentBlockId: childBlocks[1].id, Index: 1, Block: grandChildBlockRequest2},
	}

	grandChildrenBlocks, err := store.InsertBlocks(payload)
	require.NoError(t, err)

	movedBlockId := childBlocks[0].id
	newParentBlockId := childBlocks[1].id
	movePayload := movePayload{NewParentId: newParentBlockId, Index: 0}
	err = store.MoveBlock(movedBlockId, movePayload)
	require.NoError(t, err)

	assert.Equal(t, newParentBlockId, store.parentsCache[movedBlockId])
	newParentBlock, _, _, err := store.findBlockById(newParentBlockId)
	require.NoError(t, err)
	require.Len(t, newParentBlock.subblocks.values, 2)
	assert.Equal(t, newParentBlock.subblocks.keys[0], movedBlockId)
	assert.Equal(t, newParentBlock.subblocks.keys[1], grandChildrenBlocks[1].id)

	require.Len(t, newParentBlock.subblocks.OrderedValues()[0].subblocks.OrderedValues(), 1)
	assert.Equal(t, newParentBlock.subblocks.OrderedValues()[0].subblocks.OrderedValues()[0], grandChildrenBlocks[0])
}

func TestInMemoryStore_MoveParentToChild_Err(t *testing.T) {
	store := NewInMemoryStore()

	blockRequest1 := blockRequest{Content: "Block 1"}
	blockRequest2 := blockRequest{Content: "Block 2"}

	payload := []insertOperation{
		{ParentBlockId: root, Index: 0, Block: blockRequest1},
		{ParentBlockId: root, Index: 1, Block: blockRequest2},
	}

	blocks, err := store.InsertBlocks(payload)
	require.NoError(t, err)

	childBlockRequest1 := blockRequest{Content: "Child Block 1"}
	childBlockRequest2 := blockRequest{Content: "Child Block 2"}

	payload = []insertOperation{
		{ParentBlockId: blocks[0].id, Index: 0, Block: childBlockRequest1},
		{ParentBlockId: blocks[1].id, Index: 1, Block: childBlockRequest2},
	}

	childBlocks, err := store.InsertBlocks(payload)
	require.NoError(t, err)

	grandChildBlockRequest1 := blockRequest{Content: "Grand Child Block 1"}
	grandChildBlockRequest2 := blockRequest{Content: "Grand Child Block 2"}

	payload = []insertOperation{
		{ParentBlockId: childBlocks[0].id, Index: 0, Block: grandChildBlockRequest1},
		{ParentBlockId: childBlocks[1].id, Index: 1, Block: grandChildBlockRequest2},
	}

	grandChildrenBlocks, err := store.InsertBlocks(payload)
	require.NoError(t, err)

	movedBlockId := childBlocks[0].id
	newParentBlockId := grandChildrenBlocks[0].id
	movePayload := movePayload{NewParentId: newParentBlockId, Index: 0}
	err = store.MoveBlock(movedBlockId, movePayload)
	assert.Equal(t, errBlockMovedToItsChild, err)
}

func TestInMemoryStore_Fetch(t *testing.T) {
	store := NewInMemoryStore()

	blockRequest1 := blockRequest{Content: "Block 1"}
	blockRequest2 := blockRequest{Content: "Block 2"}
	blockRequest3 := blockRequest{Content: "Block 3"}

	payload := []insertOperation{
		{ParentBlockId: root, Index: 0, Block: blockRequest1},
		{ParentBlockId: root, Index: 1, Block: blockRequest2},
		{ParentBlockId: root, Index: 2, Block: blockRequest3},
	}

	_, err := store.InsertBlocks(payload)
	require.NoError(t, err)

	fetchedBlocks := store.FetchBlocks([]id{1, 2})

	require.Len(t, fetchedBlocks, 2)
	assert.Equal(t, blockRequest1.Content, fetchedBlocks[0].content)
	assert.Equal(t, blockRequest2.Content, fetchedBlocks[1].content)
}

func TestInMemoryStore_Export(t *testing.T) {
	store := NewInMemoryStore()

	blockRequest1 := blockRequest{Content: "Block 1"}
	blockRequest2 := blockRequest{Content: "Block 2"}

	payload := []insertOperation{
		{ParentBlockId: root, Index: 0, Block: blockRequest1},
		{ParentBlockId: root, Index: 1, Block: blockRequest2},
	}

	blocks, err := store.InsertBlocks(payload)
	require.NoError(t, err)

	childBlockRequest1 := blockRequest{Content: "Child Block 1"}
	childBlockRequest2 := blockRequest{Content: "Child Block 2"}

	payload = []insertOperation{
		{ParentBlockId: blocks[0].id, Index: 0, Block: childBlockRequest1},
		{ParentBlockId: blocks[1].id, Index: 1, Block: childBlockRequest2},
	}

	childBlocks, err := store.InsertBlocks(payload)
	require.NoError(t, err)

	grandChildBlockRequest1 := blockRequest{Content: "Grand Child Block 1"}
	grandChildBlockRequest2 := blockRequest{Content: "Grand Child Block 2"}

	payload = []insertOperation{
		{ParentBlockId: childBlocks[0].id, Index: 0, Block: grandChildBlockRequest1},
		{ParentBlockId: childBlocks[1].id, Index: 1, Block: grandChildBlockRequest2},
	}

	_, err = store.InsertBlocks(payload)
	require.NoError(t, err)

	result := store.Export()
	assert.Equal(t,
		`Block 1
  Child Block 1
    Grand Child Block 1
Block 2
  Child Block 2
    Grand Child Block 2
`,
		result)
}
