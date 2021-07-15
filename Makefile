GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
PROJECT_DIR=WebhookServer
RUN_PATH=$(PROJECT_DIR)/cmd/server
TEST_PATH=$(PROJECT_DIR)/test/...
BINARY_NAME=webhookserver
VAR=CGO_ENABLED=0 GOOS=linux
FLAGS=-a -installsuffix cgo
USER=yulypso
IMAGE_VERSION=v0.0.6

all: deps clean test build

build:
		$(shell $(VAR) $(GOBUILD) $(FLAGS) -o $(RUN_PATH)/$(BINARY_NAME) $(RUN_PATH)/*.go)
		$(shell docker build -t $(USER)/$(BINARY_NAME):$(IMAGE_VERSION) .)
		@echo "Go build done!"

clean:
		$(GOCLEAN)
		rm -f $(RUN_PATH)/$(BINARY_NAME)

deps:
		$(GOGET) ./...

test:
		@$(GOTEST) ./$(TEST_PATH) -v