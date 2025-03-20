#!/bin/bash

# 获取当前脚本文件所在路径的绝对路径
currentPath=$(realpath "$(dirname "$0")")
echo "currentPath:${currentPath}"
#项目绝对路径
projectPath=$(dirname "${currentPath}")
echo "projectPath:${projectPath}"

cd "${projectPath}"/client.simulator/main || exit
# 编译 main.go
go build

# 获取 返回值
if [ $? -ne 0 ]; then
    echo -e "\e[91mbuild client.simulator failed.\e[0m"
    exit 1
fi

cd - || exit

# 将 client.simulator.exe 移动到 bin 目录
mv ./main/main.exe "${projectPath}"/client.simulator/bin/client.simulator.exe
echo -e "\e[92mmv main.exe to "${projectPath}"/client.simulator/bin/ successfully.\e[0m"

echo -e "\e[92mbuild client.simulator successfully.\e[0m"