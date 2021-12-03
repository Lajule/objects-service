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
  -addr string
    	TCP address (default ":8080")
  -cert string
    	File that contains X.509 certificate
  -d string
    	Object root directory (default "./data")
  -key string
    	File that contains X.509 key
  -m	Store objects in memory ?
```

[1]: https://github.com/gin-gonic/gin
[2]: https://github.com/spf13/afero
[3]: https://github.com/google/wire
