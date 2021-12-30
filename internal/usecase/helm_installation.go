package usecase

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"

	deployv1alpha1 "github.com/tiagoangelozup/charles-alpha/api/v1alpha1"
	"github.com/tiagoangelozup/charles-alpha/internal/runtime"
)

type HelmInstallation struct{}

func (hi *HelmInstallation) EnsureHelmInstallation(ctx context.Context, module *deployv1alpha1.Module) (ctrl.Result, error) {
	// helm install logic here
	return runtime.Finish()
}
