package crafttask

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type Server struct {
	store *InMemoryStore
}

func NewServer(store *InMemoryStore) Server {
	return Server{
		store,
	}
}

// InsertBlocks inserts a list of new blocks to the document
func (s Server) InsertBlocks(w http.ResponseWriter, r *http.Request) {
	var insertPayload insertPayload
	err := json.NewDecoder(r.Body).Decode(&insertPayload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// validate payload
	s.store.insert(insertPayload) // return value

	w.WriteHeader(http.StatusCreated)
}

// DeleteBlocks deletes a list of existing blocks from the document
func (s Server) DeleteBlocks(w http.ResponseWriter, r *http.Request) {
	idsRaw := mux.Vars(r)["blockIds"]
	idsSplit := strings.Split(idsRaw, ",")
	ids := make([]uint64, len(idsSplit))
	for _, idRaw := range idsSplit {
		id, err := strconv.ParseUint(idRaw, 10, 64)
		if err != nil {
			http.Error(w, "block id parameter not an id", http.StatusBadRequest)
			return
		}
		ids = append(ids, id)
	}
	// validate payload
	s.store.delete(ids)

	w.WriteHeader(http.StatusNoContent)
}

// FetchBlocksByID fetches a list of existing blocks by their ID from the document
func (s Server) FetchBlocksByID(w http.ResponseWriter, r *http.Request) {
	idsRaw := mux.Vars(r)["blockIds"]
	idsSplit := strings.Split(idsRaw, ",")
	ids := make([]uint64, len(idsSplit))
	for _, idRaw := range idsSplit {
		id, err := strconv.ParseUint(idRaw, 10, 64)
		if err != nil {
			http.Error(w, "block id parameter not an id", http.StatusBadRequest)
			return
		}
		ids = append(ids, id)
	}
	blocks := s.store.fetch(ids)
	// Respond with the fetched blocks
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blocksToResponse(blocks))
	w.WriteHeader(http.StatusOK)
}

// DuplicateBlock duplicates an existing block with all of its subblocks
func (s Server) DuplicateBlock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idRaw := params["id"]
	id, err := strconv.ParseUint(idRaw, 10, 64)
	if err != nil {
		http.Error(w, "block id paramter", http.StatusBadRequest)
		return
	}
	// validate payload

	block, err := s.store.duplicate(id)
	if err != nil {
		if errors.Is(err, errBlockDoesNotExist) {
			http.Error(w, "block to duplicate does not exist", http.StatusNotFound)
			return
		}
		http.Error(w, "unexpected error occured", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(blockToResponse(block))
	w.WriteHeader(http.StatusOK)
}

// MoveBlock moves an existing block to another position in the document
func (s Server) MoveBlock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idRaw := params["id"]
	id, parseErr := strconv.ParseUint(idRaw, 10, 64)
	if parseErr != nil {
		http.Error(w, "block id paramter", http.StatusBadRequest)
		return
	}
	var movePayload movePayload
	decodingErr := json.NewDecoder(r.Body).Decode(&movePayload)
	if decodingErr != nil {
		http.Error(w, decodingErr.Error(), http.StatusBadRequest)
		return
	}
	// validate payload

	err := s.store.move(id, movePayload)
	if err != nil {
		if errors.Is(err, errBlockDoesNotExist) {
			http.Error(w, "block to move does not exist", http.StatusNotFound)
			return
		} else if errors.Is(err, errParentBlockDoesNotExist) {
			http.Error(w, "parent block to move to is invalid or does not exist", http.StatusBadRequest)
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

// ExportDocument exports the full document to a single string
func (s Server) ExportDocument(w http.ResponseWriter, r *http.Request) {
	// Implement logic to export the full document to a single string
	// ...

	// Respond with the exported document
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, "Exported document content")
}

func blocksToResponse(blocks []block) []blockResponse {
	toReturn := make([]blockResponse, len(blocks))
	for _, block := range blocks {
		reponse := blockToResponse(block)
		toReturn = append(toReturn, reponse)
	}
	return toReturn
}

func blockToResponse(block block) blockResponse {
	return blockResponse{
		Id:        block.id,
		Content:   block.content,
		Subblocks: blocksToResponse(block.subblocks.OrderedValues()),
	}
}
