# objects-service

An HTTP service that is able to store objects organized by buckets. This service is built with [gin][1] web framework and with [afero][2] filesystem abstraction.

## Build

Build the service with:

```sh
go build
```

> The service use [Wire][3] automated initialization.

## Usage

```
./objects-service -h
Usage of ./objects-service:
  -base-path string
        Store base path (default "./data")
  -bind-addr string
        Use specified network interface (default ":8080")
  -ca-cert string
        File that contains list of trusted SSL Certificate Authorities
  -client-cert string
        File that contains X.509 certificate
  -client-key string
        File that contains X.509 key
  -mem-map-fs
        Use memory backed filesystem
```

[1]: https://github.com/gin-gonic/gin
[2]: https://github.com/spf13/afero
[3]: https://github.com/google/wire
