# objects-service

An HTTP service that is able to store objects organized by buckets. This service is built with [gin][1] web framework and with [afero][2] filesystem abstraction.

## Build

Build the service with:

```sh
go build -ldflags="-s -X 'main.Version=0.1.0'"
```

> The service use [Wire][3] automated initialization.

## Usage

```
./objects-service -h
Usage: objects-service [--base-path BASE-PATH] [--mem-map-fs] [--bind-addr BIND-ADDR] [--ca-cert CA-CERT] [--client-cert CLIENT-CERT] [--client-key CLIENT-KEY]

Options:
  --base-path BASE-PATH
                         Directory that contains buckets [default: ./data]
  --mem-map-fs           Use memory backed filesystem
  --bind-addr BIND-ADDR
                         Use specified network interface [default: :8080]
  --ca-cert CA-CERT      File that contains list of trusted SSL Certificate Authorities
  --client-cert CLIENT-CERT
                         File that contains X.509 certificate
  --client-key CLIENT-KEY
                         File that contains X.509 key
  --help, -h             display this help and exit
```

[1]: https://github.com/gin-gonic/gin
[2]: https://github.com/spf13/afero
[3]: https://github.com/google/wire
