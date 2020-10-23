PACKAGE=github.com/aureleoules/epitaf
MAIN_PACKAGE=$(PACKAGE)/cmd/epitaf
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_NAME=epitaf
STATIC_ARGS=CGO_ENABLED=0 GOOS=linux GOARCH=amd64
VERSION=`git describe --always`
COMPILE_FLAGS=-ldflags="-w -s -X $(PACKAGE)/cmd.version=$(VERSION)"

all: build
build:
	@echo "Compiling $(PACKAGE)..."
	@$(STATIC_ARGS) $(GOBUILD) -o $(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "Compiled $(BINARY_NAME)."
prod:
	@echo "Compiling $(PACKAGE) for production..."
	@$(STATIC_ARGS) $(GOBUILD) -o $(BINARY_NAME) -v $(COMPILE_FLAGS) $(MAIN_PACKAGE)
	@echo "Compiled $(BINARY_NAME)."
test: 
	$(GOTEST) -v ./...
clean: 
	@$(GOCLEAN)
	@rm -f $(BINARY_NAME)