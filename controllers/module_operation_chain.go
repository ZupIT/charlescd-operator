package controllers

import (
	"context"

	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"

	deployv1alpha1 "github.com/tiagoangelozup/charles-alpha/api/v1alpha1"
	"github.com/tiagoangelozup/charles-alpha/internal/runtime"
	"github.com/tiagoangelozup/charles-alpha/internal/tracing"
)

type ModuleOperation func(ctx context.Context, module *deployv1alpha1.Module) (ctrl.Result, error)

type ModuleOperationChain struct{ operations []ModuleOperation }

func NewModuleOperationChain(operations ...ModuleOperation) *ModuleOperationChain {
	return &ModuleOperationChain{operations: operations}
}

func (m *ModuleOperationChain) Handle(ctx context.Context, module *deployv1alpha1.Module) (ctrl.Result, error) {
	span := tracing.SpanFromContext(ctx)
	l := func() logr.Logger {
		if span != nil {
			return logger.WithValues("trace", span)
		}
		return logger
	}()

	for _, op := range m.operations {
		result, err := op(ctx, module)
		if err != nil || result.RequeueAfter > 0 || result.Requeue {
			if span != nil {
				span.SetError(err)
			}
			return result, err
		}
	}

	l.Info("Successfully reconciled!")
	return runtime.Finish()
}
