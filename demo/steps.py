import sys
import time

from kartouche import Steps

s = Steps(sys.argv[1], ["connect", "fetch", "build", "deploy"])
for _ in range(3):
    time.sleep(0.3)
    s.next()
time.sleep(0.3)
s.done()
