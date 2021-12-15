//go:build wireinject
// +build wireinject

package main

import (
	"crypto/tls"
	"net"

	"github.com/google/wire"
	"go.uber.org/zap"

	"github.com/Lajule/objects-service/pkg/groups"
	"github.com/Lajule/objects-service/pkg/service"
	"github.com/Lajule/objects-service/pkg/store"
)

// InitializeService initializes a new service
func InitializeService(basePath string, memMapFs bool, tcpAddr *net.TCPAddr, tlsConfig *tls.Config, logger *zap.Logger) *service.Service {
	wire.Build(store.New, groups.Set, groups.New, service.New)
	return &service.Service{}
}
