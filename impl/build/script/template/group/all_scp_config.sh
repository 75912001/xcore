#!/bin/bash

source ./base.sh

zoneID=$(basename "$PWD")

scp -P ${PORT} ../config/* ${USER}@${IP}:/data/yoozoogame/config${zoneID}/
echo "scp to ${zoneID} ... done."
