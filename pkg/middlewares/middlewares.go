package middlewares

import (
	"github.com/google/wire"
)

// Set is used by wire
var Set = wire.NewSet(NewLogger)
