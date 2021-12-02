//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"go.uber.org/zap"

	"github.com/Lajule/objects-service/pkg/service"
	"github.com/Lajule/objects-service/pkg/store"
)

func InitializeService(rootDir string, memory bool, port int, logger *zap.Logger) *service.Service {
	wire.Build(store.NewStore, service.NewService)
	return &service.Service{}
}
