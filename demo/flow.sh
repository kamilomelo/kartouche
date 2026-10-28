#!/usr/bin/env bash
source "$(dirname "${BASH_SOURCE[0]}")/../bash/kartouche.sh"
printf 'Loading configuration...\n' >&2
kt_banner
printf 'Checking inputs...\n' >&2
kt_card
printf 'Reviewing results...\n' >&2
kt_summary rail connect ok fetch failed deploy skipped
printf 'Continuing...\n' >&2
