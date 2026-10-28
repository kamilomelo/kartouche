package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "Loading configuration...")
	Banner("")
	fmt.Fprintln(os.Stderr, "Checking inputs...")
	Card()
	fmt.Fprintln(os.Stderr, "Reviewing results...")
	Summary(Rail, []SummaryItem{{"connect", OkOutcome}, {"fetch", FailedOutcome}, {"deploy", SkippedOutcome}})
	fmt.Fprintln(os.Stderr, "Continuing...")
}
