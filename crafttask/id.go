package crafttask

import (
	"strconv"
	"sync/atomic"
)

type id uint64

type idGenerator interface {
	getNewId() id
}

type inMemoryIdGenerator struct {
	currentId atomic.Uint64
}

func newInMemoryIdGenerator() *inMemoryIdGenerator {
	return &inMemoryIdGenerator{}
}

func (i *inMemoryIdGenerator) getNewId() id {
	return id(i.currentId.Add(1))
}

func idFromString(rawId string) (id, error) {
	parsedId, err := strconv.ParseUint(rawId, 10, 64)
	if err != nil {
		return 0, err
	}
	return id(parsedId), nil
}
