#!/usr/bin/env sh
# Developer-only helper: build this checkout into the tool's cache so any agent
# using the plugin runs your local changes instead of the released version.
set -eu

binary_name="go-modern-guidelines"

if [ -n "${XDG_CACHE_HOME:-}" ]; then
	cache_root="${XDG_CACHE_HOME}/go-modern-guidelines"
elif [ -n "${HOME:-}" ]; then
	cache_root="${HOME}/.cache/go-modern-guidelines"
else
	echo "go-modern-guidelines: HOME or XDG_CACHE_HOME must be set" >&2
	exit 1
fi

dev_dir="${cache_root}/dev"
dev_binary="${dev_dir}/${binary_name}"

# This script lives at <repo>/scripts/dev-install.sh, so the module root is the
# parent of its own directory.
script_dir="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd -P)"
module_dir="$(CDPATH= cd -- "${script_dir}/.." && pwd -P)"

case "${1:-install}" in
install)
	if ! command -v go >/dev/null 2>&1; then
		echo "go-modern-guidelines: Go toolchain is required to build the dev binary" >&2
		exit 1
	fi
	if [ ! -f "${module_dir}/go.mod" ]; then
		echo "go-modern-guidelines: no Go module found at ${module_dir}" >&2
		exit 1
	fi
	mkdir -p "${dev_dir}"
	staged="${dev_binary}.tmp.$$"
	trap 'rm -f "${staged}"' EXIT HUP INT TERM
	(
		cd "${module_dir}"
		GOFLAGS= GOWORK=off CGO_ENABLED=0 go build -o "${staged}" .
	)
	mv "${staged}" "${dev_binary}"
	trap - EXIT HUP INT TERM
	echo "go-modern-guidelines: installed dev build to ${dev_binary}" >&2
	echo "go-modern-guidelines: set GO_MODERN_GUIDELINES_DEV=1 to use it" >&2
	;;
uninstall)
	rm -rf "${dev_dir}"
	echo "go-modern-guidelines: removed dev build (${dev_binary})" >&2
	;;
*)
	echo "usage: sh scripts/dev-install.sh [install|uninstall]" >&2
	exit 2
	;;
esac
