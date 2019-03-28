BINARY = amazon-reviews-crawler-go
VET_REPORT = vet.report
GOARCH = amd64

VERSION=0.0.1
COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

# Symlink into GOPATH
GITHUB_USERNAME=anzellai
BUILD_DIR=${GOPATH}/src/github.com/${GITHUB_USERNAME}/${BINARY}
CURRENT_DIR=$(shell pwd)
BUILD_DIR_LINK=$(shell readlink ${BUILD_DIR})
# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-s -w -X main.VERSION=${VERSION} -X main.COMMIT=${COMMIT} -X main.BRANCH=${BRANCH}"

# Build the project
all: clean vet linux darwin windows

linux:
	GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BUILD_DIR}/bin/${BINARY}-linux cmd/crawler/*.go

darwin:
	GOOS=darwin GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BUILD_DIR}/bin/${BINARY}-darwin cmd/crawler/*.go

windows:
	GOOS=windows GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BUILD_DIR}/bin/${BINARY}-windows cmd/crawler/*.go

vet:
	-cd ${BUILD_DIR}; \
	go vet ./... > ${VET_REPORT} 2>&1 ; \
	cd - >/dev/null

fmt:
	cd ${BUILD_DIR}; \
	go fmt $$(go list ./... | grep -v /vendor/) ; \
	cd - >/dev/null

build:
	make clean; \
	make fmt; \
	make vet; \
	make linux; \
	make darwin;
	make windows;

clean:
	-rm -f ${BUILD_DIR}/${VET_REPORT}
	-rm -rf ${BUILD_DIR}/bin

.PHONY: linux darwin windows vet fmt clean build
