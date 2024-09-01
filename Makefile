ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
MAKEFLAGS += --always-make
GOARCH ?= $(shell go env GOARCH)

build:
	go build -o dorisdump main.go

build-darwin:
	GOOS=darwin CGO_ENABLED=1 go build -o dorisdump-darwin-$(GOARCH) -a -trimpath
	tar czf dorisdump-darwin-$(GOARCH).tar.gz dorisdump-darwin-$(GOARCH)

build-linux:
	GOOS=linux CGO_ENABLED=1 go build -o dorisdump-linux-$(GOARCH) -a -trimpath -installsuffix cgo -ldflags "-linkmode external -extldflags -static"
	tar czf dorisdump-linux-$(GOARCH).tar.gz dorisdump-linux-$(GOARCH)

run:
	@go run main.go $(ARGS)

gen:
	@go generate ./...

fmt:
	@go fmt .
	@goimports -l -w -local "github.com/Thearas/dorisdump" .

addcmd:
	cobra-cli --license apache --author "Thearas thearas850@gmail.com" add $(ARGS)

%:
	@:
