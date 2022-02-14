// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package controllers

import (
	client2 "github.com/manifestival/client-go-client"
	"github.com/tiagoangelozup/charles-alpha/internal/client"
	"github.com/tiagoangelozup/charles-alpha/internal/object"
	"github.com/tiagoangelozup/charles-alpha/internal/resources"
	"github.com/tiagoangelozup/charles-alpha/internal/runtime"
	"github.com/tiagoangelozup/charles-alpha/pkg/filter"
	"github.com/tiagoangelozup/charles-alpha/pkg/module"
	"github.com/tiagoangelozup/charles-alpha/pkg/transformer"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// Injectors from wire.go:

func createReconcilers(managerManager manager.Manager) ([]Reconciler, error) {
	clientClient := runtime.Client(managerManager)
	clientModule := client.NewModule(clientClient)
	status := module.NewStatus(clientModule)
	gitRepository := &filter.GitRepository{}
	filters := &module.Filters{
		GitRepository: gitRepository,
	}
	scheme := runtime.Scheme(managerManager)
	unstructuredConverter := &object.UnstructuredConverter{
		Scheme: scheme,
	}
	transformerGitRepository := transformer.NewGitRepository(unstructuredConverter)
	reference := &object.Reference{
		Scheme: scheme,
	}
	metadata := transformer.NewMetadata(reference)
	transformers := &module.Transformers{
		GitRepository: transformerGitRepository,
		Metadata:      metadata,
	}
	config := runtime.Config(managerManager)
	manifestivalClient, err := client2.NewClient(config)
	if err != nil {
		return nil, err
	}
	manifests := &resources.Manifests{
		Client: manifestivalClient,
	}
	desiredState := module.NewDesiredState(filters, transformers, manifests)
	clientGitRepository := client.NewGitRepository(clientClient)
	artifactDownload := module.NewArtifactDownload(clientGitRepository, clientModule)
	moduleReconciler := NewModuleReconciler(status, desiredState, artifactDownload, clientModule)
	v := reconcilers(moduleReconciler)
	return v, nil
}
