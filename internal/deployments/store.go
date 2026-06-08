package deployments

import (
	"CloudHub/db/generated"
	"context"
	"database/sql"
)

type Store struct {
	db      *sql.DB
	queries *generated.Queries
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		queries: generated.New(db),
	}
}

func (s *Store) CreateNewDeployment(ctx context.Context, gitUrl string) (generated.Deployment, error) {
	newDeployment, err := s.queries.CreateNewDeployment(ctx, gitUrl)
	if err != nil {
		return generated.Deployment{}, err
	}
	return newDeployment, nil

}
