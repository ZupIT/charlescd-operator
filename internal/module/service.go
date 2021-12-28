package module

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"

	deployv1alpha1 "github.com/tiagoangelozup/charles-alpha/api/v1alpha1"
)

type Service struct {
	Client client.Client
}

func (s *Service) GetModule(ctx context.Context, key client.ObjectKey) (*deployv1alpha1.Module, error) {
	m := new(deployv1alpha1.Module)
	if err := s.Client.Get(ctx, key, m); err != nil {
		return nil, fmt.Errorf("failed to lookup resource: %w", err)
	}
	return m, nil
}
