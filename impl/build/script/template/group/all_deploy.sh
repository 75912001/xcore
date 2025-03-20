#!/bin/bash

scriptPath=$(dirname "$PWD")/script
source $scriptPath/group.base.sh

source ./base.sh

currentPath=$(cd $(dirname ${BASH_SOURCE[0]}); pwd)
groupID=$(basename "$PWD")

echo "deploy groupID:${groupID} service ..."

CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o ${exeName} ../../server/main/main.go

echo "build is complete"

ssh ${USER}@${IP} -p ${PORT} "
mkdir -p /data/game/group${groupID}
"

scp -P ${PORT} ${exeName} ${USER}@${IP}:/data/game/group${groupID}/${exeName}


# 循环遍历数组元素
for serviceName in "${serviceArr[@]}"; do
    echo "当前服务名: $serviceName"
done

scp -P ${PORT} $(dirname ${scriptPath})/bench/*.tpl ${USER}@${IP}:/data/game/group${groupID}/


for elementName in `ls $currentPath`
do
  dir_or_file=$currentPath"/"$elementName
  if [ -d $dir_or_file ]
  then
    name=`basename "$dir_or_file"`
    cd $dir_or_file
    currentPathService=$(cd $dir_or_file;pwd)
    cd - >> /dev/null 2>&1
    for elementID in `ls $currentPathService`
    do
      dir_or_fileID=$currentPathService"/"$elementID
      if [ -d $dir_or_fileID ]
      then
        id=`basename "$dir_or_fileID"`
        ssh ${USER}@${IP} -p ${PORT} "
        \cp -sf /data/game/group${groupID}/${exeName} /data/game/group${groupID}/${name}${id}/${groupID}.${name}.${id}.exe
        \cp -sf /data/game/group${groupID}/bench.json /data/game/group${groupID}/${name}${id}/bench.json
        "
        echo "deploy service done. ${groupID}.${name}.${id} "
      fi
    done
  fi
done