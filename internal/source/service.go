package source

import (
	"context"
	"fmt"

	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Service struct {
	Client client.Client
}

func (s *Service) GetGitRepository(ctx context.Context, key client.ObjectKey) (*sourcev1.GitRepository, error) {
	m := new(sourcev1.GitRepository)
	if err := s.Client.Get(ctx, key, m); err != nil {
		return nil, fmt.Errorf("failed to lookup resource: %w", err)
	}
	return m, nil
}
