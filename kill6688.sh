#!/bin/sh
## 执行netstat命令，查找监听在6688端口的进程
pid=$(netstat -lnp | grep ':6688 ' | awk '{print $7}' | cut -d'/' -f1)

if [ -n "$pid" ]; then
# 终止对应进程
sudo kill -9 "$pid"
echo "Process with PID $pid has been killed"
else
echo "No process listening on port 6688"
fi
