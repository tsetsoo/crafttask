package crafttask

type insertOperation struct {
	ParentBlockId uint64
	Block         blockRequest
	Index         int
}

type insertPayload struct {
	InsertOperations []insertOperation
}

type movePayload struct {
	NewParentId uint64
	Index       int
}

type blockResponse struct {
	Id        uint64
	Content   string
	Subblocks []blockResponse
}

type blockRequest struct {
	Content   string
	Subblocks []blockRequest //TODO not implemented
}
