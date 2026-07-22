#!/bin/sh

# Print the Go language version that is safe to use for the project rooted at
# the current working directory.

start=${1:-.}

if [ -f "$start" ]; then
    start=$(dirname "$start")
fi
start=$(CDPATH= cd "$start" 2>/dev/null && pwd -P) || {
    echo unknown
    exit 0
}

go_directive() {
    awk '$1 == "go" && $2 ~ /^[0-9]+\.[0-9]+/ { print $2; exit }' "$1" 2>/dev/null
}

# The nearest go.mod up the tree is the active module; its version wins.
directory=$start
while :; do
    if [ -f "$directory/go.mod" ]; then
        version=$(go_directive "$directory/go.mod")
        echo "${version:-unknown}"
        exit 0
    fi
    parent=$(dirname "$directory")
    if [ "$parent" = "$directory" ]; then
        break
    fi
    directory=$parent
done

# No enclosing module: use the lowest version across modules below the root.
find "$start" \
    -type d \( -name .git -o -name vendor -o -name testdata -o -name node_modules \) -prune \
    -o -type f -name go.mod -print 2>/dev/null |
    while IFS= read -r mod_file; do
        go_directive "$mod_file"
    done |
    awk -F. '
        NF >= 2 {
            major = $1 + 0
            minor = $2 + 0
            patch = $3 + 0
            if (!found || major < m1 || (major == m1 && (minor < m2 || (minor == m2 && patch < m3)))) {
                m1 = major; m2 = minor; m3 = patch
                minimum = $0
                found = 1
            }
        }
        END { print (found ? minimum : "unknown") }
    '
