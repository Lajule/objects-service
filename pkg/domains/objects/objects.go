package objects

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/Lajule/objects-service/pkg/service"
)

// CreateOrReplace creates or replaces an object
func CreateOrReplace(c *gin.Context) {
	s := c.MustGet("service").(*service.Service)

	bucket := c.Param("bucket")
	objectID := c.Param("objectID")

	defer c.Request.Body.Close()
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		s.Logger.Error("Can not read request body", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	s.Logger.Info("Creating object",
		zap.String("bucket", bucket),
		zap.String("objectID", objectID),
		zap.ByteString("data", data))

	if err := s.Store.CreateBucketIfNotExists(bucket); err != nil {
		s.Logger.Error("Can not create bucket if not exists", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	object, err := s.Store.CreateOrOpenObject(bucket, objectID)
	if err != nil {
		s.Logger.Error("Can not create or open object", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	defer object.Close()
	if _, err := object.WriteString(string(data)); err != nil {
		s.Logger.Error("Can not write object", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, struct {
		ID string `json:"id"`
	}{
		ID: objectID,
	})
}

// Get get an object
func Get(c *gin.Context) {
	s := c.MustGet("service").(*service.Service)

	bucket := c.Param("bucket")
	objectID := c.Param("objectID")

	s.Logger.Info("Getting object",
		zap.String("bucket", bucket),
		zap.String("objectID", objectID))

	object, err := s.Store.GetObjectIfExists(bucket, objectID)
	if err != nil {
		s.Logger.Error("Can not get object if exists", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	if object == nil {
		s.Logger.Info("Object not exists")
		c.Status(http.StatusNotFound)
		return
	}

	defer object.Close()
	data, err := io.ReadAll(object)
	if err != nil {
		s.Logger.Error("Can not read object", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	c.String(http.StatusOK, string(data))
}

// Delete deletes an object
func Delete(c *gin.Context) {
	s := c.MustGet("service").(*service.Service)

	bucket := c.Param("bucket")
	objectID := c.Param("objectID")

	s.Logger.Info("Deleting object",
		zap.String("bucket", bucket),
		zap.String("objectID", objectID))

	removed, err := s.Store.RemoveObjectIfExists(bucket, objectID)
	if err != nil {
		s.Logger.Error("Can not remove object if exists", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	if !removed {
		s.Logger.Info("Object not exists")
		c.Status(http.StatusBadRequest)
		return
	}

	c.Status(http.StatusOK)
}
