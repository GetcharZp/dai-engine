#!/bin/sh

curl -i -X POST \
http://127.0.0.1:13100/my/My/SayHello \
-H 'Content-Type: application/json' \
-d '{"name":"get"}'