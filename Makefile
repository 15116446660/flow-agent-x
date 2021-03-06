PROJECT=flow-agent-x

LINUX_AMD64     := GOOS=linux GOARCH=amd64
MAC_AMD64       := GOOS=darwin GOARCH=amd64

GO		    	:= go
GOGEN			:= $(GO) generate ./...
GOBUILD_LINUX   := $(LINUX_AMD64) $(GO) build -o bin/$(PROJECT)-linux -v
GOBUILD_MAC     := $(MAC_AMD64) $(GO) build -o bin/$(PROJECT)-mac -v
GOTEST      	:= $(GO) test ./... -v
GOENV			:= -e GOCACHE=/ws/.cache -e GOPATH=/ws/.vender -e GO111MODULE=on

CURRENT_DIR 	:= $(shell pwd)
DOCKER_IMG		:= flowci/gosdk:1.0
DOCKER_RUN 		:= docker run -it --rm -v $(CURRENT_DIR):/ws $(GOENV) -w /ws --network host $(DOCKER_IMG) /bin/bash -c

DOCKER_BUILD 	:= ./build.sh

.PHONY: build protogen test docker clean cleanall

build:
	$(DOCKER_RUN) "$(GOGEN) && $(GOBUILD_LINUX) && $(GOBUILD_MAC)"

test:
	$(DOCKER_RUN) "$(GOTEST)"

docker: build
	$(DOCKER_BUILD) $(tag)

clean:
	$(GO) clean -i ./...
	rm -rf bin

cleanall: clean
	rm -rf .cache
	rm -rf .vender
