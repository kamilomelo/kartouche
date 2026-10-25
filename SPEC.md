# kartouche: a design system of copy-paste templates for my CLI scripts

> Working name: **kartouche**. Renaming is a single find-and-replace.
> This is the **active plan**. The heavier library/package plan is archived in `SPEC-library.md` as a future upgrade path. Do not build that one.

## 0. Instructions for the implementing agent

1. Read this whole document first.
2. Build the repo in the phases of section 9, in order.
3. **This is not a library.** Nothing here is published, versioned for consumers, or imported by other repos. Every file is a template that gets **copied** into another project and edited there.
4. **No dependencies**, standard library only, in every language.
5. **Keep the code short and boring.** The size target is about 150 lines per template file (Rust and Go about 200). This is a **goal, not an acceptance criterion**: a file may exceed it if correctness or the expected output requires, in which case say so in the final report. The hard rules are: no dependencies, one self-contained file per language, the markers of section 5.1, and the expected output. If a feature is not in this document, do not add it.
6. The file `expected/demo.txt` is the source of truth for how every template must look. All four languages must reproduce it byte for byte when color is off.
7. Reference platform: **Arch Linux**, bash 5.x, and the **current Arch versions** of Python (3.14), Rust and Go. Do not add compatibility code for other platforms or for older language versions; modern syntax and standard-library features are fine.
8. The repo already exists at **https://github.com/kamilomelo/kartouche** (see section 3.1 for the remote, SSH key and what to commit).

---

## 1. Idea

I write many small CLI scripts in Bash, Rust, Python and Go. I want them all to look like one family: the same banner, the same card, the same step counter. This repo is the **design system** for that: a written look (`DESIGN.md` content is section 4) plus one small template file per language that implements it.

To use it: copy `kartouche.<ext>` into a script, edit the config block at the top, call the functions. The copy belongs to that script. If I later improve kartouche, I ask my coding agent to pull the new template into the other scripts (section 6).

### Principles

- **Copy, then own.** No linking between repos, no install step, no version management.
- **Update by agent.** Files are structured so an agent can replace the shared code and keep my edits (section 5.1).
- **One accent color.** The design has exactly one real color; everything else follows the terminal theme (section 4.6).
- **Safe in pipes and logs.** Never emit escape codes when output is not a terminal.

## 2. Non-goals

- No package publishing (crates.io, PyPI, Go module proxy), no version tags for consumers.
- No dependencies, optional or otherwise. No feature flags.
- No CI matrix, no cross-platform support beyond Arch Linux.
- No Unicode width handling. The card assumes every character is one column wide; text with wide characters (CJK, emoji) will misalign its right border, and keeping the text single-width is the caller's responsibility. The banner and steps are left-aligned, so they are unaffected. No terminal-width truncation either.
- No ASCII fallback beyond the single `TERM=linux` / `TERM=dumb` rule in section 4.7. Otherwise UTF-8 terminals are assumed.
- No color-name parser or color-depth detection: the accent is a raw constant.
- No spinners, prompts, themes or config files.

---

## 3. Repo layout

```
kartouche/
├── README.md          # screenshots as text blocks, 4 quickstarts, "how to update a script"
├── DESIGN.md          # the look (copy of section 4 of this file)
├── UPDATING.md        # instructions for a coding agent to refresh a script (section 6)
├── LICENSE            # MIT (section 8)
├── bash/kartouche.sh
├── rust/kartouche.rs
├── python/kartouche.py
├── go/kartouche.go
├── demo/
│   ├── demo.sh
│   ├── demo.rs
│   ├── demo.py
│   ├── demo.go
│   ├── steps.sh        # live Steps run: `steps.<ext> STYLE` (used by preview and phase 4 checks)
│   ├── steps.rs
│   ├── steps.py
│   └── steps.go
├── expected/
│   ├── demo.txt        # exact output of every demo with NO_COLOR=1 and a Unicode TERM
│   └── demo-ascii.txt  # same output with TERM=linux (ASCII glyphs)
├── check.sh           # runs all four demos, diffs against the expected files
└── preview.sh         # runs the demos with color so I can look at them
```

One file per language contains **everything** (banner, card, steps). Unused functions are harmless, so copies never need pruning, which keeps updates a clean replace.

### 3.1 GitHub repo and git setup

- **Repo:** https://github.com/kamilomelo/kartouche (owner `kamilomelo`, MIT, default branch `main`).
- **Remote (SSH):** `git@github.com:kamilomelo/kartouche.git`
- **SSH key:** the key `~/.ssh/id_github-km` is already authorized for this repo. Use it without creating or modifying keys. If the default SSH config does not select it, either add this to `~/.ssh/config` (and ask before editing that file) or use it for a single command:
  ```
  Host github.com
      IdentityFile ~/.ssh/id_github-km
      IdentitiesOnly yes
  ```
  ```
  GIT_SSH_COMMAND='ssh -i ~/.ssh/id_github-km -o IdentitiesOnly=yes' git push
  ```
- **Setup (initialize in place; do not clone into this directory):** the working directory already contains `SPEC.md` and `SPEC-library.md` and is **not** a git repo, so a clone into it would fail. Instead, from the working directory:
  1. `git init -b main`
  2. `git remote add origin git@github.com:kamilomelo/kartouche.git`
  3. `git fetch origin` (with the SSH key as above).
  4. If `origin/main` does not exist (empty remote), continue to the first commit.
  5. If `origin/main` exists (for example a GitHub-generated README or LICENSE), inspect its tracked paths and contents (`git ls-tree -r --name-only origin/main` and, where needed, `git show origin/main:<path>`) against the local files **before** resetting or restoring anything. If a remote path overlaps a local path other than the planned `README.md` and `LICENSE`, stop and ask me how to reconcile it; do not overwrite or discard local work. Then run `git reset origin/main` (the default mixed reset, **not** `--hard`) so history continues from the remote while working files stay untouched. Our `README.md` and `LICENSE` replace any remote versions; restore other remote files **only if the path is absent locally**. **Never force-push.**
  6. Keep `SPEC.md` and `SPEC-library.md` where they are and include them in the Phase 1 commit.
  Never commit keys or any file from `~/.ssh`.
- **Commits:** one commit at the end of each phase in section 9, with a short imperative message (`Add bash and python templates`). **Do not push** unless I ask.
- **GitHub description** (one line, set by me in the repo settings): *Copy-paste terminal banners, cards and step counters for Bash, Rust, Python and Go: one look for all my CLI scripts.*
- **GitHub topics:** `cli`, `tui`, `terminal`, `bash`, `rust`, `python`, `golang`, `ansi`, `design-system`, `templates`.
- **Raw file URLs** (used by `UPDATING.md` and the README): `https://raw.githubusercontent.com/kamilomelo/kartouche/main/<language folder>/kartouche.<ext>`, for example `https://raw.githubusercontent.com/kamilomelo/kartouche/main/bash/kartouche.sh`.

---

## 4. Design (the look)

Monospace, UTF-8. `·` in this section marks nothing unless inside literal output.

### 4.1 Banner

```
╭─ ● app_name  v0.1.0
╰─❯ app_description
```

- Line 1: `╭─`, space, accent `●`, space, **name** (bold), two spaces, `v` + version (dim).
- Line 2: `╰─❯`, space, description (default color).
- Frame glyphs (`╭─`, `╰─❯`, `├─`) are dim.
- If the version already starts with `v`, do not add another. If the version is empty, omit it and its two spaces. If the description is empty, omit line 2.
- **Status** is an optional fourth input of the banner (`banner(status)`), not a separate call. When given and non-empty it is printed between line 1 and the description line, as a dim `├─` followed by the text in the default color:
  ```
  ╭─ ● fleetty  v0.2.0
  ├─ 12 hosts · 11 up · 1 down
  ╰─❯ your fleet, one terminal
  ```
  If the description is empty, the status becomes the last line and uses `╰─❯` instead of `├─`. To show a changed status later, print the banner again.

### 4.2 Card

```
╭─ lumen ───────────  v0.3.1 ─╮
│  tiny notes, zero friction  │
╰─────────────────────────────╯
```

- Width `w = max(len(description) + 4, left + right + 2)` where `left = len("─ " + name + " ")` and `right = len("  v" + version + " ─")`. Lengths are character counts, not display widths: the card is the only element with a right edge, and it is the caller's job to pass single-width text.
- Top: `╭`, `─ `, **name** (bold), space, `─` × `(w - left - right)`, two spaces, `v` + version, ` ─`, `╮`.
- Body: `│`, two spaces, description padded to `w - 4`, two spaces, `│`.
- Bottom: `╰`, `─` × `w`, `╯`. Frame is dim.

### 4.3 Steps

Steps replace the banner's description line. While steps run, the output is: the **title line** (`╭─ ● name  vX`, printed once) followed by a **region** that changes. `n` is the current step (1-based), `N` the total. Frames below are for `n=3`, `N=4`, labels `connect fetch build deploy`, current label `build`.

**inline** (default)
```
╭─ ● app_name  v0.1.0
╰─❯ 3/4 build
```

**dots**: `●` × n (accent), `○` × (N − n) (dim), two spaces, counter (dim). Above 12 steps, use `line`.
```
╭─ ● app_name  v0.1.0
├─ ●●●○  3/4
╰─❯ build
```

**crumbs**: all labels joined by ` › `. Done and pending labels dim, current label bold accent.
```
╭─ ● app_name  v0.1.0
├─ connect › fetch › build › deploy
╰─❯ 3/4 build
```

**rail**: one line per step. `✓` (accent) done, `◉` (accent, bold label) current, `○` (dim) pending.
```
╭─ ● app_name  v0.1.0
├─ ✓ connect
├─ ✓ fetch
├─ ◉ build
├─ ○ deploy
╰─❯ 3/4
```

**line**: bar of 20 cells. `filled = round(20 * n / N)`. If `filled == 0`, 20 × `─`. If `filled >= 20`, 20 × `━`. Otherwise `━` × (filled − 1), one `╸`, then `─` × (20 − filled). Filled part accent, rest dim, then two spaces and the counter.
```
╭─ ● app_name  v0.1.0
├─ ━━━━━━━━━━━━━━╸─────  3/4
╰─❯ build
```

**log**: finished steps stay as lines `├─ i/N label` + space + `·` × `(longest_label - len(label) + 7)` + ` ok`. The last line is the active step. `ok` is accent.
```
╭─ ● app_name  v0.1.0
├─ 1/4 connect ······· ok
├─ 2/4 fetch ········· ok
╰─❯ 3/4 build
```

Rounding for `line` is round-half-up (`floor(x + 0.5)`).

**The all-done frame.** For `N > 0`, `render_frame` accepts `n` from `1` to `N + 1`. `n = N + 1` is the "all done" frame that `done` uses (section 5.5): every step is shown as done and none is current, the counter reads `N/N`, and the `label` argument is ignored. `dots` shows `N` filled dots, `line` a full bar of 20 `━`, `crumbs` has no highlighted label (all labels dim), `rail` shows `✓` for every step, and `log` shows `N` lines ending in `ok`. Every style ends with the same replaceable final line `╰─❯ N/N`; `done` replaces that line with the end line (section 4.4). For `N = 0` there is no frame; see section 5.5.

### 4.4 End states

Replace the last line of the region:

```
╰─ ✓ done in 4.2s
╰─ ✗ failed at 3/4 build
```

`✓` is accent. `✗` is red. Elapsed time is optional text supplied by the caller or computed by the template as seconds with one decimal (`4.2s`).

### 4.5 Spacing and text rules

- Always two spaces between name and version, two spaces before a counter.
- Newlines in any text field are replaced by a single space.
- Text is printed literally. Never interpret `%`, `\`, `$` or backslash sequences in user text (Bash: use `printf '%s'`).

### 4.6 Color

| Role | Style |
|---|---|
| frame glyphs, version, counter, pending items | dim `\e[2m` |
| name | bold `\e[1m` |
| accent: `●`, `✓`, `◉`, filled bar, `ok`, current crumb | the config constant `ACCENT` |
| failure: `✗` | red `\e[31m` |
| every styled run ends with | reset `\e[0m` |

- `ACCENT` is a raw SGR parameter string chosen by hand when the template is copied: `36` (cyan, default), `95` (bright magenta), `38;5;141` (256-color purple), and so on. There is no name parser, no capability detection and no fallback logic.
- Approved accent suggestions to list in `README.md` and `DESIGN.md`:

| Look | `ACCENT` |
|---|---|
| cyan (default) | `36` |
| magenta | `35` |
| yellow / amber | `33` |
| green | `32` |
| blue | `34` |
| soft purple | `38;5;141` |
| sky blue | `38;5;81` |
| warm orange | `38;5;214` |
| mint | `38;5;78` |
| coral | `38;5;203` |

- Readability note for the docs: `30` (black) and `90` (gray) vanish on some themes, so avoid them as accents.
- One accent per app is the intended use. The frame, glyphs and layout stay identical across apps.

### 4.7 Console fallback and append-only mode (for installers)

Some scripts run on the bare Linux console (for example an Arch install from the live ISO, where `TERM=linux`), where the font may lack the rounded corners and symbols, and where other programs (`pacstrap`, `arch-chroot`) print between steps and would corrupt an in-place redraw.

**ASCII glyphs.** The config constant `GLYPHS` (default `auto`) picks the glyph set: `unicode` always uses the Unicode glyphs, `ascii` always uses the ASCII twins, and `auto` uses ASCII when `TERM` is `linux` or `dumb` and Unicode otherwise. Layout, spacing, colors and line counts stay identical.

`TERM` is read **every time something is rendered**, never cached. If a script changes `TERM` midway (or a nested program does), the new glyph set applies from the next line on. A banner or region is never half one set and half the other, because each render call reads `TERM` once and uses it for the whole call. Scripts that know where they run should skip the guessing: an installer meant for the Linux console, including inside `tmux` or `screen` on it (where `TERM` becomes `screen` or `tmux-256color` and `auto` would wrongly pick Unicode), sets `GLYPHS="ascii"`.

| Unicode | ASCII | Unicode | ASCII | Unicode | ASCII |
|---|---|---|---|---|---|
| `╭` | `,` | `╮` | `.` | `●` | `*` |
| `╰` | `` ` `` | `╯` | `'` | `○` | `o` |
| `─` | `-` | `│` | `\|` | `◉` | `@` |
| `├` | `\|` | `❯` | `>` | `✓` | `+` |
| `›` | `>` | `━` | `=` | `✗` | `x` |
| `╸` | `=` | `·` | `.` | | |

**Append-only steps.** The config constant `APPEND_ONLY` (default false) forces the non-redraw behavior of section 5.3 rule 4 even on a terminal: the title line is printed once, then one new line per step, never moving the cursor. Installers should set it to true.

**Colors on the console.** The Linux console only approximates 256-color values, so scripts meant for it should use a plain accent (`30` to `37` or `90` to `97`; avoid `30` and `90`).

---

## 5. Template contract

### 5.1 File structure (the key to easy updates)

Every `kartouche.<ext>` file has exactly this shape:

```
<header comment: name, `SPDX-License-Identifier: MIT`, `Copyright (c) 2026 Kamilo Melo`, the repo URL https://github.com/kamilomelo/kartouche, "copy this file into your project", KARTOUCHE_TEMPLATE_VERSION>

<<< CONFIG: edit these (an updater must keep this region) >>>
NAME, VERSION, DESCRIPTION, ACCENT, APPEND_ONLY, GLYPHS   (and nothing else)
<<< END CONFIG >>>

<<< KARTOUCHE CODE v<N>: do not edit (an updater replaces this region) >>>
...all functions...
<<< END KARTOUCHE CODE >>>
```

- Markers are real comment lines, identical in every language apart from the comment prefix (`#`, `//`). Use these exact texts:
  - `kartouche:config begin` / `kartouche:config end`
  - `kartouche:code begin v1` / `kartouche:code end`
- The number after `v` in the code marker is the template revision. Bump it whenever the code region changes. The header repeats it as `KARTOUCHE_TEMPLATE_VERSION`; an updater must keep these two values in sync (section 6).
- Config values are plain constants. The code region only reads them.
- **Go syntax exception:** `package main` and a standard-library `import (...)` block go between the header comment and the config region; Go requires imports before constant declarations. The imports are not part of the editable config. An updater must reconcile that import block when a new Go code revision adds or removes imports (section 6).
- An updater can therefore replace everything between the code markers, update only the header's `KARTOUCHE_TEMPLATE_VERSION` value to match, and leave the config region and the rest of the host script untouched.

### 5.2 Functions per language

Same names and behavior everywhere, in each language's casing.

**Two layers.**

- **Render functions** return text and write nothing. They take explicit values (name, version, description, ...) and a **destination stream**, which is used *only* to decide color and glyphs for text meant for that stream (section 5.3). They read `ACCENT` and `GLYPHS` from the config region and `TERM` and `NO_COLOR` from the environment, but never `NAME`, `VERSION` or `DESCRIPTION`. The print functions, `Steps`, the demos and tests are all built on them. Returned text has lines joined by `\n` and **no trailing newline** (Bash prints it followed by one newline, because shell functions return text on stdout).
- **Print functions** and **`Steps`** read `NAME`, `VERSION`, `DESCRIPTION` (and `ACCENT`, `APPEND_ONLY`, `GLYPHS`) from the config region and always write to **stderr**. A script that needs stdout, or a different name, calls the render functions and writes the text itself.

| Function | What it does |
|---|---|
| `render_banner` | banner: 2 lines, or 3 with a status line (4.1) |
| `render_card` | card: 3 lines (4.2) |
| `render_title` | the title line only: `╭─ ● name  vX` |
| `render_frame` | the steps region for step `n` of `N` (4.3, including the all-done frame `n = N + 1`) |
| `render_end` | the end line: `╰─ ✓ text` when ok, `╰─ ✗ text` when not (4.4); the caller supplies the whole text, such as `done in 4.2s` |
| `banner` / `card` | render with the config values, then print |
| `Steps` | the lifecycle of section 5.5 |

**Bash** (`FD` is `1` or `2`; color is decided with `[[ -t FD ]]`; `STYLE` is one of `inline dots crumbs rail line log`)
```bash
kt_render_banner FD NAME VERSION DESCRIPTION [STATUS]
kt_render_card   FD NAME VERSION DESCRIPTION
kt_render_title  FD NAME VERSION
kt_render_frame  FD STYLE N TOTAL LABEL LABELS...
kt_render_end    FD OK TEXT              # OK is 1 (✓) or 0 (✗)

kt_banner [STATUS]
kt_card
kt_steps_start STYLE LABELS...
kt_steps_next
kt_steps_done [ELAPSED_TEXT]
kt_steps_fail [REASON]
```

**Rust** (`Stream` decides the TTY check; `Style` is an enum)
```rust
pub enum Stream { Stdout, Stderr }
pub enum Style { Inline, Dots, Crumbs, Rail, Line, Log }

pub fn render_banner(name: &str, version: &str, description: &str, status: Option<&str>, stream: Stream) -> String
pub fn render_card(name: &str, version: &str, description: &str, stream: Stream) -> String
pub fn render_title(name: &str, version: &str, stream: Stream) -> String
pub fn render_frame(style: Style, n: usize, total: usize, label: &str, labels: &[&str], stream: Stream) -> String
pub fn render_end(ok: bool, text: &str, stream: Stream) -> String

pub fn banner(status: Option<&str>)
pub fn card()
impl Steps {
    pub fn start(style: Style, labels: &[&str]) -> Steps
    pub fn next(&mut self)
    pub fn done(&mut self, elapsed: Option<&str>)
    pub fn fail(&mut self, reason: Option<&str>)
}
```

**Python** (`stream` is a text file object such as `sys.stdout`; `style` is a string)
```python
def render_banner(name, version, description, status, stream) -> str
def render_card(name, version, description, stream) -> str
def render_title(name, version, stream) -> str
def render_frame(style, n, total, label, labels, stream) -> str
def render_end(ok, text, stream) -> str

def banner(status: str | None = None) -> None
def card() -> None
class Steps:
    def __init__(self, style: str, labels: list[str]) -> None   # starts immediately
    def next(self) -> None
    def done(self, elapsed: str | None = None) -> None
    def fail(self, reason: str | None = None) -> None
```

**Go** (`f` is `os.Stdout` or `os.Stderr`)
```go
type Style int // Inline, Dots, Crumbs, Rail, Line, Log

func RenderBanner(name, version, description, status string, f *os.File) string
func RenderCard(name, version, description string, f *os.File) string
func RenderTitle(name, version string, f *os.File) string
func RenderFrame(style Style, n, total int, label string, labels []string, f *os.File) string
func RenderEnd(ok bool, text string, f *os.File) string

func Banner(status string)                           // "" means no status line
func Card()
func NewSteps(style Style, labels ...string) *Steps  // starts immediately
func (s *Steps) Next()
func (s *Steps) Done(elapsed string)                 // "" means measured seconds
func (s *Steps) Fail(reason string)                  // "" means no reason
func (s *Steps) Close()
```

Rules shared by all languages:

- Style names: `inline dots crumbs rail line log`.
- Python: `Steps` also works as a context manager (`with Steps(...) as s:`); leaving normally calls `done()`, leaving with an exception calls `fail()` and re-raises.
- Rust: `Drop` on `Steps`, and Go: `Close()`, finish an unfinished run with `╰─ ✗ interrupted at n/N label`. Neither runs on a default `SIGINT`, and the templates install **no signal handlers** (5.5).
- Functions must never abort the host program: in Bash they return 0 even under `set -e` and `set -u`; in Rust, Python and Go they ignore write errors.

### 5.3 Runtime rules (small and fixed)

1. **Stream:** the print functions and `Steps` write to **stderr**, so decoration never pollutes piped stdout. Render functions write nothing; the stream they receive only drives the color and glyph decisions of rule 2.
2. **Color decision:** color is on only if the destination stream is a terminal **and** `NO_COLOR` is unset or empty. Otherwise no color or style sequence is ever emitted, not even a reset. (The cursor sequences of the redraw in rule 3 are not colors; they are emitted on a terminal even when `NO_COLOR` is set, and never when output is not a terminal.)
   **Glyph set:** per the `GLYPHS` constant of section 4.7 (in `auto`, ASCII when `TERM` is `linux` or `dumb`), reading `TERM` on every render. No other environment variable is read.
3. **Redraw (terminal only, and only when `APPEND_ONLY` is false):** `start` prints the title line once, then the region for step 1. On each `next`, move the cursor up by the number of lines in the previous region (`\e[<k>A`), then for each line write `\r`, clear it (`\e[2K`) and print the new line. After `done` or `fail`, redraw with the end state. No cursor hiding.
4. **Not a terminal, or `APPEND_ONLY` is true:** no cursor movement, no clearing. Print the title line once, then append one line per step when it starts: `├─ n/N label`, and finally `╰─ ✓ done in 4.2s` or `╰─ ✗ failed at n/N label`. This keeps CI logs readable.
5. **Elapsed time:** see 5.5.
6. **Safety:** user text is never used as a format string; fields are printed with `%s`-style literal output.

### 5.4 Per-language notes

- **Bash**: bash 5.x only. Single file meant to be `source`d; sourcing prints nothing and changes no shell options. Public functions are prefixed `kt_`, internals `_kt_`. Config variables are `KT_NAME`, `KT_VERSION`, `KT_DESCRIPTION`, `KT_ACCENT`, `KT_APPEND_ONLY`, `KT_GLYPHS`.
- **Rust**: usable as `mod kartouche;` or pasted into a single-file script. Config are `const NAME`, `VERSION`, `DESCRIPTION`, `ACCENT: &str`. Use `std::io::IsTerminal`. No `unsafe`, no crates. Public items are `pub`; the file starts with `#![allow(dead_code)]` after the header comment so unused functions do not warn. It defines **no `fn main`**.
- **Python**: current Arch Python (3.14); modern syntax is fine (`X | None`, `match`, and so on). Stdlib only, type hints. Config are module constants. Terminal check with `sys.stderr.isatty()`.
- **Go**: current Arch Go; modern features (generics, `slices`, `min`/`max` builtins) are fine. `package main` by default (the header says to change it when pasted into another package). On Linux, check whether the destination file is a terminal using a `TCGETS` ioctl on its file descriptor (standard-library `syscall` and `unsafe` are permitted for this check); a successful ioctl means terminal, and an error means not a terminal. `(*os.File).Stat()` with `os.ModeCharDevice` alone is **not** sufficient: `/dev/null` is a character device but not a terminal. The file defines **no `func main`** so a demo or host program can supply its own. `Steps` is not concurrency-safe and the doc comment says so.

### 5.5 Steps lifecycle

State: `labels` (the plan; `N` is its length), `n` (current step, 1-based), a `finished` flag, the start time, and the number of region lines currently on screen.

1. **start** (`kt_steps_start`, `Steps::start`, `Steps(...)`, `NewSteps`): record the start time and set `n = 1` when `N > 0` (or `n = 0` when `N = 0`). Print the title line, then **show step 1 immediately** when `N > 0`: in redraw mode the region `render_frame(style, 1, N, labels[0], labels)`; in append-only mode the line `├─ 1/N <labels[0]>`. An empty `labels` is allowed: only the title line is printed, no region is displayed, and `next` does nothing.
2. **next** takes **no arguments**. If `finished`, do nothing. If `n < N`, set `n = n + 1` (the increment happens first, then the new step is rendered): redraw mode replaces the previous region with the region for the new `n`; append-only mode appends `├─ n/N <labels[n-1]>`. If `n == N` (already on the last step), do nothing: `n` never exceeds `N` and nothing is printed. The label shown is always `labels[n-1]`; there is no per-call label.
3. **done(elapsed)**: if `finished`, do nothing; otherwise set `finished`. The end text is `done in <elapsed>`, where `<elapsed>` is the caller's string when one is given and otherwise the seconds since `start` with one decimal and an `s` (`4.2s`). For `N > 0`, redraw mode replaces the region with `render_frame(style, N + 1, N, "", labels)` whose **last line is replaced** by `render_end(true, text)`. For `N = 0`, there is no region to redraw: print `render_end(true, text)` on a new line after the title, without cursor movement. Append-only mode also appends `render_end(true, text)`. Calling `done` before the last step is allowed and shows every step as complete.
4. **fail(reason)**: if `finished`, do nothing; otherwise set `finished`. The end text is `failed at n/N <labels[n-1]>` using the current `n`, followed by ` · <reason>` when a non-empty reason is given. With an empty `labels` the text is just `failed` plus the optional reason. For `N > 0`, redraw mode replaces the region with `render_frame(style, n, N, labels[n-1], labels)` whose last line is replaced by `render_end(false, text)`. For `N = 0`, print `render_end(false, text)` after the title without cursor movement. Append-only mode appends `render_end(false, text)`. No elapsed time is shown on failure.
5. After `done` or `fail` the output ends with a newline and every later call on that `Steps` is ignored.
6. **Unfinished runs:** a `Steps` that is dropped (Rust `Drop`), closed (Go `Close`) or left without `done`/`fail` finishes with the end text `interrupted at n/N <labels[n-1]>` (shown with `✗`). With empty `labels`, the text is just `interrupted` and the end line is appended after the title without cursor movement. Python's context manager calls `fail()` on an exception instead. Bash has no object and does nothing.
7. **No signal handling:** the templates install no signal handlers, never hide the cursor and never change terminal modes, so `Ctrl-C` always leaves the terminal usable. The `interrupted` line appears only when the host program unwinds normally (a default `SIGINT` kills a Rust or Go process without running `Drop` or deferred calls).

---

## 6. Updating copies (content of `UPDATING.md`)

`UPDATING.md` is addressed to a coding agent and must say, in plain steps:

1. **Locate** the copy of kartouche in the target script or repo: a file named `kartouche.<ext>`, or code between the markers `kartouche:code begin` and `kartouche:code end` pasted into a larger file.
2. **Fetch** the current template for that language from this repo: a local clone path if one exists, otherwise the raw file URL `https://raw.githubusercontent.com/kamilomelo/kartouche/main/<language folder>/kartouche.<ext>` (repo: https://github.com/kamilomelo/kartouche).
3. **Compare** the revision: if `kartouche:code begin vN` has the same `N`, stop and report "already up to date".
4. **Replace** the text between `kartouche:code begin` and `kartouche:code end` (markers included) with the same region from the new template. Also update **only** the header's `KARTOUCHE_TEMPLATE_VERSION` value to match the new code marker; never leave them out of sync.
5. **Leave untouched:** the config region (`kartouche:config`), everything else outside the code region (except that one header version value), and the host script's own calls.
6. **Check what changed:** if the new revision renamed or removed a function the host script calls, list the affected call sites and fix them. If config fields were added to the new template, add them to the host script's config region with the template's default values and tell me.
   For Go, compare the new template's import block with the copy's import block and add or remove imports required by the new code; preserve any imports used elsewhere by the host script. This is the sole exception to step 5 apart from the header version.
7. **Verify:** run the host script's own test or at least a dry run, and show me the diff.

`README.md` repeats this as a one-paragraph recipe: *"To update a script, tell your agent: update the kartouche in this script from `<path or URL>`, following `UPDATING.md`."*

---

## 7. Demo, expected output and checks

### 7.1 Demo contract

Each `demo/demo.<ext>` uses **only the render functions** of the template in the matching language folder. Render functions write nothing, so the demo prints their text to **stdout** itself, each followed by one newline. This is the whole mechanism for getting stdout output; the templates have no output override.

It prints, with `NO_COLOR=1 TERM=xterm-256color`, exactly `expected/demo.txt`, and with `NO_COLOR=1 TERM=linux`, exactly `expected/demo-ascii.txt`. The destination stream passed to every render function is stdout. In order:

1. `render_banner("app_name", "0.1.0", "app_description", no status)`.
2. A blank line, then `render_card("lumen", "0.3.1", "tiny notes, zero friction")`.
3. For each style in the order `inline dots crumbs rail line log`: a blank line, `-- <style>`, `render_title("app_name", "0.1.0")`, then `render_frame(style, 3, 4, "build", ["connect","fetch","build","deploy"])`.
4. A blank line, `-- end`, `render_end(true, "done in 4.2s")`, `render_end(false, "failed at 3/4 build")`.

(Match the blank lines exactly to section 7.2.) The demos never read the config region for name, version or description.

**Run commands** (from the repo root; `$tmp` is a fresh directory from `mktemp -d`; template files define no `main`, so the demo supplies it):

| Language | Command |
|---|---|
| Bash | `bash demo/demo.sh` (the demo starts with `source "$(dirname "${BASH_SOURCE[0]}")/../bash/kartouche.sh"`) |
| Rust | `rustc --edition 2024 -o "$tmp/demo" demo/demo.rs && "$tmp/demo"` (`demo.rs` declares `#[path = "../rust/kartouche.rs"] mod kartouche;`) |
| Python | `PYTHONPATH=python python3 demo/demo.py` (the demo does `import kartouche`) |
| Go | `cp go/kartouche.go demo/demo.go "$tmp"/ && (cd "$tmp" && go run kartouche.go demo.go)` |

Go note: `go run a.go b.go` requires all files to be in one directory, so the demo is run from a temporary copy; both files are `package main` and only `demo.go` has `func main`.

**Live steps demos.** `demo/steps.<ext> STYLE` (run the same way as above, with the style as the extra argument) creates a `Steps` with the labels `connect fetch build deploy`, sleeps 0.3 s between steps, ends with `done`, and uses the real print path, so everything goes to **stderr**. `preview.sh` and the Phase 4 checks use it.

### 7.2 `expected/demo.txt`

````text
╭─ ● app_name  v0.1.0
╰─❯ app_description

╭─ lumen ───────────  v0.3.1 ─╮
│  tiny notes, zero friction  │
╰─────────────────────────────╯

-- inline
╭─ ● app_name  v0.1.0
╰─❯ 3/4 build

-- dots
╭─ ● app_name  v0.1.0
├─ ●●●○  3/4
╰─❯ build

-- crumbs
╭─ ● app_name  v0.1.0
├─ connect › fetch › build › deploy
╰─❯ 3/4 build

-- rail
╭─ ● app_name  v0.1.0
├─ ✓ connect
├─ ✓ fetch
├─ ◉ build
├─ ○ deploy
╰─❯ 3/4

-- line
╭─ ● app_name  v0.1.0
├─ ━━━━━━━━━━━━━━╸─────  3/4
╰─❯ build

-- log
╭─ ● app_name  v0.1.0
├─ 1/4 connect ······· ok
├─ 2/4 fetch ········· ok
╰─❯ 3/4 build

-- end
╰─ ✓ done in 4.2s
╰─ ✗ failed at 3/4 build
````

(A blank line separates blocks. The file ends with a single newline.)

### 7.2b `expected/demo-ascii.txt`

The same output with the ASCII glyphs of section 4.7:

````text
,- * app_name  v0.1.0
`-> app_description

,- lumen -----------  v0.3.1 -.
|  tiny notes, zero friction  |
`-----------------------------'

-- inline
,- * app_name  v0.1.0
`-> 3/4 build

-- dots
,- * app_name  v0.1.0
|- ***o  3/4
`-> build

-- crumbs
,- * app_name  v0.1.0
|- connect > fetch > build > deploy
`-> 3/4 build

-- rail
,- * app_name  v0.1.0
|- + connect
|- + fetch
|- @ build
|- o deploy
`-> 3/4

-- line
,- * app_name  v0.1.0
|- ===============-----  3/4
`-> build

-- log
,- * app_name  v0.1.0
|- 1/4 connect ....... ok
|- 2/4 fetch ......... ok
`-> 3/4 build

-- end
`- + done in 4.2s
`- x failed at 3/4 build
````

### 7.3 Scripts

- `check.sh`: for each language in `bash rust python go`, run the demo command of 7.1 twice, with `NO_COLOR=1 TERM=xterm-256color` against `expected/demo.txt` and with `NO_COLOR=1 TERM=linux` against `expected/demo-ascii.txt`; diff the output; print one line per language (`bash ok`, `rust FAIL`, and so on); exit non-zero on any difference. Uses `bash`, `rustc`, `python3`, `go` from the Arch packages, and cleans up its temporary directories.
- `preview.sh [lang]`: run the demos of 7.1 with color on (output to the terminal, no `NO_COLOR`) so the look can be inspected, then run `steps.<ext>` for every style so the redraw can be seen. With no argument it runs all four languages, with a header per language.

`check.sh` is for building this repo and for sanity after edits. It is **not** run by anything automatically.

---

## 8. License

MIT, one `LICENSE` file at the repo root. Each template file's header contains `SPDX-License-Identifier: MIT` and the copyright line, so the notice travels with every copy. The copyright line is `Copyright (c) 2026 Kamilo Melo`; use exactly that text in `LICENSE` and in every template header.

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

---

## 9. Build phases and acceptance criteria

**Phase 1: design and expected output**
- [ ] `DESIGN.md` (copy of section 4), `LICENSE`, `README.md` skeleton (with the repo link and badge-free one-line description from 3.1).
- [ ] Git remote set to `git@github.com:kamilomelo/kartouche.git` and the phase commit made (no push).
- [ ] `expected/demo.txt` and `expected/demo-ascii.txt` exactly as in sections 7.2 and 7.2b.

**Phase 2: Bash and Python templates**
- [ ] `bash/kartouche.sh` and `python/kartouche.py` with the config/code markers (5.1), the render functions, `banner`, `card`, `Steps` and end states, following 5.2 and 5.5.
- [ ] `demo/demo.sh`, `demo/demo.py`, `demo/steps.sh` and `demo/steps.py`; the demos match both expected files (Unicode and `TERM=linux`).

**Phase 3: Rust and Go templates**
- [ ] `rust/kartouche.rs` and `go/kartouche.go` with the same structure and behavior.
- [ ] `demo/demo.rs`, `demo/demo.go`, `demo/steps.rs` and `demo/steps.go`; the demos match both expected files, run with the commands of 7.1.

**Phase 4: local acceptance checks (required; the agent can run all of these without a raw console)**

Use `script -qec '<command>' /dev/null` (util-linux, on Arch by default) or Python's `pty` module to run a command under a pseudo-terminal. `<steps>` means any language's `demo/steps.<ext>` command from 7.1, and every language must pass every check.
- [ ] **Piped output:** `<steps> rail 2>&1 | cat` contains no ESC byte (`\x1b`) and no `\r`; it is the append-only layout (title line, `├─ n/N label` per step, end line). Same with `NO_COLOR=1` on a pseudo-terminal: no color sequences.
- [ ] **Go terminal detection:** call a Go render function with an open `os.DevNull` file as its destination and `NO_COLOR` unset; the returned text has no escape sequences. A character device alone must not enable color.
- [ ] **Redraw:** under a pseudo-terminal with `TERM=xterm-256color`, `<steps> rail` output contains cursor-up (`\x1b[<k>A`, k ≥ 1) and line-clear (`\x1b[2K`) sequences and ends with the end line.
- [ ] **Linux console glyphs:** under a pseudo-terminal with `TERM=linux`, `<steps> dots` output has only 7-bit ASCII bytes once escape sequences are stripped (for example `LC_ALL=C grep -qP '[^\x00-\x7F]'` finds nothing).
- [ ] **Append-only config:** under a pseudo-terminal with `TERM=xterm-256color` and `APPEND_ONLY` set to true in a **temporary copy** of each template's config (do not change the repo templates), `<steps> rail` emits one step line per step and an end line, with no cursor-up (`\x1b[<k>A`), line-clear (`\x1b[2K`) or carriage-return (`\r`) sequences. Run the matching steps demo against that temporary template copy; this verifies the config switch independently of the non-terminal fallback.
- [ ] **Bash strict mode:** a script with `set -euo pipefail` that sources `bash/kartouche.sh` and calls every public function exits 0.
- [ ] **Single-file copy:** copying one template file alone into an empty directory and calling the banner works with no other files.

**Manual follow-up (for me, not required for the agent to finish; list these as "not performed" in the final report)**
- Look at `./preview.sh` in my own terminal.
- `Ctrl-C` in the middle of a steps run.
- The real Arch ISO console (`TERM=linux`), including inside `tmux` with `GLYPHS="ascii"`.
- An installer run with `pacstrap` output between steps and `APPEND_ONLY` set to true.

**Phase 5: docs and updater**
- [ ] `UPDATING.md` as in section 6.
- [ ] `README.md`: text-block screenshots of banner, card and all step styles, four quickstarts of under 10 lines each, the accent table, a one-line note that the card assumes single-width characters, and the "how to update a script" recipe.
- [ ] `check.sh` and `preview.sh` working (7.3).

**Done means:** `check.sh` prints four `ok` lines (each covering both the Unicode and the `TERM=linux` run), every checkbox of phases 1 to 5 is ticked (the manual follow-up in phase 4 is excluded), and the final report lists the manual follow-ups as not performed. The line-count targets of section 0 are not part of this definition.

---

## 10. Upgrade path

If copy-and-modify stops being enough (many scripts, frequent design changes, wanting one place to fix bugs), the full plan in `SPEC-library.md` turns these templates into a versioned, conformance-tested library. The look (section 4), the color roles and the demo output carry over unchanged, so nothing built here is thrown away.
