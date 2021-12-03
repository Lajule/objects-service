BINNAME := objects-service
PACKAGE := github.com/Lajule/objects-service
VERSION := 0.0.1
MAINDIR := .
TARGETS := all run debug watch generate tidy test vet lint format clean bootstrap dist

all:
	go build -ldflags="-s -X 'main.Version=$(VERSION)'" -tags "$(GOTAGS)" -o $(BINNAME) $(MAINDIR)

run:
	go run -tags "$(GOTAGS)" $(MAINDIR)

debug:
	dlv debug --build-flags "-tags '$(GOTAGS)'" $(PACKAGE)

watch:
	air -c .air.toml

generate:
	go generate ./...

tidy:
	go mod tidy

test:
	go test -tags "$(GOTAGS)" -v ./...

vet:
	go vet ./...

lint:
	golint ./...

format:
	go fmt ./...

clean:
	go clean -i -cache -testcache -modcache

bootstrap:
	find . -mindepth 1 -type d -exec sh -c "echo \"TARGETS := $(TARGETS)\n\n\\\$$(TARGETS):\n\t\\\$$(MAKE) -C .. \\\$$@\n.PHONY: \\\$$(TARGETS)\" >{}/Makefile" \;

dist:
	touch $(BINNAME).tar.gz && tar -czf $(BINNAME).tar.gz --exclude=$(BINNAME).tar.gz .

.PHONY: $(TARGETS)
