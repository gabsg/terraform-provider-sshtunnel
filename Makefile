TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=github.com
NAMESPACE=maclarensg
NAME=sshtunnel
BINARY=terraform-provider-${NAME}
VERSION=0.1.0
OS_ARCH=darwin_amd64
ARCH=linux_amd64

# Build the binary provider in locally, defaults for Macs
build-local:
	GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_linux_amd64
# Installs the local binary 
install-local: build-local
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test: 
	go test -i $(TEST) || exit 1                                                   
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4                    

# Creates a TF image with installed provider binary
build-tf:
	docker build --no-cache -t "${BINARY}:${VERSION}" \
	--build-arg  ARCH="${ARCH}" \
	--build-arg  BINARY="${BINARY}" \
	--build-arg  HOSTNAME="${HOSTNAME}" \
	--build-arg  NAMESPACE="${NAMESPACE}" \
	--build-arg  NAME="${NAME}" \
	--build-arg  VERSION="${VERSION}" \
	.
.PHONY: default build-local install-local release install test  
