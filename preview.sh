#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "${BASH_SOURCE[0]}")"
if (( $# > 1 )); then echo 'usage: ./preview.sh [bash|rust|python|go]' >&2; exit 2; fi
case ${1:-all} in
    all) languages=(bash rust python go) ;;
    bash|rust|python|go) languages=("$1") ;;
    *) echo 'usage: ./preview.sh [bash|rust|python|go]' >&2; exit 2 ;;
esac
tmp=$(mktemp -d)
trap 'rm -rf -- "$tmp"' EXIT

if [[ " ${languages[*]} " == *' rust '* ]]; then
    rustc --edition 2024 -o "$tmp/rust-demo" demo/demo.rs
    rustc --edition 2024 -o "$tmp/rust-steps" demo/steps.rs
fi
if [[ " ${languages[*]} " == *' go '* ]]; then
    mkdir "$tmp/go-demo" "$tmp/go-steps"
    cp go/kartouche.go demo/demo.go "$tmp/go-demo/"
    cp go/kartouche.go demo/steps.go "$tmp/go-steps/"
fi

run() {
    local lang=$1 kind=$2 style=${3:-}
    case "$lang:$kind" in
        bash:demo) bash demo/demo.sh ;;
        bash:steps) bash demo/steps.sh "$style" ;;
        python:demo) PYTHONPATH=python python3 demo/demo.py ;;
        python:steps) PYTHONPATH=python python3 demo/steps.py "$style" ;;
        rust:demo) "$tmp/rust-demo" ;;
        rust:steps) "$tmp/rust-steps" "$style" ;;
        go:demo) (cd "$tmp/go-demo" && go run kartouche.go demo.go) ;;
        go:steps) (cd "$tmp/go-steps" && go run kartouche.go steps.go "$style") ;;
    esac
}

unset NO_COLOR
for lang in "${languages[@]}"; do
    printf '\n== %s ==\n' "$lang"
    run "$lang" demo
    for style in inline dots crumbs rail line log; do
        printf '\n-- live %s --\n' "$style"
        run "$lang" steps "$style"
    done
done
