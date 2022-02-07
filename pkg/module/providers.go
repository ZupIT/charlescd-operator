package module

import "github.com/google/wire"

var Providers = wire.NewSet(
	NewDesiredState,
	NewHelmInstallation,
	NewStatus,
	wire.Struct(new(Filters), "*"),
	wire.Struct(new(Transformers), "*"),
)
