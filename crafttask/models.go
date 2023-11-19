package crafttask

const root = uint64(0)

type block struct {
	id        uint64
	content   string
	subblocks *orderedMapOfBlocks
}

type document struct {
	blocks *orderedMapOfBlocks
}
