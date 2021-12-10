package groups

import (
	"go.uber.org/zap"

	"github.com/Lajule/objects-service/pkg/groups/objects"
	"github.com/Lajule/objects-service/pkg/service"
)

// New creates all groups
func New(logger *zap.Logger) []*service.Group {
	logger.Info("Creating groups")

	return []*service.Group{
		objects.Group,
	}
}
