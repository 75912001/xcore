#!/bin/bash

serviceName="gateway"
groupId=1
serviceId=1

preName="${groupId}.${serviceName}.${serviceId}"

# 获取当前脚本文件所在路径的绝对路径
currentPath=$(realpath "$(dirname "$0")")
#项目绝对路径
projectPath=$(dirname "${currentPath}")
echo "projectPath:${projectPath}"

./${preName}.exe ${groupId} ${serviceName} ${serviceId}

exit 0

