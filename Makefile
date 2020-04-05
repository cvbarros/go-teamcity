# Project name.
PROJECT_NAME = go-teamcity

# Makefile parameters.
TAG ?= $(shell git describe)

# General.
SHELL = /bin/bash
TOPDIR = $(shell git rev-parse --show-toplevel)

# Project specifics.
BUILD_DIR = dist
PLATFORMS = linux darwin
OS = $(word 1, $@)
GOOS = $(shell uname -s | tr A-Z a-z)
GOARCH = amd64
CONTAINER_NAME = teamcity_server
INTEGRATION_TEST_DIR = integration_tests
TEAMCITY_DATA_DIR = $(INTEGRATION_TEST_DIR)/data_dir
TEAMCITY_HOST = http://localhost:8112
TEAMCITY_VERSION ?= "2019.1.1"
GO111MODULE ?= "on"

default: build

.PHONY: help
help: # Display help
	@awk -F ':|##' \
		'/^[^\t].+?:.*?##/ {
			printf "\033[36m%-30s\033[0m %s\n", $$1, $$NF \
		}' $(MAKEFILE_LIST) | sort

.PHONY: build
build: ## Build the project for the current platform
	mkdir -p $(BUILD_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BUILD_DIR)/$(PROJECT_NAME)-$(TAG)-$(GOOS)-$(GOARCH)

.PHONY: ci
ci: test ## Run all the CI targets

.PHONY: start-docker
start-docker: ## Starts up docker container running TeamCity Server
	@test -d  $(TEAMCITY_DATA_DIR) || tar xfz $(INTEGRATION_TEST_DIR)/teamcity_data.tar.gz -C $(INTEGRATION_TEST_DIR)
	@curl -sL https://download.octopusdeploy.com/octopus-teamcity/4.42.1/Octopus.TeamCity.zip -o $(TEAMCITY_DATA_DIR)/plugins/Octopus.TeamCity.zip
	@test -n "$$(docker ps -q -f name=$(CONTAINER_NAME))" || docker run --rm -d \
		--name $(CONTAINER_NAME) \
		-v $(PWD)/$(TEAMCITY_DATA_DIR):/data/teamcity_server/datadir \
		-v $(PWD)/$(INTEGRATION_TEST_DIR)/log_dir:/opt/teamcity/logs \
		-p 8112:8111 \
		jetbrains/teamcity-server:$(TEAMCITY_VERSION)
	@echo -n "Teamcity server is booting (this may take a while)..."
	@until $$(curl -o /dev/null -sfI $(TEAMCITY_HOST)/login.html);do echo -n ".";sleep 5;done

.PHONY: test
test: start-docker ## Run the unit tests
	@export TEAMCITY_ADDR=$(TEAMCITY_HOST) \
		&& GO111MODULE=$(GO111MODULE) go test -v -failfast -timeout 180s ./...

.PHONY: clean
clean: clean-code clean-docker ## Clean all resources (!DESTRUCTIVE!)

.PHONY: clean-code
clean-code: ## Remove unwanted files in this project (!DESTRUCTIVE!
	@cd $(TOPDIR) && git clean -ffdx && git reset --hard

.PHONY: clean-docker
clean-docker: ## Remove the docker container if it is running
	@docker rm -f $(CONTAINER_NAME)

.PHONY: dist
dist: $(PLATFORMS) ## Package the project for all available platforms

.PHONY: setup
setup: ## Setup the full environment
	dep ensure

.PHONY: $(PLATFORMS)
$(PLATFORMS): # Build the project for all available platforms
	mkdir -p $(BUILD_DIR)
	GOOS=$(OS) GOARCH=$(GOARCH) go build -o $(BUILD_DIR)/$(PROJECT_NAME)-$(TAG)-$(OS)-$(GOARCH)
