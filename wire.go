//go:build wireinject
// +build wireinject

package main

import (
	"crypto/tls"
	"net"

	"github.com/google/wire"
	"go.uber.org/zap"

	"github.com/Lajule/objects-service/pkg/service"
	"github.com/Lajule/objects-service/pkg/store"
)

func InitializeService(basePath string, memory bool, tcpAddr *net.TCPAddr, tlsConfig *tls.Config, logger *zap.Logger) *service.Service {
	wire.Build(store.NewStore, service.NewService)
	return &service.Service{}
}
