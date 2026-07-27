# Local development helpers for the Modern Go Guidelines CLI.
#
# `make dev-install` builds this checkout into the tool's cache so any agent
# using the plugin runs your local changes instead of the released version.
# Re-run it after editing. `make dev-uninstall` restores the released version.

.PHONY: dev-install dev-uninstall

dev-install:
	@sh scripts/dev-install.sh install

dev-uninstall:
	@sh scripts/dev-install.sh uninstall
