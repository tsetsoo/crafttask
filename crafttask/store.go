package crafttask

import (
	"fmt"
	"strings"
)

type Store interface {
	InsertBlocks(insertOperations []insertOperation) ([]block, error)
	DeleteBlocks(blocksIdsToDelete []id)
	FetchBlocks(blocksIdsToFetch []id) []block
	DuplicateBlock(blockToDuplicate id) (block, error)
	MoveBlock(blockToMove id, movePayload movePayload) error
	Export() string
}

// biderectional link
// synchronized access or simpler data strcutures?
type InMemoryStore struct {
	document     document
	parentsCache map[id]id
	idGenerator  idGenerator
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		document: document{
			blocks: NewOrderedMapOfBlocks(),
		},
		parentsCache: make(map[id]id),
		idGenerator:  newInMemoryIdGenerator(),
	}
}

// maybe we can store this instead of just a map of int to int
func (st *InMemoryStore) pathToNode(parentId id) ([]id, error) {
	pathToNode := make([]id, 0)
	for parentId != root {
		pathToNode = append(pathToNode, parentId)
		var ok bool
		parentId, ok = st.parentsCache[parentId]
		if !ok {
			return nil, errBlockDoesNotExist
		}
	}
	return pathToNode, nil
}

// assumes very large document means document with 100-1000 levels deep, not more; to discuss?
func (st *InMemoryStore) findMapByParent(parentId id) (*orderedMapOfBlocks, error) {
	mapToReturn := st.document.blocks
	blockParentList, err := st.pathToNode(parentId)
	if err != nil {
		return nil, err
	}
	for i := len(blockParentList) - 1; i >= 0; i-- {
		block, ok := mapToReturn.Get(blockParentList[i])
		if !ok {
			panic("inconsistent internal state") // this shouldn't happen, perhaps worth a rework
		}
		mapToReturn = block.subblocks
	}
	return mapToReturn, nil
}

// returning 4 parameters :(
func (st *InMemoryStore) findBlockById(blockId id) (block, int, *orderedMapOfBlocks, error) {
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
		panic("inconsistent internal state") // this shouldn't happen, perhaps worth a rework
	}
	return blockToReturn, index, mapWhereBlockIsLocated, nil
}

func (st *InMemoryStore) InsertBlocks(insertOperations []insertOperation) ([]block, error) {
	blocksToReturn := make([]block, 0, len(insertOperations))
	for _, insertOperation := range insertOperations {
		mapToInsertIn, err := st.findMapByParent(insertOperation.ParentBlockId)
		if err != nil {
			return nil, err // can be reworked to return partial success
		}
		blockId := st.idGenerator.getNewId()
		blockToAdd := block{
			id:        blockId,
			content:   insertOperation.Block.Content,
			subblocks: NewOrderedMapOfBlocks(),
		}
		blocksToReturn = append(blocksToReturn, blockToAdd)
		st.parentsCache[blockId] = insertOperation.ParentBlockId
		mapToInsertIn.Insert(blockId, insertOperation.Index, blockToAdd)
	}
	return blocksToReturn, nil
}

func (st *InMemoryStore) DeleteBlocks(idsToDelete []id) {
	for _, blockIdToDelete := range idsToDelete {
		blockToDelete, _, mapToDeleteFrom, err := st.findBlockById(blockIdToDelete)
		if err != nil {
			continue
			// already deleted; can continue
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

func (st *InMemoryStore) FetchBlocks(idsToFetch []id) []block {
	toReturn := make([]block, 0, len(idsToFetch))
	for _, id := range idsToFetch {
		blockToReturn, _, _, err := st.findBlockById(id)
		if err != nil {
			continue // we are filtering here so I think we shouldn't error out
		}
		toReturn = append(toReturn, blockToReturn)
	}
	return toReturn
}

func (st *InMemoryStore) DuplicateBlock(idToDuplicate id) (block, error) {
	blockToDuplicate, index, mapToDuplicateIn, err := st.findBlockById(idToDuplicate)
	if err != nil {
		return block{}, err
	}
	parentOfBlock, hasParent := st.parentsCache[idToDuplicate]
	if !hasParent {
		panic("inconsistent internal state") // this shouldn't happen, perhaps worth a rework
	}
	newId := st.recursiveSetNewIds(&blockToDuplicate, parentOfBlock)
	mapToDuplicateIn.Insert(newId, index+1, blockToDuplicate)
	return blockToDuplicate, nil
}

func (st *InMemoryStore) recursiveSetNewIds(blockToSetIdsTo *block, parentId id) id {
	blockId := st.idGenerator.getNewId()
	blockToSetIdsTo.id = blockId
	st.parentsCache[blockId] = parentId
	for _, subblock := range blockToSetIdsTo.subblocks.values {
		st.recursiveSetNewIds(&subblock, blockId)
	}
	return blockId
}

func (st *InMemoryStore) MoveBlock(blockId id, movePayload movePayload) error {
	consistencyCheckErr := st.blockMovedToItsChild(blockId, movePayload.NewParentId)
	if consistencyCheckErr != nil {
		return consistencyCheckErr
	}
	blockToMove, _, oldMap, findErr := st.findBlockById(blockId)
	if findErr != nil {
		return findErr
	}

	oldMap.Delete(blockId)
	st.parentsCache[blockId] = movePayload.NewParentId
	newMap, findMapErr := st.findMapByParent(movePayload.NewParentId)
	if findMapErr != nil {
		//log
		return errParentBlockDoesNotExist
	}
	newMap.Insert(blockId, movePayload.Index, blockToMove)

	return nil
}

func (st *InMemoryStore) blockMovedToItsChild(blockId, newParentId id) error {
	for newParentId != root {
		var ok bool
		newParentId, ok = st.parentsCache[newParentId]
		if !ok {
			return errBlockDoesNotExist
		}
		if newParentId == blockId {
			return errBlockMovedToItsChild
		}
	}
	return nil
}

func (st *InMemoryStore) Export() string {
	var builder strings.Builder
	for _, block := range st.document.blocks.OrderedValues() {
		addString(&builder, block, 0)
	}
	return builder.String()
}

func addString(builder *strings.Builder, block block, indentLevel int) {
	builder.WriteString(fmt.Sprintf("%s%s\n", strings.Repeat(" ", indentLevel*2), block.content))
	for _, subblock := range block.subblocks.OrderedValues() {
		addString(builder, subblock, indentLevel+1)
	}
}
