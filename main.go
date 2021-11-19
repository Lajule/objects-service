package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/afero"
	"go.uber.org/zap"
)

var (
	// Flags
	port    int
	memory  bool
	rootDir string

	logger *zap.Logger

	afs *afero.Afero
)

func init() {
	flag.IntVar(&port, "p", 8080, "HTTP port")
	flag.BoolVar(&memory, "m", false, "Store objects in memory ?")
	flag.StringVar(&rootDir, "d", "./data", "Object root directory")
}

func main() {
	flag.Parse()

	logger, _ = zap.NewProduction()
	defer logger.Sync()

	logger.Info("Starting service",
		zap.Int("port", port),
		zap.Bool("memory", memory),
		zap.String("rootDir", rootDir))

	afs = newFs()

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.PUT("/objects/:bucket/:objectID", createOrReplaceObject)
	router.GET("/objects/:bucket/:objectID", getObject)
	router.DELETE("/objects/:bucket/:objectID", deleteObject)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("listening...", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exiting")
}
