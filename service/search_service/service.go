package search_service

import "gitlab.com/ProtectIdentity/pugcha-backend/models/repositories_interface"

type SearchService struct {
	SearchRepo repositories_interface.SearchRepository
}

func NewSearch(repo repositories_interface.SearchRepository) *SearchService {
	return &SearchService{SearchRepo: repo}
}

func (s *SearchService) SearchTerm(term string) (interface{}, error) {
	output, err := s.SearchRepo.Search(term)
	if err != nil {
		return nil, err
	}

	return output, nil
}
