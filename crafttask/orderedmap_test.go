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
	assert.Equal(t, block1, val)
}

func TestInsertAndGetAndIndex(t *testing.T) {
	o := NewOrderedMapOfBlocks()
	block1 := block{content: "Block1"}
	block2 := block{content: "Block2"}

	o.Set(1, block1)
	o.Insert(2, 0, block2)

	assert.Equal(t, []id{2, 1}, o.Keys())

	val, index, exists := o.GetAndIndex(2)
	assert.True(t, exists)
	assert.Equal(t, block2, val)
	assert.Equal(t, 0, index)
}

func TestDelete(t *testing.T) {
	o := NewOrderedMapOfBlocks()
	block1 := block{content: "Block1"}

	o.Set(1, block1)
	o.Delete(1)

	_, exists := o.Get(1)
	assert.False(t, exists)
}

func TestOrderedValues(t *testing.T) {
	o := NewOrderedMapOfBlocks()
	block1 := block{content: "Block1"}
	block2 := block{content: "Block2"}

	o.Set(1, block1)
	o.Set(2, block2)

	expectedOrderedValues := []block{block1, block2}

	assert.Equal(t, expectedOrderedValues, o.OrderedValues())
}
