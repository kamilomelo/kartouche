import sys

from kartouche import banner, card, summary

print("Loading configuration...", file=sys.stderr)
banner()
print("Checking inputs...", file=sys.stderr)
card()
print("Reviewing results...", file=sys.stderr)
summary("rail", [("connect", "ok"), ("fetch", "failed"), ("deploy", "skipped")])
print("Continuing...", file=sys.stderr)
