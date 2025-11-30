ifndef GOPATH
	GOPATH := $(shell go env GOPATH)
endif
ifndef GOBIN # derive value from gopath (default to first entry, similar to 'go get')
	GOBIN := $(shell go env GOPATH | sed 's/:.*//')/bin
endif

tools = $(addprefix $(GOBIN)/, golangci-lint)

ifneq (,$(wildcard ./.local.env))
    include ./.local.env
    export
endif

.DEFAULT_GOAL := all

SYNC_ASSETS_COMMAND =	@go tool arelo\
	--target "./app/web/public" \
	--pattern "**/*.js" \
	--pattern "**/*.css" \
	--delay "100ms" \
	-- go tool templ generate --notify-proxy	

ifeq ($(OS),Windows_NT)
	SERVER_COMMAND = go tool air -c .air.win.toml
	BUILD_COMMAND = go build -o ./cloudness.exe ./cmd/app
else
	SERVER_COMMAND = go tool air -c .air.toml
	BUILD_COMMAND = go build -o ./cloudness ./cmd/app
endif

###############################################################################
#
# Initialization
#
###############################################################################

init: ## Install git hooks to perform pre-commit checks
	git config core.hooksPath .githooks
	git config commit.template .gitmessage

dep: $(deps) ## Install the deps required to generate code and build cloudness
	@echo "Installing dependencies"
	@echo "Installing go modules"
	@go mod download
	@echo "Installing go tools"
	@go install tool
	@echo "Installing npm packages"
	@npm i
	@echo "Generating templ code"
	@go tool templ generate

tools: $(tools) ## Install tools required for the build
	@echo "Installed tools"

###############################################################################
#
# dev rules
#
###############################################################################

# run: run-web run-server  ## Runs the platform server with yarn and air
templ:
	@go tool templ generate

templ-dev:
	@go tool templ generate --watch --proxy="http://localhost:8000" --open-browser=false

# run air to detect any go file changes to re-build and re-run the server.
server:
	${SERVER_COMMAND}

# run tailwindcss to generate the styles.css bundle in watch mode.
watch-assets:
	@npx @tailwindcss/cli -i app/web/assets/app.css -o ./app/web/public/assets/styles.css --watch=always

# run esbuild to generate the index.js bundle in watch mode.
watch-esbuild:
	@npx esbuild app/web/assets/index.js --bundle --outdir=app/web/public/assets --watch=forever

# watch for any js or css change in the assets/ folder, then reload the browser via templ proxy.
sync_assets:
	${SYNC_ASSETS_COMMAND}

# start the application in development
dev:
	@make -j5 templ-dev server watch-assets sync_assets


###############################################################################
#
# Build and testing rules
#
###############################################################################

build:
	@echo "Building Styles"
	@npm i
	@npx @tailwindcss/cli -i app/web/assets/app.css -o ./app/web/public/assets/styles.css
	@echo "Building Javascript"
	@npx esbuild app/web/assets/index.js --bundle --outdir=app/web/public/assets
	@echo "Building App"
	${BUILD_COMMAND}
	@echo "compiled you application with all its assets to a single binary"

test: generate  ## Run the go tests
	@echo "Running tests"
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

###############################################################################
#
# Code Formatting and linting
#
###############################################################################

format: # Format go code and error if any changes are made
	@echo "Formating ..."
	@go tool goimports -w .
	@go tool gci write --skip-generated --custom-order -s standard -s "prefix(github.com/cloudness-io/cloudness)" -s default -s blank -s dot .
	@go tool templ fmt .
	@echo "Formatting complete"

sec:
	@echo "Vulnerability detection $(1)"
	@go tool govulncheck ./...

lint: tools generate # lint the golang code
	@echo "Linting $(1)"
	@golangci-lint run --timeout=3m --verbose

###############################################################################
# Code Generation
#
# Some code generation can be slow, so we only run it if
# the source file has changed.
###############################################################################

generate: wire templ
	@echo "Generated Code"

wire: 
	@sh ./scripts/wire/wire.sh

###############################################################################
# Install Tools and deps
#
# These targets specify the full path to where the tool is installed
# If the tool already exists it wont be re-installed.
###############################################################################

update-tools: delete-tools $(tools) ## Update the tools by deleting and re-installing

delete-tools: ## Delete the tools
	@rm $(tools) || true

# Install golangci-lint
$(GOBIN)/golangci-lint:
	@echo "🔘 Installing golangci-lint... (`date '+%H:%M:%S'`)"
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) v1.56.2

help: ## show help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[$$()% 0-9a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: delete-tools update-tools help format lint swagger-gen
