#!/bin/bash

#${1} ${zoneID}
#${2} ${serviceName}
#${3} ${serviceID}
#${4} #${IP}
#${5} ${PORT}
#${6} ${USER}

#echo "stop ${name} server ..."

ssh ${6}@${4} -p ${5} "
cd /data/yoozoogame/admin
/bin/bash stop.sh ${1} ${2} ${3}
cd - >> /dev/null 2>&1
"

echo "stop service done. ${1}.${2}.${3}"
