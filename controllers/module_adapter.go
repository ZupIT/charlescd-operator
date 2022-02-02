package controllers

import "github.com/tiagoangelozup/charles-alpha/internal/usecase"

type ModuleAdapter struct {
	DesiredState     *usecase.DesiredState
	HelmInstallation *usecase.HelmInstallation
}
