package usecase

import "github.com/google/wire"

var Providers = wire.NewSet(
	wire.Struct(new(DesiredState), "*"),
	wire.Struct(new(HelmInstallation), "*"),
)
