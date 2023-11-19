package crafttask

// biderectional link
// synchronized access or simpler data strcutures?
type InMemoryStore struct {
	document     document
	parentsCache map[uint64]uint64
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		document: document{
			blocks: NewOrderedMapOfBlocks(),
		},
		parentsCache: make(map[uint64]uint64),
	}
}

// maybe we can store this instead of just a map of int to int
func (st *InMemoryStore) pathToNode(parentId uint64) ([]uint64, error) {
	pathToNode := make([]uint64, 0)
	if parentId == root {
		return pathToNode, nil
	}
	pathToNode = append(pathToNode, parentId)
	for {
		parent, hasParent := st.parentsCache[parentId]
		if !hasParent {
			return nil, errBlockDoesNotExist
		}
		if parent == root {
			break
		}
		pathToNode = append(pathToNode, parent)
	}
	return pathToNode, nil
}

// assumes very large document means document with 100-1000 levels deep, not more; to discuss?
func (st *InMemoryStore) findMapByParent(parentId uint64) (*orderedMapOfBlocks, error) {
	mapToReturn := st.document.blocks
	blockParentList, err := st.pathToNode(parentId)
	if err != nil {
		return nil, err
	}
	for i := len(blockParentList) - 1; i >= 0; i-- {
		block, ok := mapToReturn.Get(blockParentList[i])
		if !ok {
			panic("inconsistent internal state")
		}
		mapToReturn = block.subblocks
	}
	return mapToReturn, nil
}

func (st *InMemoryStore) findBlockById(blockId uint64) (block, int, *orderedMapOfBlocks, error) {
	parentOfBlock, hasParent := st.parentsCache[blockId]
	if !hasParent {
		return block{}, 0, nil, errBlockDoesNotExist
	}
	mapWhereBlockIsLocated, err := st.findMapByParent(parentOfBlock)
	if err != nil {
		return block{}, 0, nil, err
	}
	blockToReturn, index, ok := mapWhereBlockIsLocated.GetAndIndex(blockId)
	if !ok {
		panic("inconsistent internal state")
	}
	return blockToReturn, index, mapWhereBlockIsLocated, nil
}

func (st *InMemoryStore) insert(insertPayload insertPayload) error {
	for _, insertOperation := range insertPayload.InsertOperations {
		mapToInsertIn, err := st.findMapByParent(insertOperation.ParentBlockId)
		if err != nil {
			return err
		}
		blockId := getNewId()
		blockToAdd := block{
			id:        blockId,
			content:   insertOperation.Block.Content,
			subblocks: NewOrderedMapOfBlocks(),
		}
		st.parentsCache[blockId] = insertOperation.ParentBlockId
		mapToInsertIn.Insert(blockId, insertOperation.Index, blockToAdd)
	}
	return nil
}

func (st *InMemoryStore) delete(idsToDelete []uint64) {
	for _, blockIdToDelete := range idsToDelete {
		blockToDelete, _, mapToDeleteFrom, err := st.findBlockById(blockIdToDelete)
		if err != nil {
			continue
			// already delete can continue
		}
		mapToDeleteFrom.Delete(blockIdToDelete)
		delete(st.parentsCache, blockIdToDelete)
		st.recursiveDeleteParentLinks(blockToDelete)
	}
}

func (st *InMemoryStore) recursiveDeleteParentLinks(blockToDelete block) {
	for subblockId, subblockToUnlink := range blockToDelete.subblocks.values {
		delete(st.parentsCache, subblockId)
		st.recursiveDeleteParentLinks(subblockToUnlink)
	}
}

func (st *InMemoryStore) fetch(idsToFetch []uint64) []block {
	toReturn := make([]block, len(idsToFetch))
	for _, id := range idsToFetch {
		blockToReturn, _, _, err := st.findBlockById(id)
		if err != nil {
			continue
			// we are filtering here so I think we shouldn't error out
		}
		toReturn = append(toReturn, blockToReturn)
	}
	return toReturn
}

func (st *InMemoryStore) duplicate(idToDuplicate uint64) (block, error) {
	blockToDuplicate, index, mapToDuplicateIn, err := st.findBlockById(idToDuplicate)
	if err != nil {
		return block{}, err
	}
	parentOfBlock, hasParent := st.parentsCache[idToDuplicate]
	if !hasParent {
		panic("inconsistent internal state")
	}
	newId := st.recursiveSetNewIds(blockToDuplicate, parentOfBlock)
	mapToDuplicateIn.Insert(newId, index+1, blockToDuplicate)
	return blockToDuplicate, nil
}

func (st *InMemoryStore) recursiveSetNewIds(blockToSetIdsTo block, parentId uint64) uint64 {
	blockId := getNewId()
	blockToSetIdsTo.id = blockId
	st.parentsCache[blockId] = parentId
	for _, subblock := range blockToSetIdsTo.subblocks.values {
		st.recursiveSetNewIds(subblock, blockId)
	}
	return blockId
}

func (st *InMemoryStore) move(blockId uint64, movePayload movePayload) error {
	blockToMove, _, oldMap, err := st.findBlockById(blockId)
	if err != nil {
		return err
	}
	oldMap.Delete(blockId)
	st.parentsCache[blockId] = movePayload.NewParentId
	newMap, err := st.findMapByParent(movePayload.NewParentId)
	if err != nil {
		//log
		return errParentBlockDoesNotExist
	}
	newMap.Insert(blockId, movePayload.Index, blockToMove)
	return nil
}

func (st *InMemoryStore) export() {
	// this should be simple using the treelike struct
}
