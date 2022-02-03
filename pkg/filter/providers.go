package filter

import (
	"github.com/google/wire"
)

var Providers = wire.NewSet(
	wire.Struct(new(GitRepository), "*"),
)
