#!/usr/bin/env bash
set -euo pipefail

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)

if [[ "$(go env CGO_ENABLED)" != "1" ]]; then
	echo "examples-check requires CGO_ENABLED=1 for the SQLite example" >&2
	exit 1
fi

shopt -s nullglob
examples=("$repo_root"/examples/example-*.md)
if ((${#examples[@]} == 0)); then
	echo "No example Markdown files found in $repo_root/examples" >&2
	exit 1
fi

tmp_dir=$(mktemp -d)
trap 'rm -rf "$tmp_dir"' EXIT

extract_go_block() {
	local markdown=$1
	local output=$2

	if ! awk '
		/^```go[[:space:]]*$/ {
			if (in_go) bad = 1
			in_go = 1
			blocks++
			next
		}
		/^```[[:space:]]*$/ && in_go {
			in_go = 0
			closed++
			next
		}
		in_go { print }
		END {
			if (blocks != 1 || closed != 1 || in_go || bad) exit 1
		}
	' "$markdown" >"$output"; then
		printf '%s must contain exactly one complete ```go block\n' "$markdown" >&2
		return 1
	fi
}

for markdown in "${examples[@]}"; do
	name=${markdown##*/}
	name=${name%.md}
	dir="$tmp_dir/$name"
	mkdir -p "$dir"
	extract_go_block "$markdown" "$dir/api_test.go"
done

cd "$tmp_dir"
go mod init github.com/kematzy/zentests-examples-check >/dev/null 2>&1
go mod edit -go=1.25.5
go mod edit -require=github.com/kematzy/zentests@v0.0.0
go mod edit "-replace=github.com/kematzy/zentests=$repo_root"
go mod tidy

echo "Checking Markdown examples with go vet..."
go vet ./...
echo "Running Markdown examples..."
go test ./...
