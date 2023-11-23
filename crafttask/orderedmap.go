package crafttask

type orderedMapOfBlocks struct {
	keys   []id
	values map[id]block
}

func NewOrderedMapOfBlocks() *orderedMapOfBlocks {
	o := orderedMapOfBlocks{
		keys:   make([]id, 0),
		values: make(map[id]block),
	}
	return &o
}

func (o *orderedMapOfBlocks) Get(key id) (block, bool) {
	val, exists := o.values[key]
	return val, exists
}

func (o *orderedMapOfBlocks) GetAndIndex(key id) (block, int, bool) {
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

func (o *orderedMapOfBlocks) Set(key id, value block) {
	_, exists := o.values[key]
	if !exists {
		o.keys = append(o.keys, key)
	}
	o.values[key] = value
}

func (o *orderedMapOfBlocks) Insert(key id, index int, value block) {
	_, exists := o.values[key]
	if !exists {
		o.keys = insert(o.keys, index, key)
	}
	o.values[key] = value
}

func insert(a []id, index int, value id) []id {
	if index >= len(a) { // to handle sometimes incosistent state if we try to insert at an invalid index, we just append instead
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...)
	a[index] = value
	return a
}

func (o *orderedMapOfBlocks) Delete(key id) {
	_, ok := o.values[key]
	if !ok {
		return
	}
	for i, k := range o.keys {
		if k == key {
			o.keys = append(o.keys[:i], o.keys[i+1:]...)
			break
		}
	}
	delete(o.values, key)
}

func (o *orderedMapOfBlocks) Keys() []id {
	return o.keys
}

func (o *orderedMapOfBlocks) Values() map[id]block {
	return o.values
}

func (o *orderedMapOfBlocks) OrderedValues() []block {
	toReturn := make([]block, 0, len(o.values))
	for _, key := range o.keys {
		toReturn = append(toReturn, o.values[key])
	}
	return toReturn
}
