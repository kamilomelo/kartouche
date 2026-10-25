#!/usr/bin/env bash
source "$(dirname "${BASH_SOURCE[0]}")/../bash/kartouche.sh"
kt_render_banner 1 app_name 0.1.0 app_description
printf '\n'
kt_render_card 1 lumen 0.3.1 'tiny notes, zero friction'
for style in inline dots crumbs rail line log; do
    printf '\n-- %s\n' "$style"
    kt_render_title 1 app_name 0.1.0
    kt_render_frame 1 "$style" 3 4 build connect fetch build deploy
done
printf '\n-- end\n'
kt_render_end 1 1 'done in 4.2s'
kt_render_end 1 0 'failed at 3/4 build'
