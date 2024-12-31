SOURCES := $(shell \
	find . -not \( \( -name .git -o -name .go -o -name vendor -o -name '*.pb.go' -o -name '*.pb.gw.go' -o -name '*_gen.go' -o -name '*mock*' \) -prune \) \
	-name '*.go')

.PHONY: help
help: Makefile ## Show list of commands.
	@echo "Choose a command to run in "$(APP_NAME)":"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

#-------------------------------------------------------------------------------
#  Build and run
#-------------------------------------------------------------------------------

.PHONY: build
build: build-linux-x86_64 ## Build the project and generates a binary.
	@rm -f ./bin/meepow || true
	@go build -o ./bin/meepow ./

.PHONY: build-linux-x86_64
build-linux-x86_64: ## Build the project and generates a binary for x86_64 architecture.
	@rm -f ./bin/meepow-linux-x86_64 || true
	@env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o ./bin/meepow-linux-x86_64 ./

.PHONY: run/management-api
run/management-api: build ## Runs meepow management-api.
	@MEEPOW_API_PORT=3000 go run main.go start management-api -l development

.PHONY: run/bot
run/bot: build ## Runs meepow bot.
	@go run main.go start bot -l development

.PHONY: docker/build
docker/build: ## Build docker image.
	@docker build -t meepow .

#-------------------------------------------------------------------------------
#  Development
#-------------------------------------------------------------------------------
.PHONY: dev/management-api
dev/management-api: ## Runs meepow management-api in development mode.
	@MEEPOW_API_PORT=3000 go run main.go start management-api

.PHONY: dev/bot
dev/bot: ## Runs meepow bot.
	@go run main.go start bot