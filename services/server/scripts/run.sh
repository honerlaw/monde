#!/bin/sh

# start the render server in a background process
cd /base/server/renderer && nohup npm run start > /dev/null 2>&1 &

# start the main go server
/base/server/bin/server

