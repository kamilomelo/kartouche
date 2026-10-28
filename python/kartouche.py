# kartouche — copy this file into your project; https://github.com/kamilomelo/kartouche
# SPDX-License-Identifier: MIT
# Copyright (c) 2026 Kamilo Melo
# KARTOUCHE_TEMPLATE_VERSION: 4

# kartouche:config begin
NAME = "app_name"
VERSION = "0.1.0"
DESCRIPTION = "app_description"
ACCENT = "33"
APPEND_ONLY = False
GLYPHS = "auto"
PAUSE_SECONDS = 1
# kartouche:config end

# kartouche:code begin v4
import os
import sys
import time

_UNICODE = "╭╮╰╯─│├❯●○◉›━✓✗╸·"
_ASCII = ",.`'-||>*o@>=+x=."


def _glyphs() -> dict[str, str]:
    ascii_mode = GLYPHS == "ascii" or (GLYPHS == "auto" and os.environ.get("TERM") in ("linux", "dumb"))
    return dict(zip(_UNICODE, _ASCII if ascii_mode else _UNICODE))


def _text(value: str) -> str:
    return value.replace("\n", " ")


def _paint(value: str, role: str, stream) -> str:
    if not value or not stream.isatty() or os.environ.get("NO_COLOR"):
        return value
    code = {"dim": "2", "bold": "1", "accent": ACCENT, "boldaccent": "1;" + ACCENT, "green": "32", "red": "31"}[role]
    return f"\x1b[{code}m{value}\x1b[0m"


def _version(value: str) -> str:
    value = _text(value)
    return value if value.startswith("v") else "v" + value


def render_title(name, version, stream) -> str:
    g = _glyphs()
    suffix = "  " + _paint(_version(version), "dim", stream) if version else ""
    return (_paint(g["╭"] + g["─"], "dim", stream) + " " +
            _paint(g["●"], "accent", stream) + " " + _paint(_text(name), "bold", stream) + suffix)


def render_banner(name, version, description, status, stream) -> str:
    g = _glyphs()
    lines = [render_title(name, version, stream)]
    if status:
        prefix = "├─" if description else "╰─❯"
        lines.append(_paint("".join(g[c] for c in prefix), "dim", stream) + " " + _text(status))
    if description:
        lines.append(_paint("".join(g[c] for c in "╰─❯"), "dim", stream) + " " + _text(description))
    return "\n".join(lines)


def render_card(name, version, description, stream) -> str:
    g = _glyphs()
    name, description = _text(name), _text(description)
    v = _version(version)
    left, right = len("─ " + name + " "), len("  " + v + " ─")
    w = max(len(description) + 4, left + right + 2)
    dim = lambda s: _paint("".join(g.get(c, c) for c in s), "dim", stream)
    top = (dim("╭─ ") + _paint(name, "bold", stream) + dim(" " + "─" * (w-left-right) + "  ") +
           _paint(v, "dim", stream) + dim(" ─╮"))
    body = dim("│") + "  " + description.ljust(w-4) + "  " + dim("│")
    return "\n".join((top, body, dim("╰" + "─"*w + "╯")))


def render_end(ok, text, stream) -> str:
    g = _glyphs()
    return (_paint(g["╰"] + g["─"], "dim", stream) + " " +
            _paint(g["✓" if ok else "✗"], "green" if ok else "red", stream) + " " + _text(text))


def render_frame(style, n, total, label, labels, stream) -> str:
    g = _glyphs()
    labels = [_text(x) for x in labels]
    label = _text(label)
    if style == "dots" and total > 12:
        style = "line"
    prefix = _paint(g["├"] + g["─"], "dim", stream) + " "
    tail = _paint(g["╰"] + g["─"] + g["❯"], "dim", stream) + " "
    count = f"{min(n, total)}/{total}"
    done = n > total
    lines = []
    if style == "dots":
        lines.append(prefix + _paint(g["●"] * min(n, total), "accent", stream) +
                     _paint(g["○"] * max(total-n, 0), "dim", stream) + "  " + _paint(count, "dim", stream))
    elif style == "crumbs":
        parts = [_paint(x, "boldaccent" if i == n-1 else "dim", stream) for i, x in enumerate(labels)]
        lines.append(prefix + (" " + _paint(g["›"], "dim", stream) + " ").join(parts))
    elif style == "rail":
        for i, x in enumerate(labels):
            glyph, role = ("✓", "green") if i < n-1 or done else (("◉", "accent") if i == n-1 else ("○", "dim"))
            text = _paint(x, "bold", stream) if i == n-1 and not done else x
            lines.append(prefix + _paint(g[glyph], role, stream) + " " + text)
    elif style == "line":
        filled = (20 * min(n, total) * 2 + total) // (2 * total)
        bar = g["━"] * 20 if filled == 20 else (g["━"] * (filled-1) + g["╸"] if filled else "")
        lines.append(prefix + _paint(bar, "accent", stream) +
                     _paint(g["─"] * (20-filled), "dim", stream) + "  " + _paint(count, "dim", stream))
    elif style == "log":
        longest = max(map(len, labels), default=0)
        for i in range(min(n-1, total)):
            dots = g["·"] * (longest-len(labels[i])+7)
            lines.append(prefix + f"{i+1}/{total} {labels[i]} {dots} " + _paint("ok", "accent", stream))
    if style not in ("inline", "dots", "crumbs", "rail", "line", "log"):
        style = "inline"
    if done:
        lines.append(tail + count)
    elif style == "rail":
        lines.append(tail + count)
    else:
        lines.append(tail + (label if style == "dots" or style == "line" else count + " " + label))
    return "\n".join(lines)


def render_summary(style: str, name: str, version: str, items: list[tuple[str, str]], stream) -> str:
    if style not in ("rail", "log") or any(state not in ("ok", "failed", "skipped") for _, state in items):
        return ""
    g = _glyphs()
    counts = {state: sum(result == state for _, result in items) for state in ("ok", "failed", "skipped")}
    longest = max((len(_text(label)) for label, _ in items), default=0)
    lines = [render_title(name, version, stream)]
    for i, (label, state) in enumerate(items, 1):
        label = _text(label)
        prefix = _paint(g["├"] + g["─"], "dim", stream) + " "
        if style == "rail":
            glyph, role = {"ok": ("✓", "green"), "failed": ("✗", "red"), "skipped": ("○", "dim")}[state]
            lines.append(prefix + _paint(g[glyph], role, stream) + " " + label)
        else:
            role = {"ok": "accent", "failed": "red", "skipped": "dim"}[state]
            dots = g["·"] * (longest - len(label) + 7)
            lines.append(prefix + f"{i}/{len(items)} {label} {dots} " + _paint(state, role, stream))
    glyph, role = (("✗", "red") if counts["failed"] else
                   ("○", "dim") if counts["skipped"] or not items else ("✓", "green"))
    end = f'{counts["ok"]} ok {g["·"]} {counts["failed"]} failed {g["·"]} {counts["skipped"]} skipped'
    lines.append(_paint(g["╰"] + g["─"], "dim", stream) + " " + _paint(g[glyph], role, stream) + " " + end)
    return "\n".join(lines)


def _after_block() -> None:
    if PAUSE_SECONDS > 0 and sys.stderr.isatty():
        time.sleep(PAUSE_SECONDS)
    try: print(file=sys.stderr)
    except OSError: pass


def _print_block(text: str) -> None:
    if not text:
        return
    try:
        print("\n" + text, file=sys.stderr)
        sys.stderr.flush()
    except OSError: pass
    _after_block()


def summary(style: str, items: list[tuple[str, str]]) -> None:
    _print_block(render_summary(style, NAME, VERSION, items, sys.stderr))


def banner(status: str | None = None) -> None:
    _print_block(render_banner(NAME, VERSION, DESCRIPTION, status, sys.stderr))


def card() -> None:
    _print_block(render_card(NAME, VERSION, DESCRIPTION, sys.stderr))


class Steps:
    def __init__(self, style: str, labels: list[str]) -> None:
        self.style, self.labels, self.n, self.finished = style, labels, int(bool(labels)), False
        self.started, self.previous = time.monotonic(), 0
        self.redraw = sys.stderr.isatty() and not APPEND_ONLY
        self._write("\n" + render_title(NAME, VERSION, sys.stderr) + "\n")
        if labels: self._show()

    def _write(self, text: str) -> None:
        try: sys.stderr.write(text); sys.stderr.flush()
        except OSError: pass

    def _show(self, end: str | None = None) -> None:
        if not self.redraw:
            self._write(end + "\n" if end else f"{_paint(_glyphs()['├'] + _glyphs()['─'], 'dim', sys.stderr)} {self.n}/{len(self.labels)} {_text(self.labels[self.n-1])}\n")
            return
        if not self.labels:
            self._write(end + "\n" if end else "")
            return
        n = len(self.labels)+1 if end and self.finished_ok else self.n
        region = render_frame(self.style, n, len(self.labels), self.labels[self.n-1], self.labels, sys.stderr).split("\n")
        if end: region[-1] = end
        if self.previous: self._write(f"\x1b[{self.previous}A")
        for line in region: self._write(("\r\x1b[2K" if self.previous else "") + line + "\n")
        self.previous = len(region)

    def next(self) -> None:
        if not self.finished and self.n < len(self.labels):
            self.n += 1
            self._show()

    def done(self, elapsed: str | None = None) -> None:
        if self.finished: return
        self.finished, self.finished_ok = True, True
        text = "done in " + (_text(elapsed) if elapsed is not None else f"{time.monotonic()-self.started:.1f}s")
        self._show(render_end(True, text, sys.stderr))
        _after_block()

    def fail(self, reason: str | None = None) -> None:
        if self.finished: return
        self.finished, self.finished_ok = True, False
        text = f"failed at {self.n}/{len(self.labels)} {_text(self.labels[self.n-1])}" if self.labels else "failed"
        if reason: text += " · " + _text(reason)
        self._show(render_end(False, text, sys.stderr))
        _after_block()

    def __enter__(self): return self

    def __exit__(self, kind, value, traceback):
        if kind: self.fail()
        else: self.done()
        return False
# kartouche:code end
