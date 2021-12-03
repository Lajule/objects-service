// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"github.com/Lajule/objects-service/pkg/service"
	"github.com/Lajule/objects-service/pkg/store"
	"go.uber.org/zap"
)

// Injectors from wire.go:

func InitializeService(rootDir string, memory bool, port int, logger *zap.Logger) *service.Service {
	storeStore := store.NewStore(rootDir, memory, logger)
	serviceService := service.NewService(port, storeStore, logger)
	return serviceService
}
