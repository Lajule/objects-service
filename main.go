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
	"go.uber.org/zap"
)

var logger *zap.Logger

func main() {
	port := flag.Int("p", 8080, "HTTP port")
	memory := flag.Bool("m", false, "Store objects in memory ?")
	rootDir := flag.String("d", "./data", "Object root directory")
	flag.Parse()

	logger, _ = zap.NewProduction()
	defer logger.Sync()

	logger.Info("Starting service",
		zap.Int("port", *port),
		zap.Bool("memory", *memory),
		zap.String("rootDir", *rootDir))

	store := newStore(*memory, *rootDir)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(func(c *gin.Context) {
		c.Set("store", store)
	})

	router.PUT("/objects/:bucket/:objectID", createOrReplaceObject)
	router.GET("/objects/:bucket/:objectID", getObject)
	router.DELETE("/objects/:bucket/:objectID", deleteObject)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
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
