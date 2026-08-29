# Maintenance helpers for the Modern Go Guidelines CLI.
#
# `make generate-features` rebuilds FEATURES.md from the guideline source data.
# `make test` runs all tests.
#
# `make dev-install` builds this checkout into the tool's cache so any agent
# using the plugin runs your local changes instead of the released version.
# Re-run it after editing. `make dev-uninstall` restores the released version.

.PHONY: generate-features test dev-install dev-uninstall

generate-features:
	@go generate ./internal/guidelines

test:
	@go test ./...

dev-install:
	@sh scripts/dev-install.sh install

dev-uninstall:
	@sh scripts/dev-install.sh uninstall
