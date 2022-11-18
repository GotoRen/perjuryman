GOCMD=go
GORUN=$(GOCMD) run
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GODOC=$(GOCMD)doc
GOLANGCI=golangci-lint
DUMPCMD=tcpdump
OPENSSL=openssl

INTERFACE=eth0

ROOT_CERTIFICATE_NAME=ca.pem
ROOT_CERTIFICATE_PRIVATEKEY_NAME=ca.key
SERVER_CERTIFICATE_NAME=server.local.pem
SERVER_CERTIFICATE_PRIVATEKEY_NAME=server.local.key



all: build run

build: ## go build
	$(GOBUILD) ./cmd/server/main.go

build/pi: ## go build for armv7l
	GOOS=linux GOARCH=arm GOARM=7 $(GOBUILD) ./cmd/main.go

build/pi64: ## go build for aarch64
	GOOS=linux GOARCH=arm64 $(GOBUILD) ./cmd/main.go

run: ## go run
	$(GORUN) ./cmd/server/main.go

test: ## go test
	$(GOTEST) -v ./...

benchmark: ## go benchmark
	$(GOTEST) -bench -v ./...

lint: ## golang-ci lint
	$(GOLANGCI) run --config=.golangci.yaml ./...

doc: ## godoc http:6060
	$(GODOC) -http=:6060

coverage: ## coverage
	go test -coverprofile=cover.out ./...
	go tool cover -html=cover.out -o cover.html

listen: ## check listen port
	lsof -i -P
	
dump: ## dump packet
	$(DUMPCMD) -i ${INTERFACE}

dump/out: ## dump packet and output to file
	$(DUMPCMD) -i ${INTERFACE} -w ./debug/dumpfile-server.pcapng

create/certificate: ## generate certificate
	$(GORUN) ./cmd/certificate/main.go
	@cp ca.pem ../client; echo "[INFO] 🏃 Sent the root certificate to the client..."

show/root-cert: ## show root certificate
	${OPENSSL} x509 -text -noout -in ${ROOT_CERTIFICATE_NAME}

show/server-cert: ## show server certificate
	${OPENSSL} x509 -text -noout -in ${SERVER_CERTIFICATE_NAME}

show/root-cert-privkey: ## show root certificate private key
	${OPENSSL} rsa -text -noout -in ${ROOT_CERTIFICATE_PRIVATEKEY_NAME}

show/server-cert-privkey: ## show server certificate private key
	${OPENSSL} rsa -text -noout -in ${SERVER_CERTIFICATE_PRIVATEKEY_NAME}

verify/root-cert: ## perjury: verify root certificate
	${OPENSSL} verify -CAfile ${ROOT_CERTIFICATE_NAME} ${ROOT_CERTIFICATE_NAME}

verify/server-cert: ## verify server certificate
	${OPENSSL} verify -CAfile ${ROOT_CERTIFICATE_NAME} ${SERVER_CERTIFICATE_NAME}



help: ## Display this help screen
	@grep -E '^[a-zA-Z/_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
