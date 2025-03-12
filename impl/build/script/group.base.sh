#!/bin/bash

exeName="project.server.exe"

serviceNameLogic="logic"
serviceNameGateway="gateway"
serviceNameLogin="login"

#服务名数组
serviceArr=(${serviceNameLogic} ${serviceNameGateway} ${serviceNameLogin})

#开启服务顺序
startServiceArr=(${serviceNameLogic} ${serviceNameGateway} ${serviceNameLogin})

#关闭服务顺序
stopServiceArr=(${serviceNameLogin} ${serviceNameGateway} ${serviceNameLogic})