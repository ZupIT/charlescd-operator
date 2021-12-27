// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package controllers

import (
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Injectors from wire.go:

func createReconcilers(clientClient client.Client, scheme *runtime.Scheme) ([]Reconciler, error) {
	moduleReconciler := &ModuleReconciler{}
	v := reconcilers(moduleReconciler)
	return v, nil
}
