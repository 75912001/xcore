#!/bin/bash

#${1} ${zoneID}
#${2} ${serviceName}
#${3} ${serviceID}
#${4} #${IP}
#${5} ${PORT}
#${6} ${USER}

#echo "start ${name} service ..."

output=$(ssh ${6}@${4} -p ${5} "
cd /data/yoozoogame/admin
/bin/bash start.sh ${1} ${2} ${3}
cd - >> /dev/null 2>&1
sleep 1
echo "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++"
ps -ef | grep -v grep | grep ${1}.${2}
echo "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++"
")

#服务是否正常启动
if echo "${output}" | grep -q "${2}"; then
  echo -e "start service \033[32m [success] \033[0m ${1}.${2}.${3} "
else
  echo -e "start service \033[31m [failed] \033[0m ${1}.${2}.${3} "
  exit 1
fi