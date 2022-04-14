CURRENT_WORKING_DIR=$(shell pwd)

#------------------------------------------------------------------
# Project build information
#------------------------------------------------------------------
PROJNAME          		:= javatruststoreinit
VENDOR            		:= gavinmcnair
MAINTAINER        		:= Gavin McNair

CIRCLE_BUILD_NUM      ?= "unknown"
VERSION               := 1.1.$(CIRCLE_BUILD_NUM)

GIT_REPO          		:= github.com/$(VENDOR)/$(PROJNAME)
GIT_SHA           		:= $(shell git rev-parse --verify HEAD)
BUILD_DATE        		:= $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

DOCKER_USERNAME     	:= gavinmcnair
DOCKER_PASSWORD     	?="unknown"

# Construct docker image name.
IMAGE             		:= docker.io/${DOCKER_USERNAME}/$(PROJNAME):$(VERSION)

#------------------------------------------------------------------
# Go configuration
#------------------------------------------------------------------
GOCMD             		:= go
GOFMT             		:= gofmt
BIN               		:= bin

#------------------------------------------------------------------
# Build targets
#------------------------------------------------------------------

.PHONY: fmt
fmt: ## Run go fmt against code
	$(GOFMT) -w main.go
	$(GOCMD) fmt *.go

.PHONY: test
test: fmt ## Run tests
	$(GOCMD) test . -coverprofile cover.out

.PHONY: javatruststoreinit
javatruststoreinit: fmt
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOCMD) build -o $(BIN)/linux/$(PROJNAME) $(GIT_REPO)
	env GOOS=darwin GOARCH=amd64 $(GOCMD) build -o $(BIN)/darwin/$(PROJNAME) $(GIT_REPO)

.PHONY: javatruststoreinit-linux
javatruststoreinit-linux: fmt
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOCMD) build -o $(BIN)/linux/$(PROJNAME) $(GIT_REPO)

#------------------------------------------------------------------
# CI targets
#------------------------------------------------------------------

.PHONY: build
build:
	docker build \
    --build-arg git_repository=`git config --get remote.origin.url` \
    --build-arg git_branch=`git rev-parse --abbrev-ref HEAD` \
    --build-arg git_commit=`git rev-parse HEAD` \
    --build-arg built_on=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
    -t $(IMAGE) .

.PHONY: login
login:
	docker login -u $(DOCKER_USERNAME) -p $(DOCKER_PASSWORD)

push-image:
	docker push $(IMAGE)
	docker rmi $(IMAGE)

logout:
	docker logout
