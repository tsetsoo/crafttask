package crafttask

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Server struct {
	api API
}

func NewServer(api API) Server {
	return Server{api: api}
}

// these all should be prefixed for the document, but since there is only one, there is no need now (most apparent on the export)
func (s Server) Run() error {
	r := mux.NewRouter()

	r.HandleFunc("/blocks/bulk-insert", s.api.InsertBlocks).Methods("POST")
	r.HandleFunc("/blocks", s.api.DeleteBlocks).Methods("DELETE")
	r.HandleFunc("/blocks", s.api.FetchBlocksByID).Methods("GET")
	r.HandleFunc("/blocks/{id}/duplicate", s.api.DuplicateBlock).Methods("POST")
	r.HandleFunc("/blocks/{id}/move", s.api.MoveBlock).Methods("POST")
	r.HandleFunc("/export", s.api.ExportDocument).Methods("GET")

	handler := cors.AllowAll().Handler(r)
	http.Handle("/", handler)

	fmt.Println("Server listening on :8080")
	return http.ListenAndServe(":8080", nil)
}
