package main

import (
	"flag"

	"go.uber.org/zap"
)

func main() {
	rootDir := flag.String("d", "./data", "Object root directory")
	memory := flag.Bool("m", false, "Store objects in memory ?")
	port := flag.Int("p", 8080, "HTTP port")
	flag.Parse()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	srv := InitializeService(*rootDir, *memory, *port, logger)
	srv.Start()
}
