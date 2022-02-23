package module

import (
	"context"

	"github.com/angelokurtis/reconciler"
	"github.com/go-logr/logr"
	mf "github.com/manifestival/manifestival"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	deployv1alpha1 "github.com/tiagoangelozup/charles-alpha/api/v1alpha1"
	"github.com/tiagoangelozup/charles-alpha/internal/tracing"
)

type (
	ObjectConverter interface {
		FromUnstructured(in *unstructured.Unstructured, out interface{}) error
	}
	CheckComponents struct {
		reconciler.Funcs
		manifest ManifestReader
		object   ObjectConverter
		status   StatusWriter
	}
)

func NewCheckComponents(manifest ManifestReader, object ObjectConverter, status StatusWriter) *CheckComponents {
	return &CheckComponents{manifest: manifest, object: object, status: status}
}

func (c *CheckComponents) Reconcile(ctx context.Context, obj client.Object) (ctrl.Result, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	resources := resourcesFromContext(ctx)
	if module, ok := obj.(*deployv1alpha1.Module); ok && len(resources) > 0 {
		manifest, err := c.manifest.FromUnstructured(ctx, resources)
		if err != nil {
			return c.RequeueOnErr(ctx, err)
		}
		return c.reconcile(ctx, module, manifest)
	}
	return c.Next(ctx, obj)
}

func (c *CheckComponents) reconcile(ctx context.Context, module *deployv1alpha1.Module, manifest mf.Manifest) (ctrl.Result, error) {
	log := logr.FromContextOrDiscard(ctx)

	components := make([]*deployv1alpha1.Component, 0, 0)
	for _, u := range manifest.Filter(mf.ByKind("Deployment")).Resources() {
		deploy := &appsv1.Deployment{}
		if err := c.object.FromUnstructured(&u, deploy); err != nil {
			return c.RequeueOnErr(ctx, err)
		}
		component := &deployv1alpha1.Component{Name: deploy.GetName()}
		for _, container := range deploy.Spec.Template.Spec.Containers {
			component.Containers = append(component.Containers, &deployv1alpha1.Container{
				Name:  container.Name,
				Image: container.Image,
			})
		}
		components = append(components, component)
	}

	total := len(components)
	if total > 0 {
		log.Info("Deployable components found", "total", total)
	} else {
		log.Info("No deployable components were found")
	}

	if module.SetComponents(components) {
		return c.status.UpdateModuleStatus(ctx, module)
	}
	return c.Next(ctx, module)
}