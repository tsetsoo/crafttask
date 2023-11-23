package crafttask

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	o := NewOrderedMapOfBlocks()
	block1 := block{content: "Block1"}
	o.Set(1, block1)

	val, exists := o.Get(1)
	assert.True(t, exists, "Key should exist")
	assert.Equal(t, block1, val, "Values should be equal")
}

func TestInsertAndGetAndIndex(t *testing.T) {
	o := NewOrderedMapOfBlocks()
	block1 := block{content: "Block1"}
	block2 := block{content: "Block2"}

	o.Set(1, block1)
	o.Insert(2, 1, block2)

	keys := o.Keys()
	assert.Equal(t, []id{1, 2}, keys)

	val, index, exists := o.GetAndIndex(2)
	assert.True(t, exists, "Key should exist")
	assert.Equal(t, block2, val, "Values should be equal")
	assert.Equal(t, 1, index, "Index should be 1")
}

func TestDelete(t *testing.T) {
	o := NewOrderedMapOfBlocks()
	block1 := block{content: "Block1"}

	o.Set(1, block1)
	o.Delete(1)

	_, exists := o.Get(1)
	assert.False(t, exists, "Key should not exist after deletion")
}

func TestOrderedValues(t *testing.T) {
	o := NewOrderedMapOfBlocks()
	block1 := block{content: "Block1"}
	block2 := block{content: "Block2"}

	o.Set(1, block1)
	o.Set(2, block2)

	expectedOrderedValues := []block{block1, block2}

	assert.Equal(t, expectedOrderedValues, o.OrderedValues(), "OrderedValues should be equal")
}
