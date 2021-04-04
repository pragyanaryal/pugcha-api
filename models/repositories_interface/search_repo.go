package repositories_interface

type SearchRepository interface {
	Search(string) (interface{}, error)
}

