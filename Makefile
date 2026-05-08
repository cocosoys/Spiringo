GO ?= go
ENV ?= dev
IMAGE ?= spiringo:latest
SERVERLESS_IMAGE ?= spiringo-serverless:latest
DOCKER ?= docker
DOCKER_COMPOSE ?= docker compose
HELM ?= helm
KUBECTL ?= kubectl
GOCACHE_DIR ?= .gocache

.PHONY: test build build-cli build-serverless module crud generate payment-channel oauth-provider migrate migrate-up migrate-create docker-build docker-build-serverless docker-up docker-dev docker-down helm-lint helm-template helm-deploy k8s-apply clean

test:
	$(GO) test ./...

build:
	$(GO) build -o bin/spiringo ./cmd/spiringo

build-cli:
	$(GO) build -o bin/spiringo-cli ./cmd/spiringo-cli

build-serverless:
	$(GO) build -o bin/spiringo-serverless ./cmd/spiringo-serverless

module:
	$(GO) run ./cmd/spiringo-cli module $(MODULE)

crud:
	$(GO) run ./cmd/spiringo-cli crud $(MODULE) $(MODEL)

generate:
	TYPE=$(TYPE) MODULE=$(MODULE) MODEL=$(MODEL) sh scripts/generate.sh

payment-channel:
	$(GO) run ./cmd/spiringo-cli payment-channel $(CHANNEL)

oauth-provider:
	$(GO) run ./cmd/spiringo-cli oauth-provider $(PROVIDER)

migrate: migrate-up

migrate-up:
	$(GO) run ./cmd/spiringo migrate up -env $(ENV)

migrate-create:
	$(GO) run ./cmd/spiringo-cli migrate create $(NAME)

docker-build:
	$(DOCKER) build -t $(IMAGE) -f deployments/docker/Dockerfile .

docker-build-serverless:
	$(DOCKER) build -t $(SERVERLESS_IMAGE) -f deployments/docker/Dockerfile.serverless .

docker-up:
	$(DOCKER_COMPOSE) -f deployments/docker-compose/docker-compose.yaml up -d

docker-dev:
	$(DOCKER_COMPOSE) -f deployments/docker-compose/docker-compose.dev.yaml up -d

docker-down:
	$(DOCKER_COMPOSE) -f deployments/docker-compose/docker-compose.yaml -f deployments/docker-compose/docker-compose.dev.yaml down

helm-lint:
	$(HELM) lint deployments/helm/spiringo -f deployments/helm/spiringo/values-$(ENV).yaml

helm-template:
	$(HELM) template spiringo deployments/helm/spiringo -f deployments/helm/spiringo/values-$(ENV).yaml

helm-deploy:
	$(HELM) upgrade --install spiringo deployments/helm/spiringo -f deployments/helm/spiringo/values-$(ENV).yaml

k8s-apply:
	$(KUBECTL) apply -f deployments/kubernetes

clean:
	$(GO) clean
	$(RM) -r bin $(GOCACHE_DIR)
