#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "${BASH_SOURCE[0]}")"
tmp=$(mktemp -d)
trap 'rm -rf -- "$tmp"' EXIT

rustc --edition 2024 -o "$tmp/rust-demo" demo/demo.rs
mkdir "$tmp/go-demo"
cp go/kartouche.go demo/demo.go "$tmp/go-demo/"

demo() {
    case $1 in
        bash) bash demo/demo.sh ;;
        rust) "$tmp/rust-demo" ;;
        python) PYTHONPATH=python python3 demo/demo.py ;;
        go) (cd "$tmp/go-demo" && go run kartouche.go demo.go) ;;
    esac
}

failed=0
for lang in bash rust python go; do
    ok=true
    for mode in unicode ascii; do
        if [[ $mode == unicode ]]; then term=xterm-256color; expected=expected/demo.txt
        else term=linux; expected=expected/demo-ascii.txt; fi
        if ! NO_COLOR=1 TERM="$term" demo "$lang" > "$tmp/$lang-$mode"; then
            ok=false
        elif ! diff -u "$expected" "$tmp/$lang-$mode" >&2; then
            ok=false
        fi
    done
    if [[ $ok == true ]]; then printf '%s ok\n' "$lang"
    else printf '%s FAIL\n' "$lang"; failed=1; fi
done
exit "$failed"
