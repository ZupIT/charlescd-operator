package client

import (
	"context"
	"fmt"

	"github.com/angelokurtis/reconciler"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	deployv1alpha1 "github.com/tiagoangelozup/charles-alpha/api/v1alpha1"
	"github.com/tiagoangelozup/charles-alpha/internal/tracing"
)

type Module struct {
	reconciler.Funcs
	client client.Client
}

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

func (s *Module) UpdateModuleStatus(ctx context.Context, module *deployv1alpha1.Module) (ctrl.Result, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()
	log := logr.FromContextOrDiscard(ctx)

	m, err := s.GetModule(ctx, types.NamespacedName{Namespace: module.GetNamespace(), Name: module.GetName()})
	if err != nil {
		return s.RequeueOnErr(ctx, fmt.Errorf("failed to update Module status: %w", err))
	}

	patch := client.MergeFrom(m.DeepCopy())
	m.Status = module.Status
	diff, _ := patch.Data(module)
	log.Info("Status changed", "diff", string(diff))

	if err = s.client.Status().Patch(ctx, m, patch); err != nil {
		return s.RequeueOnErr(ctx, fmt.Errorf("failed to update Module status: %w", err))
	}

	return s.Finish(ctx)
}
