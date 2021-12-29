package predicate

import (
	"github.com/google/wire"
)

var Providers = wire.NewSet(
	wire.Struct(new(RepoStatus), "*"),
	wire.Struct(new(Module), "*"),
)
