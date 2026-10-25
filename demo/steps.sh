#!/usr/bin/env bash
source "$(dirname "${BASH_SOURCE[0]}")/../bash/kartouche.sh"
kt_steps_start "${1:-inline}" connect fetch build deploy
for _ in 1 2 3; do
    sleep 0.3
    kt_steps_next
done
sleep 0.3
kt_steps_done
