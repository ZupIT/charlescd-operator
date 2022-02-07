package client

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"

	deployv1alpha1 "github.com/tiagoangelozup/charles-alpha/api/v1alpha1"
)

type Module struct{ client client.Client }

func NewModule(client client.Client) *Module {
	return &Module{client: client}
}

func (s *Module) GetModule(ctx context.Context, key client.ObjectKey) (*deployv1alpha1.Module, error) {
	m := new(deployv1alpha1.Module)
	if err := s.client.Get(ctx, key, m); err != nil {
		return nil, fmt.Errorf("failed to lookup resource: %w", err)
	}
	return m, nil
}

func (s *Module) UpdateModuleStatus(ctx context.Context, module *deployv1alpha1.Module) error {
	if err := s.client.Status().Update(ctx, module); err != nil {
		return fmt.Errorf("failed to update Module status: %w", err)
	}
	return nil
}
