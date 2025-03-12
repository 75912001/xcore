#!/bin/bash

source group.base.sh
source ./base.sh

if [ "${serviceName}" == "${serviceNameGM}" ] || [ "${serviceName}" == "${serviceNameServerList}" ];then
  ssh ${USER}@${IP} -p ${PORT} "
  cd /data/yoozoogame/code/zone${zoneID}/${serviceName}${serviceID}
  ./input_etcd.sh
  cd -
  "
fi

source ../../../script/start.sh ${zoneID} ${serviceName} ${serviceID} ${IP} ${PORT} ${USER}