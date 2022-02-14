package resources

import (
	"github.com/google/wire"
	mfc "github.com/manifestival/client-go-client"
)

var Providers = wire.NewSet(
	mfc.NewClient,
	wire.Struct(new(Manifests), "*"),
)
