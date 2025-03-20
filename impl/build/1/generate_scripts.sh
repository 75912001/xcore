#!/bin/bash

currentPath=$(cd $(dirname ${BASH_SOURCE[0]}); pwd)
source $(dirname $currentPath)/script/group.base.sh

mkdir -p $currentPath/${serviceNameLogic}/1
mkdir -p $currentPath/${serviceNameLogic}/2
mkdir -p $currentPath/${serviceNameGateway}/1
mkdir -p $currentPath/${serviceNameGateway}/2
mkdir -p $currentPath/${serviceNameLogin}/1
mkdir -p $currentPath/${serviceNameLogin}/2

groupID=$(basename "$PWD")

serviceTemplateSrcPathFiles=$(dirname "$PWD")/script/template/service/*
serviceTemplateDstPathPrefix=$(dirname "$PWD")

groupTemplateSrcPathFiles=$(dirname "$PWD")/script/template/group/*
groupTemplateDstPathPrefix="$PWD"

# 生成每个服务的脚本

for elementName in `ls $currentPath`
do
  dir_or_file=$currentPath"/"$elementName
  if [ -d $dir_or_file ]
  then
    name=`basename "$dir_or_file"`
    cd $dir_or_file
    currPathService=$(cd $dir_or_file;pwd)
    cd - >> /dev/null 2>&1
    for elementID in `ls $currPathService`
    do
      dir_or_fileID=$currPathService"/"$elementID
      if [ -d $dir_or_fileID ]
      then
        id=`basename "$dir_or_fileID"`
        \cp ${serviceTemplateSrcPathFiles} ${serviceTemplateDstPathPrefix}/${groupID}/${name}/${id}
      fi
    done
  fi
done

# 生成 all*.sh 脚本

\cp ${groupTemplateSrcPathFiles} ${groupTemplateDstPathPrefix}
