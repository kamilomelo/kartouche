#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "${BASH_SOURCE[0]}")"
tmp=$(mktemp -d)
trap 'rm -rf -- "$tmp"' EXIT

rustc --edition 2024 -o "$tmp/rust-demo" demo/demo.rs
rustc --edition 2024 -o "$tmp/rust-flow" demo/flow.rs
mkdir "$tmp/go-demo" "$tmp/go-flow"
cp go/kartouche.go demo/demo.go "$tmp/go-demo/"
cp go/kartouche.go demo/flow.go "$tmp/go-flow/"

demo() {
    case $1 in
        bash) bash demo/demo.sh ;;
        rust) "$tmp/rust-demo" ;;
        python) PYTHONPATH=python python3 demo/demo.py ;;
        go) (cd "$tmp/go-demo" && go run kartouche.go demo.go) ;;
    esac
}

flow() {
    case $1 in
        bash) bash demo/flow.sh ;;
        rust) "$tmp/rust-flow" ;;
        python) PYTHONPATH=python python3 demo/flow.py ;;
        go) (cd "$tmp/go-flow" && go run kartouche.go flow.go) ;;
    esac
}

check_bad_bash_count() {
    bash -c 'source bash/kartouche.sh; kt_render_frame 1 line "$1" "$2" build connect fetch build deploy' _ "$1" "$2" > "$tmp/bad-output" 2> "$tmp/bad-error" &&
        [[ ! -s $tmp/bad-output && ! -s $tmp/bad-error ]]
}

failed=0
for lang in bash rust python go; do
    ok=true
    for mode in unicode ascii; do
        if [[ $mode == unicode ]]; then term=xterm-256color; expected=expected/demo.txt; flow_expected=expected/flow.txt
        else term=linux; expected=expected/demo-ascii.txt; flow_expected=expected/flow-ascii.txt; fi
        if ! NO_COLOR=1 TERM="$term" demo "$lang" > "$tmp/$lang-$mode"; then
            ok=false
        elif ! diff -u "$expected" "$tmp/$lang-$mode" >&2; then
            ok=false
        fi
        if ! NO_COLOR=1 TERM="$term" flow "$lang" > "$tmp/$lang-$mode-flow-stdout" 2> "$tmp/$lang-$mode-flow"; then
            ok=false
        elif [[ -s $tmp/$lang-$mode-flow-stdout ]] || ! diff -u "$flow_expected" "$tmp/$lang-$mode-flow" >&2; then
            ok=false
        fi
    done
    if [[ $lang == bash ]]; then
        probe='arr[$(printf INJECTION >&2; echo 0)]'
        if ! check_bad_bash_count "$probe" 4 || ! check_bad_bash_count 1 "$probe" || ! check_bad_bash_count 1 0; then
            ok=false
        fi
    fi
    if [[ $ok == true ]]; then printf '%s ok\n' "$lang"
    else printf '%s FAIL\n' "$lang"; failed=1; fi
done
exit "$failed"
