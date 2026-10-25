import sys

import kartouche as k

out = sys.stdout
print(k.render_banner("app_name", "0.1.0", "app_description", None, out))
print()
print(k.render_card("lumen", "0.3.1", "tiny notes, zero friction", out))
for style in ("inline", "dots", "crumbs", "rail", "line", "log"):
    print("\n--", style)
    print(k.render_title("app_name", "0.1.0", out))
    print(k.render_frame(style, 3, 4, "build", ["connect", "fetch", "build", "deploy"], out))
print("\n-- end")
print(k.render_end(True, "done in 4.2s", out))
print(k.render_end(False, "failed at 3/4 build", out))
