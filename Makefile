ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
MAKEFLAGS += --always-make
GOARCH ?= $(shell go env GOARCH)
CXX := g++

build:
	go build -o dorisdump main.go

build-hyper:
	CXX=$(CXX) go build -tags chimera -o dorisdump main.go

build-darwin:
	CXX=$(CXX) GOOS=darwin CGO_ENABLED=1 GOEXPERIMENT=newinliner CGO_LDFLAGS='-static-libstdc++' go build -tags chimera -o dorisdump-darwin-$(GOARCH) -a -trimpath -ldflags "-w -s"
	tar czf dorisdump-darwin-$(GOARCH).tar.gz dorisdump-darwin-$(GOARCH)

build-linux:
	GOOS=linux CGO_ENABLED=1 GOEXPERIMENT=newinliner go build -tags chimera -o dorisdump-linux-$(GOARCH) -a -trimpath -installsuffix cgo -ldflags "-w -s -linkmode external -extldflags -static"
	tar czf dorisdump-linux-$(GOARCH).tar.gz dorisdump-linux-$(GOARCH)

run:
	@go run main.go $(ARGS)

test:
	@go test -v ./...

install: build
	cp dorisdump /usr/local/bin

gen:
	@go generate ./...

fmt:
	@go fmt .
	@goimports -l -w -local "github.com/Thearas/dorisdump" .

addcmd:
	cobra-cli --license apache --author "Thearas thearas850@gmail.com" add $(ARGS)

%:
	@:
