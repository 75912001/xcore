#!/bin/bash

groupID=$(basename $(dirname $(dirname "$PWD")))
serviceName=$(basename $(dirname "$PWD"))
serviceID=$(basename "$PWD")
