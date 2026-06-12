package deployments

import (
	"CloudHub/db/generated"
	"context"
	"database/sql"

	"github.com/google/uuid"
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
func (s *Store) GetAllDeployments(ctx context.Context) ([]generated.Deployment, error) {
	allDeployments, err := s.queries.GetAllDeployments(ctx)
	if err != nil {
		return nil, err
	}
	return allDeployments, nil
}

func (s *Store) GetDeploymentById(ctx context.Context, deploymentID uuid.UUID) (generated.Deployment, error) {
	deployment, err := s.queries.GetDeploymentById(ctx, deploymentID)
	if err != nil {
		return generated.Deployment{}, err
	}
	return deployment, nil
}
func (s *Store) UpdateDeploymentStatusToBuilding(ctx context.Context, deploymentID uuid.UUID) (generated.Deployment, error) {
	deployment, err := s.queries.ChangeDeploymentStatus(ctx, deploymentID)
	if err != nil {
		return generated.Deployment{}, err
	}
	return deployment, nil
}
