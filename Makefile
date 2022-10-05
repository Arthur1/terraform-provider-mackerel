export TF_ACC=1

# Values to install the provider locally for testing purposes
HOSTNAME=registry.terraform.io
NAMESPACE=mackerelio-labs
NAME=mackerel
BINARY=terraform-provider-${NAME}
VERSION=9.9.9
OS_ARCH=darwin_amd64 # change according to your environment

.PHONY: test
test:
	go test ./... -v -timeout 120m -coverprofile coverage.txt -covermode atomic

.PHONY: testacc
testacc:
	TF_ACC=1 go test -v ./mackerel/... -run $(TESTS) -timeout 120m

.PHONY: local-build
local-build:
	go build -o ${BINARY}

.PHONY: local-install
local-install: local-build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
