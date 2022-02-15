package module

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/angelokurtis/reconciler"
	"github.com/go-logr/logr"
	"github.com/hashicorp/go-getter"
	mf "github.com/manifestival/manifestival"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	deployv1alpha1 "github.com/tiagoangelozup/charles-alpha/api/v1alpha1"
	"github.com/tiagoangelozup/charles-alpha/internal/tracing"
)

type (
	HelmClient interface {
		Template(ctx context.Context, name, chart string, values runtime.RawExtension) (mf.Manifest, error)
	}
	HelmValidation struct {
		reconciler.Funcs
		helm HelmClient
	}
)

func NewHelmValidation(helm HelmClient) *HelmValidation {
	return &HelmValidation{helm: helm}
}

func (h *HelmValidation) Reconcile(ctx context.Context, obj client.Object) (ctrl.Result, error) {
	if module, ok := obj.(*deployv1alpha1.Module); ok {
		return h.reconcile(ctx, module)
	}
	return h.Next(ctx, obj)
}

func (h *HelmValidation) reconcile(ctx context.Context, module *deployv1alpha1.Module) (ctrl.Result, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()
	log := logr.FromContextOrDiscard(ctx)

	if module.Status.Source == nil || module.Status.Source.Path == "" {
		return h.Next(ctx, module)
	}

	origin := module.Status.Source.Path
	destination := origin[0 : len(origin)-len(".tar.gz")]

	if module.Spec.Repository.Git != nil && module.Spec.Repository.Git.Path != "" {
		origin += "//" + module.Spec.Repository.Git.Path
		destination = filepath.Join(destination, module.Spec.Repository.Git.Path)
	}

	if err := getter.GetAny(destination, origin); err != nil {
		return h.RequeueOnErr(ctx, fmt.Errorf("error extracting Source artifact: %w", err))
	}

	manifest, err := h.helm.Template(ctx, module.GetName(), destination, module.Spec.Values)
	if err != nil {
		return h.RequeueOnErr(ctx, fmt.Errorf("error rendering Helm chart templates: %w", err))
	}

	if _, err = manifest.DryRun(); err != nil {
		return h.RequeueOnErr(ctx, fmt.Errorf("error rendering Helm chart templates: %w", err))
	}

	log.Info("Helm chart successfully rendered")
	return h.Next(ctx, module)
}
