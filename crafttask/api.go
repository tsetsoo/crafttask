package crafttask

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type API struct {
	store Store
}

func NewAPI(store Store) API {
	return API{
		store,
	}
}

// InsertBlocks inserts a list of new blocks to the document
func (s API) InsertBlocks(w http.ResponseWriter, r *http.Request) {
	var insertPayload []insertOperation
	decodeErr := json.NewDecoder(r.Body).Decode(&insertPayload)
	if decodeErr != nil {
		http.Error(w, decodeErr.Error(), http.StatusBadRequest)
		return
	}
	// validate payload
	blocks, insertErr := s.store.InsertBlocks(insertPayload)
	if insertErr != nil {
		http.Error(w, insertErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(blocksToResponse(blocks))
}

func (s API) DeleteBlocks(w http.ResponseWriter, r *http.Request) {
	idsRaw := r.URL.Query().Get("blockIds")
	idsSplit := strings.Split(idsRaw, ",")
	ids := make([]id, len(idsSplit))
	for _, idRaw := range idsSplit {
		id, err := idFromString(idRaw)
		if err != nil {
			http.Error(w, "block id parameter not an id", http.StatusBadRequest)
			return
		}
		ids = append(ids, id)
	}
	// validate payload
	s.store.DeleteBlocks(ids)

	w.WriteHeader(http.StatusNoContent)
}

func (s API) FetchBlocksByID(w http.ResponseWriter, r *http.Request) {
	idsRaw := r.URL.Query().Get("blockIds")
	idsSplit := strings.Split(idsRaw, ",")
	ids := make([]id, len(idsSplit))
	for _, idRaw := range idsSplit {
		id, err := idFromString(idRaw)
		if err != nil {
			http.Error(w, "block id parameter not an id", http.StatusBadRequest)
			return
		}
		ids = append(ids, id)
	}
	blocks := s.store.FetchBlocks(ids)
	// Respond with the fetched blocks
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(blocksToResponse(blocks))
}

func (s API) DuplicateBlock(w http.ResponseWriter, r *http.Request) {
	idRaw := mux.Vars(r)["id"]
	id, err := idFromString(idRaw)
	if err != nil {
		http.Error(w, "block id paramter", http.StatusBadRequest)
		return
	}
	// validate payload

	block, err := s.store.DuplicateBlock(id)
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
}

func (s API) MoveBlock(w http.ResponseWriter, r *http.Request) {
	idRaw := mux.Vars(r)["id"]
	id, parseErr := idFromString(idRaw)
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
	// also validate blockToMove != newParentId

	err := s.store.MoveBlock(id, movePayload)
	if err != nil {
		if errors.Is(err, errBlockDoesNotExist) {
			http.Error(w, "block to move does not exist", http.StatusNotFound)
			return
		} else if errors.Is(err, errParentBlockDoesNotExist) {
			http.Error(w, "parent block to move to is invalid or does not exist", http.StatusBadRequest)
			return
		} else if errors.Is(err, errBlockMovedToItsChild) {
			http.Error(w, "block to move is a parent of the block to move to", http.StatusBadRequest)
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s API) ExportDocument(w http.ResponseWriter, r *http.Request) {
	content := s.store.Export()
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, content)
}

func blocksToResponse(blocks []block) []blockResponse {
	toReturn := make([]blockResponse, 0, len(blocks))
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
