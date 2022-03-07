package module

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"

	charlescdv1alpha1 "github.com/tiagoangelozup/charles-alpha/api/v1alpha1"
)

type StatusWriter interface {
	UpdateModuleStatus(ctx context.Context, module *charlescdv1alpha1.Module) (ctrl.Result, error)
}
