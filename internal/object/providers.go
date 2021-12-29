package object

import "github.com/google/wire"

var Providers = wire.NewSet(
	wire.Struct(new(UnstructuredConverter), "*"),
	wire.Struct(new(Reference), "*"),
)
