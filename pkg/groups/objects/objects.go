package objects

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/Lajule/objects-service/pkg/middlewares"
	"github.com/Lajule/objects-service/pkg/service"
	"github.com/Lajule/objects-service/pkg/store"
)

// Group is an alias to service.Group
type Group service.Group

// Params contains group's params
type Params struct {
	// Logger gives access to logger
	Logger *zap.Logger

	// Store gives access to store
	Store *store.Store
}

// NewObjects creates objects group
func NewObjects(requestLogger middlewares.Logger, st *store.Store, logger *zap.Logger) *Group {
	return &Group{
		Name:        "/objects",
		Middlewares: []gin.HandlerFunc{gin.HandlerFunc(requestLogger)},
		Routes: []*service.Route{
			&service.Route{
				Path:         "/:bucket/:objectID",
				Method:       http.MethodPut,
				HandlerFuncs: []gin.HandlerFunc{createOrReplace},
			},
			&service.Route{
				Path:         "/:bucket/:objectID",
				Method:       http.MethodGet,
				HandlerFuncs: []gin.HandlerFunc{get},
			},
			&service.Route{
				Path:         "/:bucket/:objectID",
				Method:       http.MethodDelete,
				HandlerFuncs: []gin.HandlerFunc{deleteObject},
			},
		},
		Params: &Params{
			Logger: logger.Named("objects"),
			Store:  st,
		},
	}
}

func createOrReplace(c *gin.Context) {
	params := c.MustGet("params").(*Params)

	bucket := c.Param("bucket")
	objectID := c.Param("objectID")

	defer c.Request.Body.Close()
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		params.Logger.Error("can not read request body", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	params.Logger.Info("creating object",
		zap.String("bucket", bucket),
		zap.String("objectID", objectID),
		zap.ByteString("data", data))

	if err := params.Store.CreateBucketIfNotExists(bucket); err != nil {
		params.Logger.Error("can not create bucket if not exists", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	object, err := params.Store.CreateOrOpenObject(bucket, objectID)
	if err != nil {
		params.Logger.Error("can not create or open object", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	defer object.Close()
	if _, err := object.WriteString(string(data)); err != nil {
		params.Logger.Error("can not write object", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, struct {
		ID string `json:"id"`
	}{
		ID: objectID,
	})
}

func get(c *gin.Context) {
	params := c.MustGet("params").(*Params)

	bucket := c.Param("bucket")
	objectID := c.Param("objectID")

	params.Logger.Info("getting object",
		zap.String("bucket", bucket),
		zap.String("objectID", objectID))

	object, err := params.Store.GetObjectIfExists(bucket, objectID)
	if err != nil {
		params.Logger.Error("can not get object if exists", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	if object == nil {
		params.Logger.Info("object not exists")
		c.Status(http.StatusNotFound)
		return
	}

	defer object.Close()
	data, err := io.ReadAll(object)
	if err != nil {
		params.Logger.Error("can not read object", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	c.String(http.StatusOK, string(data))
}

func deleteObject(c *gin.Context) {
	params := c.MustGet("params").(*Params)

	bucket := c.Param("bucket")
	objectID := c.Param("objectID")

	params.Logger.Info("deleting object",
		zap.String("bucket", bucket),
		zap.String("objectID", objectID))

	removed, err := params.Store.RemoveObjectIfExists(bucket, objectID)
	if err != nil {
		params.Logger.Error("can not remove object if exists", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	if !removed {
		params.Logger.Info("object not exists")
		c.Status(http.StatusNotFound)
		return
	}

	c.Status(http.StatusOK)
}
