package crafttask

const root = id(0)

type block struct {
	id        id
	content   string
	subblocks *orderedMapOfBlocks
}

type document struct {
	blocks *orderedMapOfBlocks
}
