#!/bin/bash
# 测试代码

#当前脚本文件所在路径的绝对路径
currentPath=$(realpath "$(dirname "$0")")
#项目绝对路径
projectPath="$(dirname "${currentPath}")"
echo -e "\e[92m项目路径 ${projectPath}\e[0m"

#测试代码
echo -e "\e[93m======测试代码... \e[0m"

# todo menglc [调用所有测试用例]

echo -e "\e[92m======测试代码完成\e[0m"