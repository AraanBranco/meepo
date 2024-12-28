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
	@rm -f ./bin/meepo || true
	@go build -o ./bin/meepo ./

.PHONY: build-linux-x86_64
build-linux-x86_64: ## Build the project and generates a binary for x86_64 architecture.
	@rm -f ./bin/meepo-linux-x86_64 || true
	@env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o ./bin/meepo-linux-x86_64 ./

.PHONY: run/management-api
run/management-api: build ## Runs meepo management-api.
	@MEEPO_API_PORT=3000 go run main.go start management-api -l development

.PHONY: run/bot
run/bot: build ## Runs meepo bot.
	@go run main.go start bot

#-------------------------------------------------------------------------------
#  Development
#-------------------------------------------------------------------------------
.PHONY: dev/management-api
dev/management-api: ## Runs meepo management-api in development mode.
	@MEEPO_API_PORT=3000 go run main.go start management-api

.PHONY: run/bot
dev/bot: ## Runs meepo bot.
	@go run main.go start bot