package module

import (
	"context"

	deployv1alpha1 "github.com/tiagoangelozup/charles-alpha/api/v1alpha1"
)

type StatusWriter interface {
	UpdateModuleStatus(ctx context.Context, module *deployv1alpha1.Module) error
}
