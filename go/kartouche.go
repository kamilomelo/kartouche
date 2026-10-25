// kartouche — copy this file into your project; https://github.com/kamilomelo/kartouche
// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Kamilo Melo
// KARTOUCHE_TEMPLATE_VERSION: 3
// Change package main to match your package when copying this file elsewhere.
package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"
	"unicode/utf8"
	"unsafe"
)

// kartouche:config begin
const NAME = "app_name"
const VERSION = "0.1.0"
const DESCRIPTION = "app_description"
const ACCENT = "38;5;244"
const APPEND_ONLY = false
const GLYPHS = "auto"

// kartouche:config end

// kartouche:code begin v3
type Style int

const (
	Inline Style = iota
	Dots
	Crumbs
	Rail
	Line
	Log
)

func terminal(f *os.File) bool {
	var term syscall.Termios
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), uintptr(syscall.TCGETS), uintptr(unsafe.Pointer(&term)))
	return err == 0
}
func glyphs(s string) string {
	if GLYPHS != "ascii" && !(GLYPHS == "auto" && (os.Getenv("TERM") == "linux" || os.Getenv("TERM") == "dumb")) {
		return s
	}
	unicode := []rune("╭╮╰╯─│├❯●○◉›━✓✗╸·")
	ascii := []rune(",.`'-||>*o@>=+x=.")
	return strings.Map(func(c rune) rune {
		for i, r := range unicode {
			if r == c {
				return ascii[i]
			}
		}
		return c
	}, s)
}
func clean(s string) string { return strings.ReplaceAll(s, "\n", " ") }
func paint(s, role string, f *os.File) string {
	if s == "" || !terminal(f) || os.Getenv("NO_COLOR") != "" {
		return s
	}
	code := "31"
	switch role {
	case "dim":
		code = "2"
	case "bold":
		code = "1"
	case "accent":
		code = ACCENT
	case "both":
		code = "1;" + ACCENT
	case "green":
		code = "32"
	}
	return "\x1b[" + code + "m" + s + "\x1b[0m"
}
func ver(s string) string {
	s = clean(s)
	if strings.HasPrefix(s, "v") {
		return s
	}
	return "v" + s
}
func prefix(f *os.File) string { return paint(glyphs("├─"), "dim", f) + " " }
func tail(f *os.File) string   { return paint(glyphs("╰─❯"), "dim", f) + " " }

func RenderTitle(name, version string, f *os.File) string {
	suffix := ""
	if version != "" {
		suffix = "  " + paint(ver(version), "dim", f)
	}
	return paint(glyphs("╭─"), "dim", f) + " " + paint(glyphs("●"), "accent", f) + " " + paint(clean(name), "bold", f) + suffix
}
func RenderBanner(name, version, description, status string, f *os.File) string {
	lines := []string{RenderTitle(name, version, f)}
	if status != "" {
		head := "├─"
		if description == "" {
			head = "╰─❯"
		}
		lines = append(lines, paint(glyphs(head), "dim", f)+" "+clean(status))
	}
	if description != "" {
		lines = append(lines, tail(f)+clean(description))
	}
	return strings.Join(lines, "\n")
}
func RenderCard(name, version, description string, f *os.File) string {
	name, version, description = clean(name), ver(version), clean(description)
	left, right := utf8.RuneCountInString(name)+3, utf8.RuneCountInString(version)+4
	w := max(utf8.RuneCountInString(description)+4, left+right+2)
	top := paint(glyphs("╭─ "), "dim", f) + paint(name, "bold", f) + paint(" "+glyphs(strings.Repeat("─", w-left-right))+"  ", "dim", f) + paint(version, "dim", f) + paint(glyphs(" ─╮"), "dim", f)
	body := paint(glyphs("│"), "dim", f) + "  " + description + strings.Repeat(" ", w-4-utf8.RuneCountInString(description)) + "  " + paint(glyphs("│"), "dim", f)
	return top + "\n" + body + "\n" + paint(glyphs("╰"+strings.Repeat("─", w)+"╯"), "dim", f)
}
func RenderEnd(ok bool, text string, f *os.File) string {
	symbol, role := "✗", "red"
	if ok {
		symbol, role = "✓", "green"
	}
	return paint(glyphs("╰─"), "dim", f) + " " + paint(glyphs(symbol), role, f) + " " + clean(text)
}
func RenderFrame(style Style, n, total int, label string, labels []string, f *os.File) string {
	if style == Dots && total > 12 {
		style = Line
	}
	done := n > total
	count := fmt.Sprintf("%d/%d", min(n, total), total)
	lines := []string{}
	switch style {
	case Dots:
		lines = append(lines, prefix(f)+paint(glyphs(strings.Repeat("●", min(n, total))), "accent", f)+paint(glyphs(strings.Repeat("○", max(total-n, 0))), "dim", f)+"  "+paint(count, "dim", f))
	case Crumbs:
		parts := make([]string, len(labels))
		for i, x := range labels {
			role := "dim"
			if i == n-1 && !done {
				role = "both"
			}
			parts[i] = paint(clean(x), role, f)
		}
		lines = append(lines, prefix(f)+strings.Join(parts, " "+paint(glyphs("›"), "dim", f)+" "))
	case Rail:
		for i, x := range labels {
			symbol, role := "○", "dim"
			if i < n-1 || done {
				symbol, role = "✓", "green"
			} else if i == n-1 {
				symbol, role = "◉", "accent"
			}
			text := clean(x)
			if i == n-1 && !done {
				text = paint(text, "bold", f)
			}
			lines = append(lines, prefix(f)+paint(glyphs(symbol), role, f)+" "+text)
		}
	case Line:
		filled := (40*min(n, total) + total) / (2 * total)
		bar := ""
		if filled == 20 {
			bar = strings.Repeat("━", 20)
		} else if filled > 0 {
			bar = strings.Repeat("━", filled-1) + "╸"
		}
		lines = append(lines, prefix(f)+paint(glyphs(bar), "accent", f)+paint(glyphs(strings.Repeat("─", 20-filled)), "dim", f)+"  "+paint(count, "dim", f))
	case Log:
		longest := 0
		for _, x := range labels {
			longest = max(longest, utf8.RuneCountInString(clean(x)))
		}
		for i := 0; i < min(n-1, total); i++ {
			text := clean(labels[i])
			dots := glyphs(strings.Repeat("·", longest-utf8.RuneCountInString(text)+7))
			lines = append(lines, fmt.Sprintf("%s%d/%d %s %s %s", prefix(f), i+1, total, text, dots, paint("ok", "accent", f)))
		}
	}
	last := count
	if !done && style != Rail {
		if style == Dots || style == Line {
			last = clean(label)
		} else {
			last = count + " " + clean(label)
		}
	}
	return strings.Join(append(lines, tail(f)+last), "\n")
}
func Banner(status string) {
	_, _ = fmt.Fprintln(os.Stderr, RenderBanner(NAME, VERSION, DESCRIPTION, status, os.Stderr))
}
func Card() { _, _ = fmt.Fprintln(os.Stderr, RenderCard(NAME, VERSION, DESCRIPTION, os.Stderr)) }

// Steps is not concurrency-safe.
type Steps struct {
	style            Style
	labels           []string
	n, previous      int
	redraw, finished bool
	started          time.Time
}

func NewSteps(style Style, labels ...string) *Steps {
	s := &Steps{style: style, labels: append([]string(nil), labels...), started: time.Now(), redraw: terminal(os.Stderr) && !APPEND_ONLY}
	_, _ = fmt.Fprintln(os.Stderr, RenderTitle(NAME, VERSION, os.Stderr))
	if len(labels) > 0 {
		s.n = 1
		s.show("", false)
	}
	return s
}
func (s *Steps) show(end string, done bool) {
	if !s.redraw || len(s.labels) == 0 {
		if end != "" {
			_, _ = fmt.Fprintln(os.Stderr, end)
		} else if s.n > 0 {
			_, _ = fmt.Fprintf(os.Stderr, "%s%d/%d %s\n", prefix(os.Stderr), s.n, len(s.labels), clean(s.labels[s.n-1]))
		}
		return
	}
	n := s.n
	if done {
		n = len(s.labels) + 1
	}
	lines := strings.Split(RenderFrame(s.style, n, len(s.labels), s.labels[s.n-1], s.labels, os.Stderr), "\n")
	if end != "" {
		lines[len(lines)-1] = end
	}
	if s.previous > 0 {
		_, _ = fmt.Fprintf(os.Stderr, "\x1b[%dA", s.previous)
	}
	for _, line := range lines {
		if s.previous > 0 {
			_, _ = fmt.Fprint(os.Stderr, "\r\x1b[2K")
		}
		_, _ = fmt.Fprintln(os.Stderr, line)
	}
	s.previous = len(lines)
}
func (s *Steps) Next() {
	if !s.finished && s.n < len(s.labels) {
		s.n++
		s.show("", false)
	}
}
func (s *Steps) Done(elapsed string) {
	if s.finished {
		return
	}
	s.finished = true
	if elapsed == "" {
		elapsed = fmt.Sprintf("%.1fs", time.Since(s.started).Seconds())
	}
	s.show(RenderEnd(true, "done in "+clean(elapsed), os.Stderr), true)
}
func (s *Steps) Fail(reason string) {
	if s.finished {
		return
	}
	s.finished = true
	text := "failed"
	if len(s.labels) > 0 {
		text = fmt.Sprintf("failed at %d/%d %s", s.n, len(s.labels), clean(s.labels[s.n-1]))
	}
	if reason != "" {
		text += " · " + clean(reason)
	}
	s.show(RenderEnd(false, text, os.Stderr), false)
}
func (s *Steps) Close() {
	if s.finished {
		return
	}
	s.finished = true
	text := "interrupted"
	if len(s.labels) > 0 {
		text = fmt.Sprintf("interrupted at %d/%d %s", s.n, len(s.labels), clean(s.labels[s.n-1]))
	}
	s.show(RenderEnd(false, text, os.Stderr), false)
}

// kartouche:code end
