#!/bin/bash

source ./base.sh

currentPath=$(cd `dirname $0`;pwd)

for elementName in `ls $currentPath`
do
  dir_or_file=$currentPath"/"$elementName
  if [ -d $dir_or_file ]
  then
    name=`basename "$dir_or_file"`
    zone=${serviceNameZoneIdMap[$name]}
    cd $dir_or_file
    currentPathService=$(cd $dir_or_file;pwd)
    for elementID in `ls $currentPathService`
    do
      dir_or_fileID=$currentPathService"/"$elementID
      if [ -d $dir_or_fileID ]
      then
        id=`basename "$dir_or_fileID"`
        cd ./${id}
        sh rm_log.sh
        cd - >> /dev/null 2>&1
      fi
    done
  fi
done