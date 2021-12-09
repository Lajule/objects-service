package domains

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/Lajule/objects-service/pkg/domains/objects"
	"github.com/Lajule/objects-service/pkg/service"
)

// New creates a new engine
func New(logger *zap.Logger) []*service.Route {
	logger.Info("Creating engine")

	domains := []*service.Route{
		&service.Route{
			Path:        "/objects/:bucket/:objectID",
			Method:      http.MethodPut,
			HandlerFunc: objects.CreateOrReplace,
		},
		&service.Route{
			Path:        "/objects/:bucket/:objectID",
			Method:      http.MethodGet,
			HandlerFunc: objects.Get,
		},
		&service.Route{
			Path:        "/objects/:bucket/:objectID",
			Method:      http.MethodDelete,
			HandlerFunc: objects.Delete,
		},
	}

	return domains
}
