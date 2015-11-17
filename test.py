#!/usr/bin/env python
# encoding: utf-8

import random as _random
import sys
AEIOU = "AEIOU"
BCDFG = "BCDFGHJKLMNPQRSTVWXYZ"


def make_random_nickname(l):
    rr = ["" for i in range(_random.randint(l/2,l))]

    for i,r in enumerate(rr):
        for x in range(_random.randint(3,8)):
            rr[i] += _random.choice(AEIOU if x % 2 == 1 else BCDFG + _random.choice(["", "h","r"]))
        rr[i] = rr[i].title()

    return " ".join(rr)


for i in xrange(int(sys.argv[1])):
    print i+1
    title = make_random_nickname(8)
    with open("./content/%s.html" % title, "wb+") as f:
        data = ("<h1>%s</h1>" % title)

        data += "<p>%s</p>"% make_random_nickname(1000)
        data += "<p>%s</p>"% make_random_nickname(1000)
        data += "<p>%s</p>"% make_random_nickname(1000)
        data += "Tags:"
        for t in make_random_nickname(8).split(" "):
            data += "<span class='tag'>%s</span>&nbsp;" % t
        data += "<span class='tag'>Tag</span>&nbsp;" 

        f.write(data)
        import time
        time.sleep(0.05)
