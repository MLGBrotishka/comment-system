include .env.example
export

PROJECT_DIR = $(shell pwd)
APP_NAME=app
APP_PATH=$(PROJECT_DIR)/bin/$(APP_NAME)
GO_MAIN_PATH=$(PROJECT_DIR)/cmd/app/main.go

BUILD_FLAGS=
BUILD_DEPS=
TEST_DEPS=

ENV_FILE=$(PROJECT_DIR)/.env.example

DOCKER_PATH=$(PROJECT_DIR)/docker
DOCKER_NAME=comment-system
DOCKER_CMD=docker-compose --env-file $(ENV_FILE) -p $(DOCKER_NAME) --project-directory $(PROJECT_DIR) -f $(DOCKER_PATH)/docker-compose.yaml

# Local run

.PHONY: build
build:
	go build $(BUILD_FLAGS) -o $(APP_PATH) $(GO_MAIN_PATH) $(BUILD_DEPS)

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: run
run: build
	export IN_MEMORY_STORAGE=true && \
	$(APP_PATH)

.PHONY: clean
clean:
	go clean
	rm -f $(APP_PATH)

.PHONY: all
all: build

# Docker-compose

.PHONY: docker
docker:
	$(DOCKER_CMD) up 

.PHONY: docker-clear
docker-clear:
	$(DOCKER_CMD) down --volumes

.PHONY: docker-rebuild
docker-rebuild:
	$(DOCKER_CMD) build
	$(DOCKER_CMD) up