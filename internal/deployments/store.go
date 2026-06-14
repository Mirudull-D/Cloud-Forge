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
func (s *Store) UpdateDeploymentStatusToBuilding(
	ctx context.Context, deploymentID uuid.UUID,
) (generated.Deployment, error) {
	deployment, err := s.queries.ChangeDeploymentStatusToBuilding(ctx, deploymentID)
	if err != nil {
		return generated.Deployment{}, err
	}
	return deployment, nil
}
func (s *Store) UpdateDeploymentStatusToRunning(
	ctx context.Context,
	id uuid.UUID,
	imageName, containerId string,
	port int) error {
	args := generated.ChangeDeploymentStatusToRunningParams{
		ID: id,
		ImageName: sql.NullString{
			String: imageName,
			Valid:  true,
		},
		ContainerID: sql.NullString{
			String: containerId,
			Valid:  true,
		},
		Port: sql.NullInt32{
			Int32: int32(port),
			Valid: true,
		},
	}

	_, err := s.queries.ChangeDeploymentStatusToRunning(ctx, args)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetNextAvailablePort(
	ctx context.Context,
) (int, error) {

	var port int

	err := s.db.QueryRowContext(
		ctx,
		`
		SELECT COALESCE(MAX(port), 8081)
		FROM deployments
		WHERE port IS NOT NULL
		`,
	).Scan(&port)

	if err != nil {
		return 0, err
	}

	return port + 1, nil
}

func (s *Store) UpdateDeploymentStatusToFailed(
	ctx context.Context, deploymentID uuid.UUID, errorMsg string) error {
	_, err := s.queries.ChangeDeploymentStatusToFailed(ctx, generated.ChangeDeploymentStatusToFailedParams{
		ID: deploymentID,
		ErrorMessage: sql.NullString{
			String: errorMsg,
			Valid:  true,
		},
	})
	if err != nil {
		return err
	}
	return nil
}
