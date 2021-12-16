package groups

import (
	"github.com/google/wire"
	"go.uber.org/zap"

	"github.com/Lajule/objects-service/pkg/groups/objects"
	"github.com/Lajule/objects-service/pkg/service"
)

// Set is used by wire
var Set = wire.NewSet(NewGroups, objects.NewObjects)

// NewGroups creates all groups
func NewGroups(objectsGroup *objects.Group, logger *zap.Logger) []*service.Group {
	logger.Info("Creating groups")

	return []*service.Group{
		(*service.Group)(objectsGroup),
	}
}
