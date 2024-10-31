package etcd

import (
	"encoding/json"
	"fmt"
	"path"
	"strconv"
	xlog "xcore/lib/log"
)

// Parse 解析key
func Parse(key string) (msgType string, groupID uint32, serviceName string, serviceID uint32) {
	strServiceID := path.Base(key)
	// strServiceID 转换成 serviceID
	if serviceIDU64, err := strconv.ParseUint(strServiceID, 10, 32); err != nil {
		xlog.PrintfErr("Parse err:%v %v %v", key, strServiceID, err)
		return
	} else {
		serviceID = uint32(serviceIDU64)
	}
	key = path.Dir(key)
	key = path.Dir(key)
	serviceName = path.Base(key)
	key = path.Dir(key)
	strGroupID := path.Base(key)
	// strGroupID 转换成 groupID
	if groupIDU64, err := strconv.ParseUint(strGroupID, 10, 32); err != nil {
		xlog.PrintfErr("Parse err:%v %v %v", key, strGroupID, err)
		return
	} else {
		groupID = uint32(groupIDU64)
	}
	key = path.Dir(key)
	msgType = path.Base(key)
	return msgType, groupID, serviceName, serviceID
}

func GenKey(projectName string, etcdWatchMsgType string, groupID uint32, serviceName string, serviceID uint32) string {
	return fmt.Sprintf("/%v/%v/%v/%v/%v/",
		projectName,
		etcdWatchMsgType,
		groupID,
		serviceName,
		serviceID,
	)
}

func GenPrefixKey(projectName string) string {
	return fmt.Sprintf("/%v/", projectName)
}

func ValueJson2String(valueJson *ValueJson) string {
	bytes, err := json.Marshal(valueJson)
	if err != nil {
		xlog.PrintfErr("Error marshaling ValueJson: %v", err)
		return ""
	}
	return string(bytes)
}

func ValueString2Json(value string) *ValueJson {
	var valueJson ValueJson
	err := json.Unmarshal([]byte(value), &valueJson)
	if err != nil {
		xlog.PrintfErr("Error unmarshaling ValueJson: %v %v", value, err)
		return nil
	}
	return &valueJson
}
