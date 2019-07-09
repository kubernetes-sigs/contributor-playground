include hack/Makefile.buildinfo

GOOS ?= linux
ARCH ?= amd64
REGISTRY := hub.baidubce.com/jpaas-public
BIN := cce-cloud-controller-manager
IMAGE := $(REGISTRY)/$(BIN)
SRC_DIRS := cmd pkg # directories which hold app source (not vendored)

LDFLAGS=$(VERSION_LDFLAGS)

.PHONY: all
all: build

.PHONY: build-output
build-output:
	@mkdir -p output/

.PHONY: build
build: build-output
	@GOOS=${GOOS} GOARCH=${ARCH} go build     \
	    -i                                    \
	    -o output/${BIN}  \
	    -installsuffix "static"               \
	    -ldflags "-X main.version=${VERSION}" \
	    -ldflags "${LDFLAGS}" \
	    ./cmd/cce-cloud-controller-manager

.PHONY: local-build
local-build: build-output
	@go build \
	    -i                                    \
	    -o output/${BIN}  \
	    -installsuffix "static"               \
	    -ldflags "-X main.version=${VERSION}" \
	    -ldflags "${LDFLAGS}" \
	    ./cmd/cce-cloud-controller-manager

.PHONY: image-build
image-build: build
	docker build -t ${IMAGE}:v1.11-latest .
	docker push ${IMAGE}:v1.11-latest

.PHONY: test
test:
	@./hack/test.sh $(SRC_DIRS)

.PHONY: clean
clean:
	@rm -rf output

.PHONY: deploy
deploy:
	kubectl -n kube-system set image deployments/${BIN} ${BIN}=${IMAGE}:${VERSION}

.PHONY: version
version:
	@echo ${VERSION}
