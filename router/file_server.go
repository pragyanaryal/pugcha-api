package router

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"path/filepath"
)

func SetupFileRoutes(r *mux.Router, file string, host string) {

	fileServer := http.FileServer(neuteredFileSystem{http.Dir(file)})
	f := r.Host(host).Schemes("https").Methods(http.MethodOptions, http.MethodGet).Subrouter()
	f.PathPrefix("/contents/").Handler(http.StripPrefix("/contents/", fileServer)).Methods("GET", "OPTIONS")

	headersOk := handlers.AllowedHeaders([]string{"*"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "OPTIONS"})

	f.Use(handlers.CORS(headersOk, originsOk, methodsOk))
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}
	s, err := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}
	return f, nil
}
