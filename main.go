package main

import (
	"crypto/tls"
	"crypto/x509"
	"net"
	"os"

	"github.com/alexflint/go-arg"
	"go.uber.org/zap"
)

// Args contains command line arguments
type Args struct {
	BasePath   string `arg:"--base-path" default:"./data" help:"Directory that contains buckets"`
	MemMapFs   bool   `arg:"--mem-map-fs" help:"Use memory backed filesystem"`
	BindAddr   string `arg:"--bind-addr" default:":8080" help:"Use specified network interface"`
	CACert     string `arg:"--ca-cert" help:"File that contains list of trusted SSL Certificate Authorities"`
	ClientCert string `arg:"--client-cert" help:"File that contains X.509 certificate"`
	ClientKey  string `arg:"--client-key" help:"File that contains X.509 key"`
}

var (
	// Version contains the program version.
	Version = "development"
)

func main() {
	args := Args{}
	arg.MustParse(&args)

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	logger.Info("objects-service", zap.String("version", Version))

	tcpAddr, err := net.ResolveTCPAddr("tcp", args.BindAddr)
	if err != nil {
		logger.Fatal(err.Error())
	}

	var tlsConfig *tls.Config

	if len(args.CACert) > 0 {
		pool := x509.NewCertPool()

		data, err := os.ReadFile(args.CACert)
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

	srv := InitializeService(args.BasePath, args.MemMapFs, tcpAddr, tlsConfig, logger)
	srv.Start(args.ClientCert, args.ClientKey)
}
