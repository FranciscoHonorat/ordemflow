SERVICES = order-service inventory-service notification-service analytics-service cdc-connector
ORDER_SVC = ./services/order-service

# Build
build-check:
	cd $(ORDER_SVC) && go build ./...

build:
	cd $(ORDER_SVC) && go build -o ../../bin/order-service ./cmd/server

# Test
test:
	cd $(ORDER_SVC) && go test ./...

test-verbose:
	cd $(ORDER_SVC) && go test -v ./...


test-benchmark:
	cd $(ORDER_SVC) && go test -bench=. ./...

# Dev
run:
	cd $(ORDER_SVC) && go run ./cmd/server

# Quality
tidy:
	cd $(ORDER_SVC) && go mod tidy

lint:
	cd $(ORDER_SVC) && golangci-lint run ./...

# Roda em todos os serviços
tidy-all:
	$(foreach svc, $(SERVICES), cd ./services/$(svc) && go mod tidy;)

build-all:
	$(foreach svc, $(SERVICES), cd ./services/$(svc) && go build ./...;)

# Clean
clean:
	rm -rf bin/

.PHONY: build-check build build-all test test-verbose test-benchmark run tidy tidy-all lint clean