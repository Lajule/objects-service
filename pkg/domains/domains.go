package domains

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/Lajule/objects-service/pkg/domains/objects"
	"github.com/Lajule/objects-service/pkg/service"
)

// New creates a new engine
func New(logger *zap.Logger) []*service.Route {
	logger.Info("Creating domains")

	domains := []*service.Route{
		&service.Route{
			Path:         "/objects/:bucket/:objectID",
			Method:       http.MethodPut,
			HandlerFuncs: []gin.HandlerFunc{objects.CreateOrReplace},
		},
		&service.Route{
			Path:         "/objects/:bucket/:objectID",
			Method:       http.MethodGet,
			HandlerFuncs: []gin.HandlerFunc{objects.Get},
		},
		&service.Route{
			Path:         "/objects/:bucket/:objectID",
			Method:       http.MethodDelete,
			HandlerFuncs: []gin.HandlerFunc{objects.Delete},
		},
	}

	return domains
}
