#!/bin/bash

genMod=go_out      #google.golang.org/protobuf/cmd/protoc-gen-go
#genMod=gofast_out  #github.com/gogo/protobuf/protoc-gen-gofast

########################################################################################################################
#生成 error 文件
genProtoError(){
  #源 目录文件
  srcDirFile=${1}

  #目标目录
  desDir="./error"
  #目标文件
  desFile="error.go"
  #目标 目录文件
  desDirFile=${desDir}/${desFile}

  #生成目录
  mkdir -p ${desDir}

  echo "// Code generated by gen.sh. DO NOT EDIT." > ${desDirFile}
  # shellcheck disable=SC2129
  echo "package error" >> ${desDirFile}
  echo "import xerror \"xcore/lib/error\"" >> ${desDirFile}
  echo "var (" >> ${desDirFile}

  #临时文件
	tmpFile=${desDirFile}.tmp

  # shellcheck disable=SC2002
  cat "${srcDirFile}" | grep "EC_" | grep -v "//EC_" | grep -v "0x0000;" > ${tmpFile}

	sed -i -r 's/[\s\t]*//g' ${tmpFile}
	sed -i -r 's/\/\// \/\//g' ${tmpFile}
	sed -i -r 's/;/ tagName tagDesc/g' ${tmpFile}

	awk -F " " '{printf("%s = xerror.NewError(%s).WithName(\"%s\").WithDesc(\"%s\") \r\n",$1,$3,$1,$6)}' ${tmpFile} > ${tmpFile}.xxx
	cat ${tmpFile}.xxx > ${tmpFile}

	unlink ${tmpFile}.xxx

	sed -i -r 's/\/\///g' ${tmpFile}
	cat ${tmpFile} >> ${desDirFile}

	unlink ${tmpFile}

  echo ")" >> ${desDirFile}

  #格式化,不显示输出
  go fmt ${desDirFile} > /dev/null
}

#生成CMD
genCMD(){
  packageName=${1}
  #目标目录
  desDir="./${packageName}"
  #目标文件
  desFile="${packageName}.go"
  #目标 目录文件
  desDirFile=${desDir}/${desFile}
  #生成目录
  mkdir -p "${desDir}"
  #echo "gen CMD ... ${packageName}"
  echo "/*Code generated by gen.sh. DO NOT EDIT.*/" > "${desDirFile}"
  echo "package ${packageName}" >>  "${desDirFile}"
  # shellcheck disable=SC2002
  cat "${fileName}" | grep "#" >>  "${desDirFile}"
  # shellcheck disable=SC2086
  sed -i "s/message /const /g"  ${desDirFile}
  # shellcheck disable=SC2086
  sed -i "s/\/\//_CMD uint32 = /g"  ${desDirFile}
  sed -i "s/#/ \/\/ /g"  "${desDirFile}"

  #格式化,不显示输出
  go fmt "${desDirFile}" > /dev/null
}

########################################################################################################################
#error.code
srcDirFile="./proto/error.code.proto"
protoc --${genMod}=. ${srcDirFile}
genProtoError ${srcDirFile}

########################################################################################################################
#login.proto
fileName="./proto/service/login.proto"
protoc --${genMod}=. ${fileName}
genCMD "login"

########################################################################################################################
#gateway.proto
fileName="./proto/service/gateway.proto"
protoc --${genMod}=. ${fileName}
genCMD "gateway"
#modifyPBFile world

########################################################################################################################
#event.proto
fileName="./proto/event.proto"
protoc --${genMod}=. ${fileName}

########################################################################################################################
#struct.proto
fileName="./proto/struct.proto"
protoc --${genMod}=. ${fileName}

echo "生成 protobuf 完成"
#read input

exit 0
