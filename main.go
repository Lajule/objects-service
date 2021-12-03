package main

import (
	"flag"
	"net"
	"crypto/tls"

	"go.uber.org/zap"
)

// Version contains the program version.
var Version = "development"

func main() {
	rootDir := flag.String("d", "./data", "Object root directory")
	memory := flag.Bool("m", false, "Store objects in memory ?")
	addr := flag.String("addr", ":8080", "TCP address")
	certFile := flag.String("cert", "", "File that contains X.509 certificate")
	keyFile := flag.String("key", "", "File that contains X.509 key")
	flag.Parse()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	logger.Info("objects-service", zap.String("version", Version))

	tcpAddr, err := net.ResolveTCPAddr("tcp", *addr)
	if err != nil {
		logger.Fatal(err.Error())
	}

	var tlsConfig *tls.Config
	if *certFile != "" && *keyFile != "" {
		tlsConfig = &tls.Config{}
		tlsConfig.Certificates = make([]tls.Certificate, 1)
		tlsConfig.Certificates[0], err = tls.LoadX509KeyPair(*certFile, *keyFile)
		if err != nil {
			logger.Fatal(err.Error())
		}
	}

	srv := InitializeService(*rootDir, *memory, tcpAddr, tlsConfig, logger)
	srv.Start()
}
