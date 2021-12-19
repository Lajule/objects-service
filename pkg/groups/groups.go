package groups

import (
	"github.com/google/wire"

	"github.com/Lajule/objects-service/pkg/groups/objects"
	"github.com/Lajule/objects-service/pkg/service"
)

// Set is used by wire
var Set = wire.NewSet(NewGroups, objects.NewObjects)

// NewGroups creates all groups
func NewGroups(objectsGroup *objects.Group) []*service.Group {
	return []*service.Group{
		(*service.Group)(objectsGroup),
	}
}
