// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package controllers

import (
	"github.com/manifestival/client-go-client"
	"github.com/tiagoangelozup/charles-alpha/internal/manifests"
	"github.com/tiagoangelozup/charles-alpha/internal/module"
	"github.com/tiagoangelozup/charles-alpha/internal/object"
	"github.com/tiagoangelozup/charles-alpha/internal/runtime"
	"github.com/tiagoangelozup/charles-alpha/internal/usecase"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// Injectors from wire.go:

func createReconcilers(managerManager manager.Manager) ([]Reconciler, error) {
	config := runtime.Config(managerManager)
	manifestivalClient, err := client.NewClient(config)
	if err != nil {
		return nil, err
	}
	service := &manifests.Service{
		Client: manifestivalClient,
	}
	scheme := runtime.Scheme(managerManager)
	unstructuredConverter := &object.UnstructuredConverter{
		Scheme: scheme,
	}
	reference := &object.Reference{
		Scheme: scheme,
	}
	desiredState := &usecase.DesiredState{
		Manifests: service,
		Object:    unstructuredConverter,
		Reference: reference,
	}
	adapter := ModuleAdapter{
		DesiredState: desiredState,
	}
	clientClient := runtime.Client(managerManager)
	moduleService := &module.Service{
		Client: clientClient,
	}
	moduleReconciler := &ModuleReconciler{
		ModuleAdapter: adapter,
		ModuleGetter:  moduleService,
	}
	v := reconcilers(moduleReconciler)
	return v, nil
}
