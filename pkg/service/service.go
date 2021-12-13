package service

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/Lajule/objects-service/pkg/store"
)

// Group is router group
type Group struct {
	// Name is the name of router group
	Name string

	// Routes is the list of routes
	Routes []*Route
}

// Route contains a HTTP method, a path and a handler function
type Route struct {
	// Path is the API endpoint
	Path string

	// Method is the HTTP method of the route
	Method string

	// HandlerFuncs is the function that handles the request
	HandlerFuncs []gin.HandlerFunc
}

// Service contains an HTTP server and a store
type Service struct {
	// Store gives access to object store
	Store *store.Store

	// Logger gives access to logger
	Logger *zap.Logger

	srv *http.Server
}

// New creates a new service
func New(tcpAddr *net.TCPAddr, tlsConfig *tls.Config, groups []*Group, st *store.Store, logger *zap.Logger) *Service {
	logger.Info("Creating service",
		zap.String("address", tcpAddr.String()))

	service := &Service{
		Store:  st,
		Logger: logger.Named("service"),
	}

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	engine.Use(func(c *gin.Context) {
		c.Set("service", service)
	})

	for _, group := range groups {
		value := reflect.ValueOf(engine.Group(group.Name))
		for _, route := range group.Routes {
			in := []reflect.Value{
				reflect.ValueOf(route.Path),
			}

			for _, handlerFunc := range route.HandlerFuncs {
				in = append(in, reflect.ValueOf(handlerFunc))
			}

			f := value.MethodByName(route.Method)
			f.Call(in)
		}
	}

	service.srv = &http.Server{
		Addr:      tcpAddr.String(),
		Handler:   engine,
		TLSConfig: tlsConfig,
	}

	return service
}

// Start starts HTTP service
func (s *Service) Start(clientCert, clientKey string) {
	s.Logger.Info("Service starting")

	go func() {
		var err error
		if len(clientCert) > 0 && len(clientKey) > 0 {
			err = s.srv.ListenAndServeTLS(clientCert, clientKey)
		} else {
			err = s.srv.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			s.Logger.Fatal("Service not listening", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	s.Logger.Info("Shutting down service...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.srv.Shutdown(ctx); err != nil {
		s.Logger.Fatal("Service forced to shutdown", zap.Error(err))
	}

	s.Logger.Info("Service exiting")
}
