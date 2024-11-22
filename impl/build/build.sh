#!/bin/bash

# 获取当前脚本文件所在路径的绝对路径
currentPath=$(realpath "$(dirname "$0")")
#项目绝对路径
projectPath=$(dirname "${currentPath}")
echo "projectPath:${projectPath}"

# 编译 main.go
go build -o service.exe "${projectPath}"/service/main/main.go
# 将 service.exe 移动到 bin 目录
mv service.exe "${projectPath}"/bin/
echo -e "\e[92mmv service.exe to "${projectPath}"/bin/ successfully.\e[0m"

echo -e "\e[92mbuild service successfully.\e[0m"