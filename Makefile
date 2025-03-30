# Variables (change these as necessary)
APP_NAME := main
APP_PATH := ./cmd/app/main.go
BIN_DIR := ./bin
PKG_LIST := $(shell go list ./... | grep -v /vendor/)
GOLANGCI_LINT := golangci-lint
GOVULNCHECK := govulncheck
GOCMD := go
GOFMT := gofmt
IMAGE_NAME := enson89/order-management
TAG := latest
DOCKERFILE := docker/Dockerfile
DEPLOY_DIR := deployments
KUBECTL := kubectl

# Default target: Run checks, tests, and run the app locally
.PHONY: all
all: lint format vulncheck test run

# Build the Go binary locally
.PHONY: build
build:
	@echo "Cleaning up..."
	$(GOCMD) mod tidy
	@echo "Building binary..."
	$(GOCMD) build -o $(BIN_DIR)/$(APP_NAME) $(APP_PATH)

# Run the application locally
.PHONY: run
run: build
	@echo "Running the application locally..."
	$(BIN_DIR)/$(APP_NAME)

# Run golangci-lint to check code quality and potential issues
.PHONY: lint
lint:
	@echo "Running golangci-lint to check code quality..."
	$(GOLANGCI_LINT) run ./...

# Check Go formatting style
.PHONY: format
format:
	@echo "Checking code formatting using gofmt..."
	@$(GOFMT) -l $(PKG_LIST)
	@diff=$$($(GOFMT) -l $(PKG_LIST)); \
	if [ -n "$$diff" ]; then \
		echo "Code is not properly formatted. Run 'gofmt -w .' to fix."; \
		exit 1; \
	fi

# Check for known vulnerabilities in Go modules using govulncheck
.PHONY: vulncheck
vulncheck:
	@echo "Checking for known vulnerabilities in Go modules..."
	$(GOVULNCHECK) ./...

# Run test cases using go test
.PHONY: test
test:
	@echo "Running unit tests..."
	$(GOCMD) test ./internal/... -v

# Build the Docker image using the existing Dockerfile
.PHONY: docker-build
docker-build:
	@echo "Building Docker image using Dockerfile..."
	docker build -t $(IMAGE_NAME):$(TAG) -f $(DOCKERFILE) .

# Push the Docker image to Docker Hub or other registry
.PHONY: docker-push
docker-push:
	@echo "Pushing Docker image to registry..."
	docker push $(IMAGE_NAME):$(TAG)

# Start Minikube with Docker driver and enable Ingress
.PHONY: minikube-start
minikube-start:
	@echo "Starting Minikube with Docker driver..."
	minikube start --driver=docker
	@echo "Enabling Ingress addon..."
	minikube addons enable ingress

# Run minikube tunnel for exposing LoadBalancer services
.PHONY: minikube-tunnel
minikube-tunnel:
	@echo "Starting Minikube tunnel (requires sudo)..."
	sudo minikube tunnel

# Deploy application, Redis, and PostgreSQL manifests to Minikube
.PHONY: deploy-minikube
deploy-minikube:
	@echo "Deploying application and dependencies to Minikube..."
	$(KUBECTL) apply -f $(DEPLOY_DIR)/app.yaml
	$(KUBECTL) apply -f $(DEPLOY_DIR)/redis.yaml
	$(KUBECTL) apply -f $(DEPLOY_DIR)/postgres.yaml
	$(KUBECTL) apply -f $(DEPLOY_DIR)/ingress.yaml

# Delete Minikube and clean up resources
.PHONY: delete-minikube
delete-minikube:
	@echo "Deleting Minikube cluster and cleaning up..."
	minikube delete

# Clean up generated files and Docker images
.PHONY: clean
clean:
	@echo "Cleaning up local build files..."
	rm -f $(BIN_DIR)/$(APP_NAME)
	docker rmi $(IMAGE_NAME):$(TAG) || true