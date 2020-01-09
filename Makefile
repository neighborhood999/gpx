GO ?= go

GREEN_COLOR = \x1b[32;01m
RED_COLOR = \x1b[31;01m
BLUE_COLOR = \x1b[34;01m
END_COLOR = \x1b[0m

.PHONY: echo
echo:
	@echo "$(GREEN_COLOR)Hello, World!$(END_COLOR)"

.PHONY: test
test:
	@$(GO) test -v

.PHONY: coverage
coverage:
	@$(GO) test -v -coverprofile coverage.txt

.PHONY: clean
clean:
	$(GO) clean -modcache -x -i ./...
