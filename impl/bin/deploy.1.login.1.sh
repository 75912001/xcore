#!/bin/bash

serviceName="login"
groupId=1
serviceId=1

preName="${groupId}.${serviceName}.${serviceId}"
# 判断 ${preName}.bench.json 是否存在
# 存在-无动作
# 不存在-复制 /build/bench/${serviceName}.bench.json -> /bin/${preName}.bench.json

# 获取当前脚本文件所在路径的绝对路径
currentPath=$(realpath "$(dirname "$0")")
#项目绝对路径
projectPath=$(dirname "${currentPath}")
echo "projectPath:${projectPath}"

cp "${projectPath}"/bin/service.exe "${projectPath}"/bin/${preName}.exe

# 检查文件是否存在 0: 存在 1: 不存在
check_file_exists() {
  if [ -f "$1" ]; then
    return 0
  else
    echo -e "\e[91m$1 does not exist.\e[0m"
    return 1
  fi
}

check_file_exists "${projectPath}/bin/${preName}.bench.json"
if [ $? -ne 0 ]; then
  echo -e "\e[91mcp "${projectPath}"/build/bench/${serviceName}.bench.json  -> "${projectPath}"/bin/${preName}.bench.json\e[0m"
  cp "${projectPath}"/build/bench/${serviceName}.bench.json "${projectPath}"/bin/${preName}.bench.json
fi

exit 0

