PACKAGE=github.com/aureleoules/epitaf
MAIN_PACKAGE=$(PACKAGE)
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BUILD_DIRECTORY=./
BINARY_NAME=epitaf
STATIC_ARGS=CGO_ENABLED=0 GOOS=linux GOARCH=amd64
VERSION=`git describe --always`
COMPILE_FLAGS=-ldflags="-w -s -X $(PACKAGE)/cmd.version=$(VERSION)"

.PHONY: prep build

all: prep build

prep:
	@mkdir -p build

swag:
	@echo "Building swagger documentation..."
	@swag init --parseDependency
	@echo "Built swagger documentation."

build:
	@echo "Compiling $(PACKAGE)..."
	@$(STATIC_ARGS) $(GOBUILD) -o $(BUILD_DIRECTORY)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "Compiled $(BINARY_NAME)."
prod:
	@echo "Compiling $(PACKAGE) for production..."
	@$(STATIC_ARGS) $(GOBUILD) -o $(BUILD_DIRECTORY)/$(BINARY_NAME) -v $(COMPILE_FLAGS) $(MAIN_PACKAGE)
	@echo "Compiled $(BINARY_NAME)."
test: 
	$(GOTEST) -v ./...
clean: 
	@$(GOCLEAN)
	@rm -f $(BINARY_NAME)