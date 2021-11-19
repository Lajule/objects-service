# objects-service

HTTP service that is able to store objects organized by buckets. This service is built with [gin][1] web framework and with [afero][2] filesystem abstraction.

## Build

Build the service with:

```sh
go build
```

## Usage

```
./objects-service -h
Usage of ./objects-service:
  -d string
    	Object root directory (default "./data")
  -m	Store objects in memory ?
  -p int
    	HTTP port (default 8080)
```

[1]: https://github.com/gin-gonic/gin
[2]: https://github.com/spf13/afero