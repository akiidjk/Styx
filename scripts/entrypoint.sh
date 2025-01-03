#!/bin/sh
python3 -m http.server 8080 --bind 0.0.0.0 > /dev/null 2>&1 &
exec ./bin/styx -i eth0 -p 172.18.0.1
