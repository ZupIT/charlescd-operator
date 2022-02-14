package module

import "github.com/google/wire"

var Providers = wire.NewSet(
	NewArtifactDownload,
	NewDesiredState,
	NewHelmValidation,
	NewStatus,
	wire.Struct(new(Filters), "*"),
	wire.Struct(new(Transformers), "*"),
)
