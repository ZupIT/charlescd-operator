package client

import (
	"github.com/google/wire"
)

var Providers = wire.NewSet(
	NewGitRepository,
	NewModule,
)
