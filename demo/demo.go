package main

import (
	"fmt"
	"os"
)

func main() {
	out := os.Stdout
	fmt.Println(RenderBanner("app_name", "0.1.0", "app_description", "", out))
	fmt.Println()
	fmt.Println(RenderCard("lumen", "0.3.1", "tiny notes, zero friction", out))
	for _, pair := range []struct {
		name  string
		style Style
	}{
		{"inline", Inline}, {"dots", Dots}, {"crumbs", Crumbs}, {"rail", Rail}, {"line", Line}, {"log", Log},
	} {
		fmt.Printf("\n-- %s\n", pair.name)
		fmt.Println(RenderTitle("app_name", "0.1.0", out))
		fmt.Println(RenderFrame(pair.style, 3, 4, "build", []string{"connect", "fetch", "build", "deploy"}, out))
	}
	fmt.Println("\n-- end")
	fmt.Println(RenderEnd(true, "done in 4.2s", out))
	fmt.Println(RenderEnd(false, "failed at 3/4 build", out))
	items := []SummaryItem{{"connect", OkOutcome}, {"fetch", FailedOutcome}, {"deploy", SkippedOutcome}}
	for _, pair := range []struct {
		name  string
		style Style
	}{{"rail", Rail}, {"log", Log}} {
		fmt.Printf("\n-- summary %s\n", pair.name)
		fmt.Println(RenderSummary(pair.style, "app_name", "0.1.0", items, out))
	}
}
