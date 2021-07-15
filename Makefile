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
		$(GOGET) github.com/gorilla/mux/@v1.8.0
		$(GOGET) github.com/joho/godotenv@v1.3.0
		$(GOGET) github.com/spyzhov/ajson@v0.4.2
		$(GOGET) github.com/yalp/jsonpath@v0.0.0-20180802001716-5cc68e5049a0
		$(GOGET) k8s.io/apimachinery@v0.21.2
		$(GOGET) k8s.io/klog/v2@v2.9.0
		$(GOMOD) tidy

test:
		@$(GOTEST) ./$(TEST_PATH) -v