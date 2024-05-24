include .env.example
export

PROJECT_DIR = $(shell pwd)
APP_NAME=app
APP_PATH=$(PROJECT_DIR)/$(APP_NAME)
GO_MAIN_PATH=$(PROJECT_DIR)/cmd/app/main.go

BUILD_FLAGS=
BUILD_DEPS=
TEST_DEPS=

DOCKER_PATH=$(PROJECT_DIR)/docker
DOCKER_NAME=comment-system

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
all: test build