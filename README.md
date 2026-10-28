# kartouche

Copy-paste terminal banners, cards, step counters and command summaries for Bash, Rust, Python and Go: one look for all my CLI scripts.

https://github.com/kamilomelo/kartouche

Copy **one** `kartouche.<ext>` into your project, edit its config block, then use its print functions. Nothing is installed or linked; the copy belongs to your script. The templates need no third-party dependencies; Bash uses Arch's standard `sleep` command for the optional pause. Decoration goes to stderr; the render functions return strings for callers that want stdout. Terminal output redraws steps, while logs and pipes use append-only lines.

See [DESIGN.md](DESIGN.md) for the exact look and [SPEC.md](SPEC.md) for the full contract.

## The look

Banner and card (without color):

```text
╭─ ● app_name  v0.1.0
╰─❯ app_description

╭─ lumen ───────────  v0.3.1 ─╮
│  tiny notes, zero friction  │
╰─────────────────────────────╯
```

All six step styles, shown at step 3 of 4:

```text
inline
╭─ ● app_name  v0.1.0
╰─❯ 3/4 build

dots
╭─ ● app_name  v0.1.0
├─ ●●●○  3/4
╰─❯ build

crumbs
╭─ ● app_name  v0.1.0
├─ connect › fetch › build › deploy
╰─❯ 3/4 build

rail
╭─ ● app_name  v0.1.0
├─ ✓ connect
├─ ✓ fetch
├─ ◉ build
├─ ○ deploy
╰─❯ 3/4

line
╭─ ● app_name  v0.1.0
├─ ━━━━━━━━━━━━━━╸─────  3/4
╰─❯ build

log
╭─ ● app_name  v0.1.0
├─ 1/4 connect ······· ok
├─ 2/4 fetch ········· ok
╰─❯ 3/4 build
```

## Quickstarts

**Bash** — copy `bash/kartouche.sh` next to your script:

```bash
source ./kartouche.sh
kt_banner
kt_steps_start rail connect fetch build
kt_steps_next
kt_steps_next
kt_steps_done
```

**Rust** — copy `rust/kartouche.rs` next to `main.rs`:

```rust
mod kartouche;
use kartouche::{Steps, Style};
fn main() {
    kartouche::banner(None);
    let mut steps = Steps::start(Style::Rail, &["connect", "fetch"]);
    steps.next();
    steps.done(None);
}
```

**Python** — copy `python/kartouche.py` next to your script:

```python
from kartouche import Steps, banner
banner()
with Steps("rail", ["connect", "fetch", "build"]) as steps:
    steps.next()
    steps.next()
```

**Go** — copy `go/kartouche.go` next to `main.go`, then run `go run main.go kartouche.go`:

```go
package main
func main() {
    Banner("")
    steps := NewSteps(Rail, "connect", "fetch")
    defer steps.Close()
    steps.Next()
    steps.Done("")
}
```

If you copy the Go file into another package, change its `package main` line to match. `Steps` in Go is not concurrency-safe.

## Summaries

When several commands have already run, print a **rail** or **log** summary with one `ok`, `failed` or `skipped` outcome per label. This is separate from `Steps`: it prints once to stderr without animation, including when output is piped. The summary's last line always counts all three outcomes.

```text
rail
╭─ ● app_name  v0.1.0
├─ ✓ connect
├─ ✗ fetch
├─ ○ deploy
╰─ ✗ 1 ok · 1 failed · 1 skipped

log
╭─ ● app_name  v0.1.0
├─ 1/3 connect ······· ok
├─ 2/3 fetch ········· failed
├─ 3/3 deploy ········ skipped
╰─ ✗ 1 ok · 1 failed · 1 skipped
```

Pass finished results to the print function (shown with `rail`; use `log` for the other layout). These are calls within your script, after loading the template (inside `main` in Rust and Go):

```bash
kt_summary rail connect ok fetch failed deploy skipped
```

```rust
kartouche::summary(kartouche::Style::Rail, &[
    ("connect", kartouche::Outcome::Ok), ("fetch", kartouche::Outcome::Failed),
    ("deploy", kartouche::Outcome::Skipped),
]);
```

```python
from kartouche import summary
summary("rail", [("connect", "ok"), ("fetch", "failed"), ("deploy", "skipped")])
```

```go
Summary(Rail, []SummaryItem{{"connect", OkOutcome}, {"fetch", FailedOutcome}, {"deploy", SkippedOutcome}})
```

To build output for stdout instead, use `kt_render_summary`, `render_summary`, or `RenderSummary` in the matching language and supply stdout as the destination. Only `rail` and `log` are summary styles; the existing progress styles are unchanged.

## Spacing and pause

The **print functions** add one blank line before and after each complete banner, card, summary or steps run, so your script does not have to. A steps run gets its opening blank line at `start` and its closing blank line after `done` or `fail`; there are no blank lines between steps. Adjacent kartouche blocks each add their own blank line. Remove manually added blank lines around kartouche calls when updating an older script, or those gaps will double. The pure `render_*` functions add no surrounding spacing, so their exact output remains under your control.

On a terminal, a completed block stays visible for **one second** before kartouche writes the closing blank line and your script continues. This happens only once per block, not once per step. Edit `PAUSE_SECONDS` in a copied template (`KT_PAUSE_SECONDS` in Bash) to change the wait; fractional seconds such as `0.5` work, and `0` disables it. When stderr is piped or redirected, there is **no delay**, though the blank lines remain. The pause works with `NO_COLOR=1` and `APPEND_ONLY=true` too.

Run `./preview.sh bash` to watch the flow example and live steps in your terminal. If your script's flow messages need to stay ordered with kartouche output, write those messages to stderr as well (`>&2` in Bash).

## Accent and console

Edit `ACCENT` (or `KT_ACCENT` in Bash) to one raw SGR parameter string. `✓` stays green (`32`) and `✗` stays red regardless of the accent:

| Look | Accent | Look | Accent |
|---|---|---|---|
| yellow / amber (default) | `33` | gray | `38;5;244` |
| cyan | `36` | magenta | `35` |
| green | `32` | blue | `34` |
| soft purple | `38;5;141` | sky blue | `38;5;81` |
| warm orange | `38;5;214` | mint | `38;5;78` |
| coral | `38;5;203` | | |

Avoid `30` (black) and `90` (gray): they disappear on some terminal themes. Color is disabled on a non-terminal stream or with `NO_COLOR=1`. `GLYPHS=auto` selects ASCII when `TERM=linux` or `TERM=dumb`; an installer inside tmux on the Linux console should set `GLYPHS=ascii` and `APPEND_ONLY=true`. The card assumes single-width characters; wide characters can misalign its right edge.

## Preview and checks

`./check.sh` compares the render and print-flow demos from all four languages against their Unicode and ASCII expected files. `./preview.sh [bash|rust|python|go]` shows a live spacing/pause example and runs each step style in your terminal. The full preview takes longer because each block pauses; use a language argument for a quicker look. These scripts require Bash, rustc, Python 3 and Go; they are not run automatically.

To update a script, tell your agent: **update the kartouche in this script from `<path or URL>`, following [UPDATING.md](UPDATING.md)**. The latest raw templates are at `https://raw.githubusercontent.com/kamilomelo/kartouche/main/<language>/kartouche.<ext>`.

## License

kartouche is licensed under the [MIT License](LICENSE). Copyright (c) 2026 Kamilo Melo. You may copy and modify the templates; if you redistribute a copy, retain its copyright and SPDX header and include the MIT license text with your project.
