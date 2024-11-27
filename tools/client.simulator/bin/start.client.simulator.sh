#!/bin/bash

# 获取当前脚本文件所在路径的绝对路径
currentPath=$(realpath "$(dirname "$0")")
#项目绝对路径
projectPath=$(dirname "${currentPath}")
echo "projectPath:${projectPath}"

# 启动 client.simulator.exe
cd ${projectPath}/bin
./client.simulator.exe
cd -

exit 0

