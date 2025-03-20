#!/bin/bash

source ./base.sh

currentPath=$(cd `dirname $0`;pwd)

for elementName in "${startServiceArr[@]}"; do
  dir_or_file=$currentPath"/"$elementName
  if [ -d $dir_or_file ]
  then
    name=`basename "$dir_or_file"`
    cd $dir_or_file
    currentPathService=$(cd $dir_or_file;pwd)

    for elementID in `ls $currentPathService`
    do
      dir_or_fileID=$currentPathService"/"$elementID
      if [ -d $dir_or_fileID ]
      then
        id=`basename "$dir_or_fileID"`
        cd ./${id}
        sh start.sh
        # 某个服务启动失败  终止
        exit_status=$?
        if [ $exit_status -ne 0 ]; then
          exit 1
        fi
        cd - >> /dev/null 2>&1
      fi
    done
  fi
done