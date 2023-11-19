package main

import (
	"fmt"
	"local/CraftTask/crafttask"
	"net/http"

	"github.com/gorilla/mux"
)

// these all should be prefixed for the document, but since there is only one, there is no need now (most apparent on the export)
func main() {
	r := mux.NewRouter()

	server := crafttask.NewServer(crafttask.NewInMemoryStore())

	r.HandleFunc("/blocks", server.InsertBlocks).Methods("POST")
	r.HandleFunc("/blocks", server.DeleteBlocks).Methods("DELETE")
	r.HandleFunc("/blocks", server.FetchBlocksByID).Methods("GET")
	r.HandleFunc("/blocks/{id}/duplicate", server.DuplicateBlock).Methods("POST")
	r.HandleFunc("/blocks/{id}/move", server.MoveBlock).Methods("POST")
	r.HandleFunc("/export", server.ExportDocument).Methods("GET")

	http.Handle("/", r)

	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", nil)
}
