package transformer

import (
	"github.com/google/wire"
)

var Providers = wire.NewSet(
	NewGitRepository,
	NewMetadata,
)
