#!/bin/sh

# start the render server in a background process
cd /base/render-server && nohup npm run start > /dev/null 2>&1 &

# start the main go server
/base/bin/server
