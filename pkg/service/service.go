package service

import (
	"context"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/Lajule/objects-service/pkg/store"
)

// Service contains an HTTP server and a store
type Service struct {
	srv   *http.Server
	store *store.Store
	log   *zap.Logger
}

// NewService creates a new service
func NewService(tcpAddr *net.TCPAddr, tlsConfig *tls.Config, st *store.Store, logger *zap.Logger) *Service {
	logger.Info("Creating service",
		zap.String("address", tcpAddr.String()))

	service := &Service{
		store: st,
		log:   logger.Named("service"),
	}

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	engine.Use(func(c *gin.Context) {
		c.Set("service", service)
	})

	engine.PUT("/objects/:bucket/:objectID", createOrReplaceObject)
	engine.GET("/objects/:bucket/:objectID", getObject)
	engine.DELETE("/objects/:bucket/:objectID", deleteObject)

	service.srv = &http.Server{
		Addr:      tcpAddr.String(),
		Handler:   engine,
		TLSConfig: tlsConfig,
	}

	return service
}

func (s *Service) Start() {
	s.log.Info("Service starting")

	go func() {
		var err error
		if s.srv.TLSConfig != nil {
			err = s.srv.ListenAndServeTLS("", "")
		} else {
			err = s.srv.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			s.log.Fatal("Service not listening", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	s.log.Info("Shutting down service...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.srv.Shutdown(ctx); err != nil {
		s.log.Fatal("Service forced to shutdown", zap.Error(err))
	}

	s.log.Info("Service exiting")
}

func createOrReplaceObject(c *gin.Context) {
	s := c.MustGet("service").(*Service)

	bucket := c.Param("bucket")
	objectID := c.Param("objectID")

	defer c.Request.Body.Close()
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		s.log.Error("Can not read request body", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	s.log.Info("Creating object",
		zap.String("bucket", bucket),
		zap.String("objectID", objectID),
		zap.ByteString("data", data))

	if err := s.store.CreateBucketIfNotExists(bucket); err != nil {
		s.log.Error("Can not create bucket if not exists", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	object, err := s.store.CreateOrOpenObject(bucket, objectID)
	if err != nil {
		s.log.Error("Can not create or open object", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	defer object.Close()
	if _, err := object.WriteString(string(data)); err != nil {
		s.log.Error("Can not write object", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, struct {
		ID string `json:"id"`
	}{
		ID: objectID,
	})
}

func getObject(c *gin.Context) {
	s := c.MustGet("service").(*Service)

	bucket := c.Param("bucket")
	objectID := c.Param("objectID")

	s.log.Info("Getting object",
		zap.String("bucket", bucket),
		zap.String("objectID", objectID))

	object, err := s.store.GetObjectIfExists(bucket, objectID)
	if err != nil {
		s.log.Error("Can not get object if exists", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	if object == nil {
		s.log.Info("Object not exists")
		c.Status(http.StatusNotFound)
		return
	}

	defer object.Close()
	data, err := io.ReadAll(object)
	if err != nil {
		s.log.Error("Can not read object", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	c.String(http.StatusOK, string(data))
}

func deleteObject(c *gin.Context) {
	s := c.MustGet("service").(*Service)

	bucket := c.Param("bucket")
	objectID := c.Param("objectID")

	s.log.Info("Deleting object",
		zap.String("bucket", bucket),
		zap.String("objectID", objectID))

	removed, err := s.store.RemoveObjectIfExists(bucket, objectID)
	if err != nil {
		s.log.Error("Can not remove object if exists", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	if !removed {
		s.log.Info("Object not exists")
		c.Status(http.StatusBadRequest)
		return
	}

	c.Status(http.StatusOK)
}
