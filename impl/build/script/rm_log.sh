#!/bin/bash

#${1} ${zoneID}
#${2} ${serviceName}
#${3} ${serviceID}
#${4} #${IP}
#${5} ${PORT}
#${6} ${USER}

ssh ${6}@${4} -p ${5} "
cd /data/yoozoogame/log
rm ${1}-${2}-${3}-* -rf
cd - >> /dev/null 2>&1
"

echo "rm log done. ${1}.${2}.${3}"