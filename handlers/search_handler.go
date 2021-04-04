package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/ProtectIdentity/pugcha-backend/db/repositories"
	"gitlab.com/ProtectIdentity/pugcha-backend/erros"
	"gitlab.com/ProtectIdentity/pugcha-backend/service/search_service"
)

type SearchHandler struct {
}

func NewSearch() *SearchHandler { return &SearchHandler{} }

func (add *SearchHandler) Search(rw http.ResponseWriter, r *http.Request) {
	term := mux.Vars(r)["term"]

	service := search_service.NewSearch(repositories.SearchRepo)
	data, err := service.SearchTerm(term)
	if err != nil {
		erros.JSONError(rw, err, http.StatusNotFound)
		return
	}

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(rw).Encode(data)
}
