#!/bin/bash
# 检查代码风格

#当前脚本文件所在路径的绝对路径
currentPath=$(realpath "$(dirname "$0")")
#项目绝对路径
projectPath="$(dirname "${currentPath}")"
echo -e "\e[92m项目路径 ${projectPath}\e[0m"

#检查代码
echo -e "\e[93m======检查代码... \e[0m"
cd  "${projectPath}" || exit
echo "gofmt..."
gofmt -w .
echo "golangci-lint..."
golangci-lint run
cd - || exit
echo -e "\e[92m======检查代码完成\e[0m"