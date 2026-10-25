# kartouche

Copy-paste terminal banners, cards and step counters for Bash, Rust, Python and Go: one look for all my CLI scripts.

https://github.com/kamilomelo/kartouche

Copy **one** `kartouche.<ext>` into your project, edit its config block, then use its print functions. Nothing is installed or linked; the copy belongs to your script. All templates use only their language's standard library. Decoration goes to stderr; the render functions return strings for callers that want stdout. Terminal output redraws steps, while logs and pipes use append-only lines.

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

**Python** — copy `python/kartouche.py` next to your script:

```python
from kartouche import Steps, banner
banner()
with Steps("rail", ["connect", "fetch", "build"]) as steps:
    steps.next()
    steps.next()
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

`./check.sh` compares all four demos byte-for-byte with both expected files. `./preview.sh [bash|rust|python|go]` shows the look and runs each live step style in your terminal. These scripts require Bash, rustc, Python 3 and Go; they are not run automatically.

To update a script, tell your agent: **update the kartouche in this script from `<path or URL>`, following [UPDATING.md](UPDATING.md)**. The latest raw templates are at `https://raw.githubusercontent.com/kamilomelo/kartouche/main/<language>/kartouche.<ext>`.
