package main

import (
	"flag"

	"go.uber.org/zap"

	"github.com/Lajule/objects-service/pkg/service"
	"github.com/Lajule/objects-service/pkg/store"
)

var logger *zap.Logger

func main() {
	port := flag.Int("p", 8080, "HTTP port")
	memory := flag.Bool("m", false, "Store objects in memory ?")
	rootDir := flag.String("d", "./data", "Object root directory")
	flag.Parse()

	logger, _ = zap.NewProduction()
	defer logger.Sync()

	logger.Info("Starting ...",
		zap.Int("port", *port),
		zap.Bool("memory", *memory),
		zap.String("rootDir", *rootDir))

	st := store.NewStore(logger, *memory, *rootDir)
	srv := service.NewService(logger, st, *port)

	srv.Start()
}
