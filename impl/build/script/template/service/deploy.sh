#!/bin/bash

source group.base.sh
source ./base.sh

mainPathFile=$(dirname $(dirname $(dirname $(dirname "$PWD"))))/impl/process/main/main.go
benchJsonPathFile=$(dirname $(dirname "$PWD"))/bench.json

source ../../../script/deploy.sh ${zoneID} ${serviceName} ${serviceID} ${IP} ${PORT} ${USER} ${mainPathFile} ${benchJsonPathFile}

if [ "${serviceName}" == "${serviceNameGM}" ] || [ "${serviceName}" == "${serviceNameServerList}" ];then
  echo "deploy ${serviceName} script input_etcd.sh"
  scp -P ${PORT} input_etcd.sh ${USER}@${IP}:/data/yoozoogame/code/zone${zoneID}/${serviceName}${serviceID}/
fi

if [ "${serviceName}" == "${serviceNameServerList}" ];then
  echo "deploy ${serviceName} /data/yoozoogame/config${zoneID}/Server.xml"
  serverXmlPathFile=$(dirname $(dirname "$PWD"))/config/Server.xml
  scp -P ${PORT} ${serverXmlPathFile} ${USER}@${IP}:/data/yoozoogame/config${zoneID}/
fi