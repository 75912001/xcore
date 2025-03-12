#!/bin/bash

source ./base.sh

currentPath=$(cd `dirname $0`;pwd)

# 按倒序
for ((i=${#startServiceArr[@]}-1; i>=0; i--)); do
  elementName=${startServiceArr[i]}
  #echo $elementName
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
        sh stop.sh
        cd - >> /dev/null 2>&1
      fi
    done
  fi
done