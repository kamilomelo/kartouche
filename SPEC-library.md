> **ARCHIVED, FUTURE UPGRADE PATH.** This is the full specification for turning kartouche into a maintained, versioned library (single-file per language, conformance tests, packaging). The active plan is the template model in `SPEC.md`. If you outgrow copy-and-modify templates, hand this file to your coding agent. The visual rules in section 4 and the color rules in section 5 are shared by both plans.
>
> Note: decisions made after this file was written (Arch Linux as the reference platform, no macOS support, MIT license, the 256-color index kept) are already reflected below.

# kartouche: specification

> Working name: **kartouche** (a cartouche is the ornamental frame around a name).
> The name is a single find-and-replace away from changing (candidates: kartucho, cairn).
> Check availability on crates.io, PyPI, npm and GitHub before publishing.

A tiny, **dependency-free, embeddable** library of **terminal banners and step counters** for scripts and small TUIs, with identical output in **Bash**, **Rust**, **Python** and **Go**. Each implementation is meant to be copied into other projects as a single file.

---

## 0. Instructions for the implementing agent

Read this whole document before writing code. Then:

1. Build in the phases of section 13, in order. Finish and test each phase before the next.
2. **The expected-output files in `spec/expected/` are the source of truth.** All four implementations must reproduce them byte for byte. If code and this document disagree, fix the code, unless the document is clearly wrong; in that case fix the document and the expected files together and say so.
3. **No dependencies at all**, not even optional ones, in any language: standard library only. Each implementation lives in a single source file that can be copied into another project (section 9).
4. Keep it lean. If something is not in this spec, do not add it. Non-goals are in section 2.
5. Never print, exit, panic or throw for a rendering problem. Degrade to plain text instead.
6. Every public function needs a doc comment with a one-line example.

---

## 1. Goals

- Make small CLIs and scripts look polished with about two lines of code.
- Same look, same environment-variable behavior, same names across Bash, Rust, Python and Go.
- Safe by default: correct output in pipes, CI logs, `NO_COLOR` environments and non-UTF-8 terminals.
- **Embeddable first:** copy one file into any script or repo (`kartouche.sh`, `lib.rs`, `kartouche.py`, `kartouche.go`) and it works. Normal package-manager installs are a convenience on top.

## 2. Non-goals

- Not a TUI framework: no input handling, no layout engine, no widgets.
- No threads, timers or background animation in v0.1 (spinners are a v0.2 idea, see section 14).
- No Windows legacy console support (Windows Terminal and VT-enabled consoles only).
- No dependencies, optional or otherwise, and no feature flags. No wrapper binary.
- No themes file or config file in v0.1.

---

## 3. Components

| Component | What it is | Version |
|---|---|---|
| **Banner** | Two-line title block with accent dot, name, version, description | v0.1 |
| **Card** | Boxed variant with the name in the top border | v0.1 |
| **Steps** | A step counter shown under the banner, acting as a subtitle: where the app is going | v0.1 |
| **Status line** | Optional `├─` line for live facts (`12 hosts · 11 up`) | v0.1 |
| **Accent colors** | The 16 ANSI names plus `0-255` indexes, with a predictable fallback | v0.1 |

---

## 4. Visual specification

All examples are monospace. `·` marks nothing; spaces are real spaces.

### 4.1 Banner

```
╭─ ● app_name  v0.1.0
╰─❯ app_description
```

Rules:

- Line 1: `╭─`, space, accent `●`, space, **name** (bold), two spaces, `v` + version (dim).
- Line 2: `╰─❯`, space, description (default foreground).
- Frame characters (`╭─`, `╰─❯`, and any `├─`) are **dim**.
- If the caller passes a version that already starts with `v`, do not add a second `v`.
- If the version is empty or missing, omit the two spaces and the version.
- If the description is empty, omit line 2 entirely.
- Truncate the description with `…` so no line exceeds the terminal width (section 8.4).

With an optional status line between the two (the middle connector uses `├─`):

```
╭─ ● fleetty  v0.2.0
├─ 12 hosts · 11 up · 1 down
╰─❯ your fleet, one terminal
```

### 4.2 Card

```
╭─ lumen ────────────  v0.3.1 ─╮
│  tiny notes, zero friction   │
╰──────────────────────────────╯
```

Rules:

- The name sits in the top border; the version sits at the right of the top border.
- Inner width `w = max(display_width(description) + 4, left + right + 2)` where `left = display_width("─ " + name + " ")` and `right = display_width("  v" + version + " ─")`.
- Top fill is `w - left - right` copies of `─`. The body line pads the description to `w - 4` and has two spaces of padding on each side.
- Bottom border is `╰`, `w` copies of `─`, `╯`.
- Frame is dim, name is bold in the default foreground, the accent color applies to nothing in the card except an optional `●` before the name when `dot=true`.

### 4.3 Steps

Steps render **below the banner** as a subtitle. A run has `N` steps (known up front) or an unknown total (see 4.3.7). `n` is the current step, 1-based. There are six styles. Every style renders the same information; the user picks one with `style`.

The banner is printed **once**. Only the steps region is redrawn (section 7).

#### 4.3.1 `inline` (default)

```
╭─ ● app_name  v0.1.0
╰─❯ 2/5 building images
```

The counter `2/5` is dim, the label is normal.

#### 4.3.2 `dots`

```
╭─ ● app_name  v0.1.0
├─ ●●○○○  2/5
╰─❯ building images
```

`●` repeated `n` times (accent), `○` repeated `N - n` times (dim), two spaces, counter (dim). Cap at 12 dots; above that fall back to `line`.

#### 4.3.3 `crumbs`

```
╭─ ● app_name  v0.1.0
├─ connect › fetch › build › deploy
╰─❯ 3/4 build
```

Done labels dim, current label bold accent, pending labels dim. Needs all labels up front. If the joined line is wider than the terminal, fall back to `inline`.

#### 4.3.4 `rail`

```
╭─ ● app_name  v0.1.0
├─ ✓ connect
├─ ✓ fetch
├─ ◉ build
├─ ○ deploy
╰─❯ 3/4
```

`✓` accent, `◉` accent bold, `○` dim, label of the current step bold. Needs all labels up front. The block grows in place and is redrawn on every change.

#### 4.3.5 `line`

```
╭─ ● app_name  v0.1.0
├─ ━━━━━━━╸────────────  2/5
╰─❯ building images
```

Bar width is 20 cells. `filled = round(20 * n / N)`. If `filled < 20`, draw `filled - 1` copies of `━`, one `╸`, then `20 - filled` copies of `─`. If `filled == 20`, draw 20 `━`. `filled = 0` draws 20 `─`. Filled part is accent, the rest is dim. Then two spaces and the counter.

#### 4.3.6 `log`

```
╭─ ● app_name  v0.1.0
├─ 1/4 connect ······· ok
├─ 2/4 fetch ········· ok
╰─❯ 3/4 build
```

Each finished step stays on screen as `├─ n/N label` + a dotted leader of `·` + `ok` (accent) or `fail` (red). Leaders pad so the result column aligns at `max(label_width) + 12`. The last line is the active step. This style appends lines and never redraws earlier ones.

#### 4.3.7 Unknown total

If `N` is unknown, render the counter as `n` alone (`3 building images`) in `inline` and `log`; other styles are not available and fall back to `inline`.

### 4.4 End states

Replace the last line (`╰─❯ ...`) with one of:

```
╰─ ✓ done in 4.2s
╰─ ✗ failed at 3/4 build
```

- `✓` uses the accent color; `✗` is red (SGR 31); this is the only place red is used by default.
- `done in` is optional; show elapsed time only when the caller started a timer (`Steps` records its own start time). Format: `4.2s` under 60 s, then `1m 07s`, then `1h 02m`.
- A `fail` message may add a reason after the label, separated by ` · `.

### 4.5 ASCII fallback

When Unicode is unavailable (section 6), use this glyph table. Layout is otherwise identical.

| Unicode | ASCII | Unicode | ASCII |
|---|---|---|---|
| `╭` | `,` | `✓` | `+` |
| `╰` | `` ` `` | `✗` | `x` |
| `╮` | `.` | `›` | `>` |
| `╯` | `'` | `━` | `=` |
| `─` | `-` | `╸` | `=` |
| `│` | `\|` | `·` | `.` |
| `├` | `\|` | `…` | `...` |
| `❯` | `>` | `●` | `*` |
| `○` | `o` | `◉` | `@` |

ASCII banner:

```
,- * app_name  v0.1.0
`-> app_description
```

---

## 5. Colors

### 5.1 Roles

| Role | Style |
|---|---|
| frame (`╭ ╰ ├ ─ │ ❯` and the card border) | dim: `\e[2m` |
| name | bold: `\e[1m` |
| version, counter, pending items | dim |
| accent (`●`, `✓`, done items, filled bar) | accent color |
| failure (`✗`, `fail`) | red: `\e[31m` |
| reset | `\e[0m` after every styled run |

Never leave a style open at the end of a line. Never emit empty style sequences.

### 5.2 Accent values

The accent accepts either of:

- **Name** (case-insensitive; `-`, `_` and spaces ignored, so `bright-red`, `bright_red` and `BrightRed` match): `black red green yellow blue magenta cyan white gray` (alias `grey`, `brightblack`), `brightred brightgreen brightyellow brightblue brightmagenta brightcyan brightwhite`.
- **Index**: an integer `0` to `255`, passed straight through as a 256-color value.

| Name | SGR fg | Name | SGR fg |
|---|---|---|---|
| black | 30 | gray | 90 |
| red | 31 | brightred | 91 |
| green | 32 | brightgreen | 92 |
| yellow | 33 | brightyellow | 93 |
| blue | 34 | brightblue | 94 |
| magenta | 35 | brightmagenta | 95 |
| cyan | 36 | brightcyan | 96 |
| white | 37 | brightwhite | 97 |

Encoding: names `\e[<code>m`; index `\e[38;5;<n>m`.

**Not supported in 0.1:** hex (`#rrggbb`), truecolor (24-bit), and any nearest-color conversion. Such values are treated as unparseable (see below). They are a roadmap item (section 14).

An unparseable or out-of-range value (`256`, `-1`, `#af87ff`, `banana`) silently becomes the default, **cyan (36)**. The library never errors on a bad color.

### 5.3 Fallback rule (one line, no math)

A **256-capable** terminal is one where `TERM` contains `256color`, or `COLORTERM` is set and non-empty.

- Accent is a **name**: always used as is (every terminal supports the 16).
- Accent is an **index** and the terminal is 256-capable: used as is.
- Accent is an **index** and the terminal is **not** 256-capable: use the default, **cyan (36)**.

There is no palette table and no color-distance calculation anywhere in the library.

### 5.4 Readability notes (for the README)

- The 16 names follow the user's terminal theme, which is why they are the recommended choice for shared tools.
- `black` is invisible on dark themes, and `gray` (bright black) is invisible on some themes such as Solarized Dark. Both are accepted but discouraged as accents; the docs say so.
- Suggested accents to showcase in the README. The standard colors come first because they follow the user's terminal theme; the indexes are optional extras for a more distinctive look:

| Look | Value | Kind |
|---|---|---|
| cyan (default) | `cyan` | standard, follows theme |
| magenta | `magenta` | standard, follows theme |
| yellow / amber | `yellow` | standard, follows theme |
| green | `green` | standard, follows theme |
| blue | `blue` | standard, follows theme |
| soft purple | `141` | 256-color index |
| sky blue | `81` | 256-color index |
| warm orange | `214` | 256-color index |
| mint | `78` | 256-color index |
| coral | `203` | 256-color index |

The `examples/` scripts show one accent per fictional app (for example cyan, magenta, yellow) to demonstrate that each tool can have its own identity.

---

## 6. Environment and stream behavior

These rules are identical in every language.

### 6.1 Output stream

The library writes to **stderr by default**, so decoration never pollutes data piped from stdout. The caller may choose `stdout` instead. All terminal checks (6.2, 7) use the **stream actually written to**.

### 6.2 Color decision, evaluated in this order

1. `NO_COLOR` set and non-empty: **no color**.
2. `FORCE_COLOR` set, non-empty and not `0`: **color on**, even when not a terminal.
3. Stream is a TTY and `TERM` is not `dumb`: **color on**.
4. Otherwise: **no color**.

With color off, output contains no escape sequences at all: not even resets.

### 6.3 Unicode decision

ASCII glyphs are used when any of these holds, checked in order:

1. `KARTOUCHE_ASCII` is set and not `0`: ASCII.
2. `KARTOUCHE_ASCII=0`: Unicode.
3. `TERM=dumb`: ASCII.
4. The first non-empty of `LC_ALL`, `LC_CTYPE`, `LANG` exists and does not contain `UTF-8` or `utf8` (case-insensitive): ASCII.
5. Otherwise: Unicode.

### 6.4 Accent override

Order of precedence, highest first:

1. App-specific variable named by the caller (for example `FLEETTY_ACCENT`).
2. `KARTOUCHE_ACCENT`.
3. The accent argument passed in code.
4. Default: cyan.

Environment wins over code so end users can theme any tool without editing it.

### 6.5 Summary table

| Variable | Effect |
|---|---|
| `NO_COLOR` | disables color |
| `FORCE_COLOR` | forces color on (when not `0`) |
| `KARTOUCHE_ASCII` | forces ASCII (`1`) or Unicode (`0`) |
| `KARTOUCHE_ACCENT` | accent override |
| `<APP>_ACCENT` | app-specific accent override, name supplied by caller |
| `COLUMNS` | terminal width hint (8.4) |
| `TERM`, `COLORTERM` | 256-color capability (5.3), `dumb` detection |
| `CI` | when non-empty, steps never redraw (7.2) |

---

## 7. Redraw behavior (Steps)

### 7.1 Interactive terminal

When the stream is a TTY and `CI` is unset:

- Print the banner once and keep it.
- Print the steps region (1 to `N + 2` lines depending on style; call it `h` lines).
- On each change: move up `h` lines (`\e[<h>A`), `\r`, then for each line clear it (`\e[2K`) and rewrite it, finishing with the cursor on the last line.
- Hide the cursor (`\e[?25l`) during the run and **always** restore it (`\e[?25h`) on normal exit, error, `Ctrl-C` and `SIGTERM`.
- On finish (`done` or `fail`), do a final redraw with the end state, restore the cursor, and end with a newline.

### 7.2 Non-interactive output

When the stream is not a TTY, or `CI` is set:

- **No cursor movement and no clearing.** Output is append-only so logs stay readable.
- The banner is printed once as usual (with color per 6.2).
- Each step is appended as one line when it **starts**, using the `inline` layout regardless of the chosen style, prefixed with the same frame glyph:

```
╭─ ● app_name  v0.1.0
├─ 1/5 connecting
├─ 2/5 building images
├─ 3/5 pushing
╰─ ✓ done in 4.2s
```

(Intermediate lines use `├─`; the final line uses `╰─ ✓` or `╰─ ✗`.)

### 7.3 Cleanup rules per language

- **Bash**: install `trap` handlers for `EXIT INT TERM` that restore the cursor. **Chain** any existing trap instead of replacing it (read it with `trap -p`).
- **Rust**: `Drop` on the `Steps` guard restores the cursor. A guard dropped without `finish()` or `fail()` prints a quiet end line `╰─ ✗ interrupted at n/N label`.
- **Python**: `Steps` is a context manager; `__exit__` finishes or fails depending on whether an exception is in flight. Also register `atexit` as a safety net for the cursor.
- **Go**: `defer s.Close()` restores the cursor. While steps are active on a TTY, the library installs a `SIGINT`/`SIGTERM` handler that restores the cursor, prints the `interrupted` end line, then exits with code 130 (or 143 for `SIGTERM`). The handler is removed when the steps finish.

---

## 8. Measuring text

### 8.1 Display width

Alignment (card, log leaders, truncation) uses **display width**, not byte or char count:

- Combining marks and zero-width characters: 0.
- East Asian Wide and Fullwidth characters, and most emoji: 2.
- Everything else: 1.
- ANSI escape sequences: 0 (strip before measuring).

### 8.2 Per-language approach

No dependencies anywhere, so each language carries a **small internal width table** covering what matters: combining marks and zero-width characters (0), CJK, Hangul, fullwidth forms and common emoji ranges (2).

- **Rust**: internal table in the single file.
- **Python**: `unicodedata.east_asian_width` and `unicodedata.combining` from the stdlib.
- **Go**: internal table over `unicode` package range tables plus explicit wide ranges.
- **Bash**: `${#var}` (character count under a UTF-8 locale). Wide characters may misalign; this is the accepted limitation.

### 8.3 Frame glyph widths

All frame, status and marker glyphs in 4.5 are assumed width 1 (ambiguous-width characters are treated as narrow).

### 8.4 Terminal width

Resolved in this order, identically in every language:

1. The explicit override: Bash `KT_WIDTH`, Rust `.width(n)`, Python `width=n`, Go `.Width(n)`.
2. `COLUMNS` if it is a positive integer.
3. A stdlib-only terminal query where one exists: Python `shutil.get_terminal_size`, Bash `tput cols 2>/dev/null`.
4. **80**.

Rust and Go have no portable stdlib terminal-size call, and `unsafe`/syscalls are avoided on purpose, so they use steps 1, 2 and 4 only. Apps that care can pass the width themselves.

Truncate text to `width - 1` with `…`. If the width is below 20, skip truncation entirely.

---

## 9. Repository layout and embedding

```
kartouche/
├── README.md                 # text-block screenshots, embed instructions, 10-line quickstarts
├── SPEC.md                   # this document
├── LICENSE                   # MIT, full text in section 9.1
├── CHANGELOG.md
├── spec/
│   ├── expected/             # expected outputs, source of truth
│   │   ├── plain/            # *.txt, no escape sequences
│   │   ├── ansi/             # *.ansi, with escape sequences (FORCE_COLOR=1)
│   │   └── ascii/            # *.txt, ASCII glyph mode
│   └── cases.txt             # demo-CLI invocations, one per expected file
├── bash/
│   ├── kartouche.sh          # the whole library, single file
│   ├── kt-demo               # demo CLI (section 11)
│   └── test/                 # bats-core tests
├── rust/
│   ├── Cargo.toml
│   ├── src/lib.rs            # the whole library, single self-contained file
│   ├── examples/kt_demo.rs   # demo CLI
│   └── tests/                # expected-output and unit tests
├── python/
│   ├── pyproject.toml
│   ├── src/kartouche.py      # the whole library, single module (also runs as `python -m kartouche`)
│   └── tests/                # pytest
├── go/
│   ├── go.mod                # module path ends in /go (see tagging, section 13)
│   ├── kartouche.go          # the whole library, single file, package kartouche
│   ├── kartouche_test.go
│   └── cmd/kt-demo/main.go   # demo CLI
├── examples/                 # one realistic script per language
├── scripts/
│   ├── conformance.sh        # runs every demo CLI against spec/expected
│   └── embed-check.sh        # copies each single file into a scratch project and builds/runs it
└── .github/workflows/ci.yml
```

### Embedding rules

- **Dependencies: none.** Standard library only in every language, no optional extras, no feature flags. Dev-only tools (`bats-core`, `shellcheck`, `pytest`, `ruff`, `mypy`, `go vet`, `staticcheck`, `clippy`) are allowed in CI but never in the library.
- **One file per language** contains the entire library. It must not `import`, `use`, `require` or `source` any other file of this repo.
- **Copy-friendly:** each file starts with a header comment containing the version, the line `SPDX-License-Identifier: MIT`, the copyright line, the repo URL and the one-line instruction "copy this file into your project". The MIT notice travels with every copied file, which is what the license requires.
  - Bash: `source ./kartouche.sh`
  - Rust: `mod kartouche;` (file saved as `kartouche.rs`), or `cargo add kartouche`
  - Python: `import kartouche` (file saved as `kartouche.py`), or `pip install kartouche`
  - Go: save as `kartouche.go` in any package and change the `package` line, or `go get <module>/go`
- `scripts/embed-check.sh` proves it in CI: for each language it creates an empty project, copies in only the single file, and builds or runs a three-line script that prints a banner.
- Each file defines a `KARTOUCHE_VERSION` constant (or variable) so an embedded copy can report which version it is.

### 9.1 License

The whole project is **MIT** licensed. Single `LICENSE` file at the repo root (no dual licensing, no Apache file). The copyright line is `Copyright (c) 2026 Kamilo Melo`; use exactly that text in `LICENSE` and in every file header.

```
MIT License

Copyright (c) 2026 Kamilo Melo

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

Where the license must also be declared:

| Place | Declaration |
|---|---|
| `rust/Cargo.toml` | `license = "MIT"` |
| `python/pyproject.toml` | `license = "MIT"` and `license-files = ["../LICENSE"]` (or a copy of the file inside `python/` so it ships in the sdist) |
| `go/` | a copy of `LICENSE` inside `go/` so the Go module proxy and pkg.go.dev recognize it |
| `bash/kartouche.sh`, `rust/src/lib.rs`, `python/src/kartouche.py`, `go/kartouche.go` | the header comment from section 9 (`SPDX-License-Identifier: MIT` plus the copyright line) |
| `README.md` | a short "License: MIT" line at the bottom |

`scripts/embed-check.sh` also verifies that each single file still contains the `SPDX-License-Identifier: MIT` line.

---

## 10. Public API

Names differ in style per language but map one to one.

### 10.1 Bash (`bash/kartouche.sh`)

Targets **bash 5.0 or newer** (the Arch Linux default), so modern bash features (associative arrays, `${var,,}`, `mapfile`, `printf -v`) are allowed. Single file, no side effects on `source` (no output, no `set` changes, no traps until `kt_steps_start`). All public functions are prefixed `kt_`; internals `_kt_`. Functions **always return 0** and never abort a script under `set -e` or `set -u`.

```bash
source "$(dirname "$0")/kartouche.sh"

kt_banner NAME VERSION DESCRIPTION [ACCENT]        # section 4.1
kt_card   NAME VERSION DESCRIPTION [ACCENT]        # section 4.2
kt_status "12 hosts · 11 up · 1 down"              # prints one ├─ line (4.1)

kt_steps_start STYLE TOTAL [LABEL ...]             # STYLE: inline|dots|crumbs|rail|line|log
kt_step "building images"                          # advance to the next step
kt_steps_done                                      # end state ✓
kt_steps_fail "reason"                             # end state ✗

# configuration (all optional, set before the call)
KT_STREAM=stderr|stdout        # default stderr
KT_APP_ACCENT_VAR=FLEETTY_ACCENT   # app-specific override name (6.4)
KT_WIDTH=100                   # terminal width override (8.4)
```

### 10.2 Rust (crate `kartouche`)

Targets the current Arch Rust; no MSRV promise (uses `std::io::IsTerminal`). `#![forbid(unsafe_code)]`. Rendering is separate from IO so everything is testable against a `String`.

```rust
use kartouche::{Banner, Steps, Style, Accent};

// reads CARGO_PKG_NAME / CARGO_PKG_VERSION / CARGO_PKG_DESCRIPTION at the caller's compile time
kartouche::banner!();

Banner::new("fleetty")
    .version("0.2.0")
    .description("your fleet, one terminal")
    .accent("magenta".parse::<Accent>().unwrap_or_default())
    .status("12 hosts · 11 up · 1 down")
    .accent_env("FLEETTY_ACCENT")
    .print();                                   // stderr by default

let s: String = Banner::new("x").version("1").render();   // pure, no IO, uses Env::detect()
let s = Banner::new("x").render_with(&Env::plain());      // deterministic, for tests

let mut steps = Steps::new(["connect", "fetch", "build", "deploy"])
    .style(Style::Crumbs)
    .start();                                   // returns a guard, prints the first region
steps.next("connect");
steps.next("fetch");
steps.done();                                   // or steps.fail("reason")
```

Requirements:

- `Accent: FromStr + Default` (default cyan). Variants: `Ansi(u8)` for the 16 names (the SGR code) and `Indexed(u8)` for `0-255`. No RGB variant in 0.1.
- `Env::detect()` snapshots environment and TTY state; `Env::plain()` and `Env::full_color()` exist for tests.
- `Stream::{Stderr, Stdout}` selectable with `.stream(...)`.
- `Steps` implements `Drop` as described in 7.3.
- No features, no dependencies. `src/lib.rs` is one self-contained file, so it can also be dropped into a project as `mod kartouche;`.
- `.width(n)` overrides terminal width (8.4).
- Publishable to crates.io with docs.rs examples that compile (`cargo test --doc`).

### 10.3 Python (package `kartouche`)

Targets the current Arch Python (3.14), no support for older versions, stdlib only, one module `kartouche.py`, fully typed (inline annotations, and a `py.typed` marker when packaged), passes `mypy --strict`.

```python
from kartouche import banner, card, status, Steps

banner("fleetty", "0.2.0", "your fleet, one terminal", accent="magenta")
banner(description="...")          # name/version default from importlib.metadata of the caller's package when available
card("lumen", "0.3.1", "tiny notes, zero friction")
status("12 hosts · 11 up · 1 down")

with Steps(["connect", "fetch", "build", "deploy"], style="crumbs") as s:
    s.next("connect")
    s.next("fetch")
    s.next("build")
# normal exit -> ✓ done; exception -> ✗ failed at n/N label (and the exception is re-raised)

banner(..., stream=sys.stdout, accent_env="FLEETTY_ACCENT")
```

Also: `render_banner(...) -> str` and `render_card(...) -> str` pure functions for tests. Use: `pip install kartouche` or copy `kartouche.py`; the demo runs via `python -m kartouche`. Width override: `width=n`.

### 10.4 Go (`go/kartouche.go`)

Targets the current Arch Go, no support for older versions, standard library only, single file, `package kartouche`. No goroutines left running after a call returns. TTY detection uses `(*os.File).Stat()` and `os.ModeCharDevice`.

```go
import "example.com/kartouche/go"   // or copy kartouche.go into your own package

kartouche.Banner("fleetty", "0.2.0", "your fleet, one terminal").
    Accent("magenta").
    AccentEnv("FLEETTY_ACCENT").
    Status("12 hosts · 11 up · 1 down").
    Print()                                  // stderr by default

out := kartouche.Banner("x", "1", "y").Render()          // pure string, uses DetectEnv()
out = kartouche.Banner("x", "1", "y").RenderWith(kartouche.PlainEnv())   // deterministic, for tests

s := kartouche.NewSteps(kartouche.Crumbs, "connect", "fetch", "build", "deploy").Start()
defer s.Close()          // restores the cursor; prints ✗ interrupted if not finished
s.Next("connect")
s.Next("fetch")
s.Done()                 // or s.Fail("reason")
```

Requirements:

- Styles are typed constants: `Inline Dots Crumbs Rail Line Log`.
- `Accent` accepts the same strings as everywhere (names and `0-255` indexes) and never returns an error for a bad value; it falls back to cyan.
- `.Stream(os.Stdout)` selects the stream, `.Width(n)` overrides width, `.Card()` renders the card instead of the banner.
- Safe to call from one goroutine at a time; document that `Steps` is not concurrency-safe.
- `go vet` and `staticcheck` clean; `go test ./...` runs expected-output tests; examples are `Example` functions so `go doc` shows them.

### 10.5 Consistency rules across languages

- Same style names: `inline dots crumbs rail line log`.
- Same environment variables and precedence (section 6).
- Same defaults: stream stderr, accent cyan, bar width 20, dots cap 12.
- Same width override and same width resolution order (8.4).
- Same behavior when arguments are empty or missing (section 4).

---

## 11. Demo CLI contract and conformance

Each language ships a tiny demo program with **this exact interface**, so one script can test all four:

```
kt-demo banner  NAME VERSION DESCRIPTION [ACCENT]
kt-demo card    NAME VERSION DESCRIPTION [ACCENT]
kt-demo steps   STYLE CURRENT TOTAL LABEL [LABEL ...]   # renders the final frame for step CURRENT of TOTAL
kt-demo end     done|fail STYLE CURRENT TOTAL LABEL     # renders the end state
kt-demo colors                                          # one banner per accent name (all 16)
```

- Demos write to **stdout** (override the default stream) so the script can capture output.
- `steps` renders statically without cursor movement (the 7.2 layout for the single frame) so expected files are stable.
- Bash: `bash/kt-demo`. Rust: `cargo run --example kt_demo --`. Python: `python -m kartouche`. Go: `go run ./go/cmd/kt-demo`.

`spec/cases.txt` has one line per case: `<expected path> | <env assignments> | <demo args>`. Example:

```
plain/banner_basic.txt   | NO_COLOR=1               | banner app_name 0.1.0 app_description
ansi/banner_magenta.ansi | FORCE_COLOR=1 TERM=xterm-256color | banner app_name 0.1.0 app_description magenta
ascii/banner_basic.txt   | NO_COLOR=1 KARTOUCHE_ASCII=1 | banner app_name 0.1.0 app_description
plain/steps_dots_2_5.txt | NO_COLOR=1               | steps dots 2 5 connect fetch "building images" push deploy
plain/card_basic.txt     | NO_COLOR=1               | card lumen 0.3.1 "tiny notes, zero friction"
```

`scripts/conformance.sh` runs every case through all four demos and diffs against the expected file. A single failing byte fails CI.

### Required test cases (minimum)

- banner: basic; no description; no version; version already prefixed with `v`; long description truncated at `COLUMNS=30`; status line; each of the 16 accent names (ansi); an index on a 256-capable terminal (`TERM=xterm-256color`); the same index on a non-256 terminal (`TERM=xterm`, must fall back to cyan); a hex value and an out-of-range index (both fall back to cyan).
- card: basic; short description (name wider than description); long name.
- steps: every style at start, middle and last step; rail and log with failure; unknown total; end states `done` and `fail`.
- environment: `NO_COLOR` beats `FORCE_COLOR`; `KARTOUCHE_ASCII=1`; `LANG=C` auto-ASCII; `TERM=dumb`.
- width: CJK description in a card.

---

## 12. Quality requirements

### Behavior

- Safe in pipes: `app | cat` and `app > file` contain no escape sequences unless `FORCE_COLOR` is set.
- Idempotent cleanup: cursor is always restored, including after `Ctrl-C`.
- No output on import or `source`.
- Never throw, panic, or call `exit` because of formatting input.
- Handle inputs containing `%`, backslashes, `$` and newlines literally (Bash must use `printf '%s'`, never interpolate user text into a format string). Newlines in fields are replaced by a single space.

### Testing

- Bash: `bats` plus `shellcheck` with no warnings, tested on the current bash 5.x.
- Rust: unit tests for color parsing and width; expected-output tests; `cargo clippy -- -D warnings`; `cargo fmt --check`.
- Python: `pytest`, `ruff`, `mypy --strict`, tested on the current Arch version.
- Go: `go test ./...`, `go vet`, `staticcheck`, tested on the current Arch version.
- CI: Linux only, with the primary job running in an `archlinux:latest` container (the reference platform) and a secondary Ubuntu job; run `scripts/conformance.sh` and `scripts/embed-check.sh`.

### Documentation

- README with: the banner and all six step styles as text blocks, a quickstart for each language (under 10 lines each), the environment variable table, and a "how to embed the single file" note for each language.
- Every public item documented.
- `examples/` has one realistic script per language that prints a banner, runs five fake steps with `sleep`, and ends with `done`.

---

## 13. Build phases and acceptance criteria

**Phase 1: shared foundations**
- [ ] Repo skeleton, `LICENSE` (MIT, section 9.1), CI with empty conformance step.
- [ ] `spec/expected` files for banner and card (plain, ansi, ascii), `spec/cases.txt`, `scripts/conformance.sh`.

**Phase 2: banner and card in Bash, Rust and Python**
- [ ] Bash, Rust, Python render banner and card matching the expected files byte for byte.
- [ ] Environment rules (section 6) fully implemented and covered by expected files.
- [ ] Accent parsing: 16 names and `0-255` indexes, with the cyan fallback of 5.3.

**Phase 3: steps**
- [ ] Six styles plus end states in Bash, Rust and Python, static frames matching the expected files.
- [ ] Redraw logic for TTY, append-only for non-TTY and `CI`.
- [ ] Cursor restore verified by a test that sends `SIGINT` to an example (or an equivalent harness).

**Phase 4: Go port** (done after the other three pass conformance, so spec mistakes are fixed once)
- [ ] `go/kartouche.go` renders banner, card, all six step styles and end states matching every expected file byte for byte.
- [ ] Cursor-restore and signal handling verified (section 7.3).
- [ ] `go vet`, `staticcheck` and `go test ./...` clean; `Example` functions present.

**Phase 5: embedding, packaging and polish**
- [ ] `scripts/embed-check.sh` passes: each single file works alone in an empty project, in all four languages.
- [ ] Each file has the header comment (including the MIT SPDX line and copyright) and `KARTOUCHE_VERSION` from section 9.
- [ ] License declared in `Cargo.toml`, `pyproject.toml` and `go/` as listed in section 9.1.
- [ ] `cargo publish --dry-run` passes; docs.rs examples compile.
- [ ] `python -m build` and `twine check` pass.
- [ ] Go: module tagged as `go/v0.1.0` (a module in a subdirectory needs the prefixed tag); `go get` works from a clean machine.
- [ ] Bash: an `install.sh` that copies `kartouche.sh` to `~/.local/share/kartouche/`.
- [ ] README complete; `CHANGELOG.md` has 0.1.0.

**Done means:** `scripts/conformance.sh` is green in CI, every checkbox above is ticked, and the four example scripts look identical in a terminal.

---

## 14. Roadmap (do not build in 0.1)

- Hex and truecolor accents, with `COLORTERM` detection and a nearest-color fallback (deliberately cut from 0.1: terminal themes remap the 16 colors, so any conversion is a guess).
- `kt_spinner` / `Spinner`: a single-line animated indicator inside the steps region.
- `kt_confirm`-style prompts in the same frame language.
- A Node port using the same expected files (single file, no dependencies).
