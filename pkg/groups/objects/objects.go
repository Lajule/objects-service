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

// NewObjects creates objects group
func NewObjects(requestLogger middlewares.Logger, logger *zap.Logger) *Group {
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
		Logger: logger.Named("objects"),
	}
}

func createOrReplace(c *gin.Context) {
	logger := c.MustGet("logger").(*zap.Logger)
	st := c.MustGet("store").(*store.Store)

	bucket := c.Param("bucket")
	objectID := c.Param("objectID")

	defer c.Request.Body.Close()
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.Error("can not read request body", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	logger.Info("creating object",
		zap.String("bucket", bucket),
		zap.String("objectID", objectID),
		zap.ByteString("data", data))

	if err := st.CreateBucketIfNotExists(bucket); err != nil {
		logger.Error("can not create bucket if not exists", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	object, err := st.CreateOrOpenObject(bucket, objectID)
	if err != nil {
		logger.Error("can not create or open object", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	defer object.Close()
	if _, err := object.WriteString(string(data)); err != nil {
		logger.Error("can not write object", zap.Error(err))
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
	logger := c.MustGet("logger").(*zap.Logger)
	st := c.MustGet("store").(*store.Store)

	bucket := c.Param("bucket")
	objectID := c.Param("objectID")

	logger.Info("getting object",
		zap.String("bucket", bucket),
		zap.String("objectID", objectID))

	object, err := st.GetObjectIfExists(bucket, objectID)
	if err != nil {
		logger.Error("can not get object if exists", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	if object == nil {
		logger.Info("object not exists")
		c.Status(http.StatusNotFound)
		return
	}

	defer object.Close()
	data, err := io.ReadAll(object)
	if err != nil {
		logger.Error("can not read object", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	c.String(http.StatusOK, string(data))
}

func deleteObject(c *gin.Context) {
	logger := c.MustGet("logger").(*zap.Logger)
	st := c.MustGet("store").(*store.Store)

	bucket := c.Param("bucket")
	objectID := c.Param("objectID")

	logger.Info("deleting object",
		zap.String("bucket", bucket),
		zap.String("objectID", objectID))

	removed, err := st.RemoveObjectIfExists(bucket, objectID)
	if err != nil {
		logger.Error("can not remove object if exists", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	if !removed {
		logger.Info("object not exists")
		c.Status(http.StatusNotFound)
		return
	}

	c.Status(http.StatusOK)
}
