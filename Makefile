## Colors
COLOR_RESET   = \033[0m
COLOR_INFO    = \033[32m
COLOR_COMMENT = \033[33m

## Variables
MAIN_PACKAGE_PATH := ./cmd/api
BINARY_NAME := api

.PHONY: help
## Help
help:
	@printf "${COLOR_COMMENT}Usage:${COLOR_RESET}\n"
	@printf " make [target] [args...]\n\n"
	@printf "${COLOR_COMMENT}Available targets:${COLOR_RESET}\n"
	@awk '/^[a-zA-Z\-\0-9\.@]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf " ${COLOR_INFO}%-16s${COLOR_RESET} %s\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

.PHONY: tidy
## format code and tidy modfile
tidy:
	go mod tidy -v
	go fmt ./...


.PHONY: audit
## run quality control checks
audit:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...


.PHONY: run_unit_tests
## run unit tests
test:
	go test -race -buildvcs -vet=off ./...

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

.PHONY: build
## build the application
build:
	go build -ldflags='-s -w' -o=./bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}


.PHONY: run
## run the application
run: build
	./bin/${BINARY_NAME} -port ${port}


.PHONY: run/live
## run the application with reloading on file changes
run/live:
	go run github.com/air-verse/air@latest \
		--build.cmd "make build" --build.bin "./bin/${BINARY_NAME}" --build.delay "100" \
		--build.exclude_dir "" \
		--build.include_ext "go,html,js" \
		--misc.clean_on_exit "true"

