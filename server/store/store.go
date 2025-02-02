package store

import "github.com/estnafinema0/Github-Trends-Aggregator/models"

type Store struct {
}

func NewStore() *Store {
	return &Store{}
}

func (s *Store) GetRepos(language string) []models.Repository {
	return []models.Repository{}
}
