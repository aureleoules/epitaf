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
RUN = docker exec db

.PHONY: prep build

all: prep build

prep:
	@mkdir -p build

# Initialize database
init:
	$(RUN) createdb epitaf
	$(RUN) psql -h localhost -U root -W -d epitaf -f /db.sql

.PHONY:stop
stop:
	./scripts/stop.sh

.PHONY:db
db:
	./scripts/dev.sh

db-update:
	$(RUN) psql -h localhost -U root -W -d epitaf -f /db.sql

db-delete: stop
	sudo rm -rf tmp

db-client:
	docker exec -it db psql -h localhost -U root -d epitaf


.PHONY: test
test:
	./scripts/test.sh

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
	
clean: 
	@$(GOCLEAN)
	@rm -f $(BINARY_NAME)