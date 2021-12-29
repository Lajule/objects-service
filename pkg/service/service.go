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
)

// Group is a router group
type Group struct {
	// Name is the name of router group
	Name string

	// Middlewares contains the router group middlewares
	Middlewares []gin.HandlerFunc

	// Routes is the list of routes
	Routes []*Route

	// SubGroups contains sub groups
	SubGroups []*Group

	// Params gives access to params
	Params interface{}
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
	srv    *http.Server
	logger *zap.Logger
}

// NewService creates a new service
func NewService(tcpAddr *net.TCPAddr, tlsConfig *tls.Config, groups []*Group, logger *zap.Logger) *Service {
	logger.Info("creating service",
		zap.String("address", tcpAddr.String()))

	service := &Service{
		logger: logger.Named("service"),
	}

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	walkthroughGroups(engine.Group("/"), groups)

	service.srv = &http.Server{
		Addr:      tcpAddr.String(),
		Handler:   engine,
		TLSConfig: tlsConfig,
	}

	return service
}

// Start starts HTTP service
func (s *Service) Start(clientCert, clientKey string) {
	s.logger.Info("service starting")

	go func() {
		var err error
		if len(clientCert) > 0 && len(clientKey) > 0 {
			err = s.srv.ListenAndServeTLS(clientCert, clientKey)
		} else {
			err = s.srv.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("service not listening", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	s.logger.Info("shutting down service...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.srv.Shutdown(ctx); err != nil {
		s.logger.Fatal("service forced to shutdown", zap.Error(err))
	}

	s.logger.Info("service exiting")
}

func walkthroughGroups(root *gin.RouterGroup, groups []*Group) {
	for _, group := range groups {
		routerGroup := root.Group(group.Name)

		routerGroup.Use(func(c *gin.Context) {
			c.Set("params", group.Params)
		})

		if group.Middlewares != nil {
			for _, middleware := range group.Middlewares {
				routerGroup.Use(middleware)
			}
		}

		value := reflect.ValueOf(routerGroup)
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

		if group.SubGroups != nil {
			walkthroughGroups(routerGroup, group.SubGroups)
		}
	}
}
