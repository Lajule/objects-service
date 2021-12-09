// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"crypto/tls"
	"github.com/Lajule/objects-service/pkg/domains"
	"github.com/Lajule/objects-service/pkg/service"
	"github.com/Lajule/objects-service/pkg/store"
	"go.uber.org/zap"
	"net"
)

// Injectors from wire.go:

// InitializeService initializes a new service
func InitializeService(basePath string, memMapFs bool, tcpAddr *net.TCPAddr, tlsConfig *tls.Config, logger *zap.Logger) *service.Service {
	v := domains.New(logger)
	storeStore := store.New(basePath, memMapFs, logger)
	serviceService := service.New(tcpAddr, tlsConfig, v, storeStore, logger)
	return serviceService
}
