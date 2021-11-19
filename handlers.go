package main

import (
	"io"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func createOrReplaceObject(c *gin.Context) {
	bucket := c.Param("bucket")
	objectID := c.Param("objectID")

	defer c.Request.Body.Close()
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.Error("Can not read request body", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	logger.Info("Creating object",
		zap.String("bucket", bucket),
		zap.String("objectID", objectID),
		zap.ByteString("data", data))

	if err := createBucketIfNotExists(filepath.Join(rootDir, bucket)); err != nil {
		logger.Error("Can not create bucket if not exists", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	object, err := createOrOpenObject(filepath.Join(rootDir, bucket, objectID))
	if err != nil {
		logger.Error("Can not create or open object", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	defer object.Close()
	if _, err := object.WriteString(string(data)); err != nil {
		logger.Error("Can not write object", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, struct {
		Id string
	}{
		Id: objectID,
	})
}

func getObject(c *gin.Context) {
	bucket := c.Param("bucket")
	objectID := c.Param("objectID")

	logger.Info("Getting object",
		zap.String("bucket", bucket),
		zap.String("objectID", objectID))

	object, err := getObjectIfExists(filepath.Join(rootDir, bucket, objectID))
	if err != nil {
		logger.Error("Can not get object if exists", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	if object == nil {
		logger.Info("Object not exists")
		c.Status(http.StatusBadRequest)
		return
	}

	defer object.Close()
	data, err := io.ReadAll(object)
	if err != nil {
		logger.Error("Can not read object", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	c.String(http.StatusOK, string(data))
}

func deleteObject(c *gin.Context) {
	bucket := c.Param("bucket")
	objectID := c.Param("objectID")

	logger.Info("Deleting object",
		zap.String("bucket", bucket),
		zap.String("objectID", objectID))

	removed, err := removeObjectIfExists(filepath.Join(rootDir, bucket, objectID))
	if err != nil {
		logger.Error("Can not remove object if exists", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	if !removed {
		logger.Info("Object not exists")
		c.Status(http.StatusBadRequest)
		return
	}

	c.Status(http.StatusOK)
}
