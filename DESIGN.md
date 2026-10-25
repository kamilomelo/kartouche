# Design (the look)

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

**rail**: one line per step. `✓` (green) done, `◉` (accent, bold label) current, `○` (dim) pending.
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

`✓` is green. `✗` is red. Elapsed time is optional text supplied by the caller or computed by the template as seconds with one decimal (`4.2s`).

### 4.5 Spacing and text rules

- Always two spaces between name and version, two spaces before a counter.
- Newlines in any text field are replaced by a single space.
- Text is printed literally. Never interpret `%`, `\`, `$` or backslash sequences in user text (Bash: use `printf '%s'`).

### 4.6 Color

| Role | Style |
|---|---|
| frame glyphs, version, counter, pending items | dim `\e[2m` |
| name | bold `\e[1m` |
| accent: `●`, `◉`, filled bar, `ok`, current crumb | the config constant `ACCENT` |
| success: every `✓` (end line and completed rail steps) | green `\e[32m` |
| failure: `✗` | red `\e[31m` |
| every styled run ends with | reset `\e[0m` |

- `ACCENT` is a raw SGR parameter string chosen by hand when the template is copied: `33` (yellow, default), `38;5;244` (gray), `36` (cyan), and so on. There is no name parser, no capability detection and no fallback logic. Success green is fixed, not configurable.
- Approved accent suggestions to list in `README.md` and `DESIGN.md`:

| Look | `ACCENT` |
|---|---|
| yellow / amber (default) | `33` |
| gray | `38;5;244` |
| cyan | `36` |
| magenta | `35` |
| green | `32` |
| blue | `34` |
| soft purple | `38;5;141` |
| sky blue | `38;5;81` |
| warm orange | `38;5;214` |
| mint | `38;5;78` |
| coral | `38;5;203` |

- Readability note for the docs: `30` (black) and `90` (gray) vanish on some themes, so avoid them as accents.
- One configurable accent per app is the intended use; success stays green and failure stays red. The frame, glyphs and layout stay identical across apps.

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
