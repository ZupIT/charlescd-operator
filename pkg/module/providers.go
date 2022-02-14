package module

import "github.com/google/wire"

var Providers = wire.NewSet(
	NewDesiredState,
	NewArtifactDownload,
	NewStatus,
	wire.Struct(new(Filters), "*"),
	wire.Struct(new(Transformers), "*"),
)
