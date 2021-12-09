all: build tag push
test: lint unit-test

NAME=rck-auth
VERSION=0.1
REGISTRY="quay.io/rcarrata"
TOOL="docker"

install:
	@go build .

build: 
	@go version
	@${TOOL} build -t localhost/${NAME}:${VERSION} .
	
tag:
	@${TOOL} tag localhost/${NAME}:${VERSION} ${REGISTRY}/${NAME}:${VERSION}

push: 
	@${TOOL} push ${REGISTRY}/${NAME}:${VERSION}