GO ?= go

GREEN_COLOR = \x1b[32;01m
RED_COLOR = \x1b[31;01m
BLUE_COLOR = \x1b[34;01m
END_COLOR = \x1b[0m

.PHONY: dep
dep:
	@echo "$(GREEN_COLOR)Installing dependencies...$(END_COLOR)"
	@$(GO) mod download
	@$(GO) mod verify

.PHONY: test
test:
	@$(GO) test -v

.PHONY: coverage
coverage:
	@$(GO) test -v -coverprofile coverage.txt

.PHONY: clean
clean:
	$(GO) clean -modcache -x -i ./...
