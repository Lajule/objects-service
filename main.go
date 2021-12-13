package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"net"
	"os"

	"go.uber.org/zap"
)

var (
	// Version contains the program version.
	Version = "development"

	// BasePath is the directory that contains buckets
	BasePath = flag.String("base-path", "./data", "Directory that contains buckets")

	// MemMapFs allows to store objects in memory
	MemMapFs = flag.Bool("mem-map-fs", false, "Use memory backed filesystem")

	// BindAddr is the network interface
	BindAddr = flag.String("bind-addr", ":8080", "Use specified network interface")

	// CACert is the file that contains list of trusted SSL Certificate Authorities
	CACert = flag.String("ca-cert", "", "File that contains list of trusted SSL Certificate Authorities")

	// ClientCert is the file that contains X.509 certificate
	ClientCert = flag.String("client-cert", "", "File that contains X.509 certificate")

	// ClientKey is the file that contains X.509 key
	ClientKey = flag.String("client-key", "", "File that contains X.509 key")
)

func main() {
	flag.Parse()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	logger.Info("objects-service", zap.String("version", Version))

	tcpAddr, err := net.ResolveTCPAddr("tcp", *BasePath)
	if err != nil {
		logger.Fatal(err.Error())
	}

	var tlsConfig *tls.Config

	if len(*CACert) > 0 {
		pool := x509.NewCertPool()

		data, err := os.ReadFile(*CACert)
		if err != nil {
			logger.Fatal(err.Error())
		}

		if !pool.AppendCertsFromPEM(data) {
			logger.Fatal("failed to parse CA certificate")
		}

		tlsConfig = &tls.Config{
			RootCAs:    pool,
			ClientAuth: tls.RequireAndVerifyClientCert,
			MinVersion: tls.VersionTLS12,
		}
	}

	srv := InitializeService(*BasePath, *MemMapFs, tcpAddr, tlsConfig, logger)
	srv.Start(*ClientCert, *ClientKey)
}
