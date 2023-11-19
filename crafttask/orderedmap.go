package crafttask

type orderedMapOfBlocks struct {
	keys   []uint64
	values map[uint64]block
}

func NewOrderedMapOfBlocks() *orderedMapOfBlocks {
	o := orderedMapOfBlocks{
		keys:   make([]uint64, 0),
		values: make(map[uint64]block),
	}
	return &o
}

func (o *orderedMapOfBlocks) Get(key uint64) (block, bool) {
	val, exists := o.values[key]
	return val, exists
}

func (o *orderedMapOfBlocks) GetAndIndex(key uint64) (block, int, bool) {
	val, exists := o.values[key]
	if !exists {
		return block{}, 0, false
	}
	index := 0
	for i, k := range o.keys {
		if k == key {
			index = i
		}
	}
	return val, index, true
}

func (o *orderedMapOfBlocks) Set(key uint64, value block) {
	_, exists := o.values[key]
	if !exists {
		o.keys = append(o.keys, key)
	}
	o.values[key] = value
}

func (o *orderedMapOfBlocks) Insert(key uint64, index int, value block) {
	_, exists := o.values[key]
	if !exists {
		o.keys = insert(o.keys, index, key)
	}
	o.values[key] = value
}

func insert(a []uint64, index int, value uint64) []uint64 {
	if index >= len(a) { // to handle sometimes incosistent state if we try to insert at an invalid index, we just append instead
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...)
	a[index] = value
	return a
}

func (o *orderedMapOfBlocks) Delete(key uint64) {
	// check key is in use
	_, ok := o.values[key]
	if !ok {
		return
	}
	// remove from keys
	for i, k := range o.keys {
		if k == key {
			o.keys = append(o.keys[:i], o.keys[i+1:]...)
			break
		}
	}
	// remove from values
	delete(o.values, key)
}

func (o *orderedMapOfBlocks) Keys() []uint64 {
	return o.keys
}

func (o *orderedMapOfBlocks) Values() map[uint64]block {
	return o.values
}

func (o *orderedMapOfBlocks) OrderedValues() []block {
	toReturn := make([]block, 0, len(o.values))
	for _, key := range o.keys {
		toReturn = append(toReturn, o.values[key])
	}
	return toReturn
}
