#!/bin/bash

#${1} ${zoneID}
#${2} ${serviceName}
#${3} ${serviceID}
#${4} #${IP}
#${5} ${PORT}
#${6} ${USER}
#${7} ${mainPathFile}
#${8} ${benchJsonPathFile}

currPath=$(cd $(dirname ${BASH_SOURCE[0]}); pwd)
source ${currPath}/base.sh

#echo "deploy ${name} service ..."

CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o ${exeName} ${7}

echo "build is complete"

ssh ${6}@${4} -p ${5} "
mkdir -p /data/yoozoogame/code/zone${1}/${2}${3}
mkdir -p /data/yoozoogame/config${1}
mkdir -p /data/yoozoogame/log
mkdir -p /data/yoozoogame/plugin
"

scp -P ${5} ${exeName} ${6}@${4}:/data/yoozoogame/code/zone${1}/${2}${3}/${exeName}
scp -P ${5} ${8} ${6}@${4}:/data/yoozoogame/code/zone${1}/${2}${3}/bench.json

ssh ${6}@${4} -p ${5} "
\cp -sf /data/yoozoogame/code/zone${1}/${2}${3}/${exeName} /data/yoozoogame/code/zone${1}/${2}${3}/${1}.${2}.${3}.exe
"

echo "deploy service done. ${1}.${2}.${3}"