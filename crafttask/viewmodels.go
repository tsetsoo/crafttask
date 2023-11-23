package crafttask

type insertOperation struct {
	ParentBlockId id
	Block         blockRequest
	Index         int
}

type movePayload struct {
	NewParentId id
	Index       int
}

type blockResponse struct {
	Id        id
	Content   string
	Subblocks []blockResponse
}

type blockRequest struct {
	Content   string
	Subblocks []blockRequest //not implemented
}
