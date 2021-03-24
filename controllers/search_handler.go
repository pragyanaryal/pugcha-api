package controllers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type SearchHandler struct{}

func NewSearch() *SearchHandler { return &SearchHandler{} }

func (add *SearchHandler) Search(rw http.ResponseWriter, r *http.Request) {
	term := mux.Vars(r)["term"]
	fmt.Println(term)
}
