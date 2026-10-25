package main

import (
    "os"
    "time"
)

func main() {
    style := Inline
    if len(os.Args) > 1 { switch os.Args[1] {
        case "dots": style = Dots; case "crumbs": style = Crumbs
        case "rail": style = Rail; case "line": style = Line; case "log": style = Log
    } }
    s := NewSteps(style, "connect", "fetch", "build", "deploy")
    defer s.Close()
    for i := 0; i < 3; i++ { time.Sleep(300*time.Millisecond); s.Next() }
    time.Sleep(300*time.Millisecond)
    s.Done("")
}
