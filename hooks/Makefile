GO=go
HOOKS=cmd/configure/configure cmd/install/install

build: $(HOOKS)

cmd/configure/configure:
	$(GO) build -o $@ ./cmd/configure

cmd/install/install:
	$(GO) build -o $@ ./cmd/install

# TODO: fixup unit testing
test:
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) vet ./...
	gofmt -l .
	[ "`gofmt -l .`" = "" ]
#	./bin/test-go-mod-tidy.sh
#	./bin/test-attribution-txt.sh

clean:
	rm -f $(HOOKS)

