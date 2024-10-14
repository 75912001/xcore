package service

const (
	UnknownMessage uint32 = 0
	LoginMessage   uint32 = 1
	GatewayMessage uint32 = 2
	LogicMessage   uint32 = 3
)

func IsLoginMessage(messageID uint32) bool {
	return 0x10000 <= messageID && messageID <= 0x1ffff
}
func IsGatewayMessage(messageID uint32) bool {
	return 0x20000 <= messageID && messageID <= 0x2ffff
}

func IsLogicMessage(messageID uint32) bool {
	return 0x30000 <= messageID && messageID <= 0x3ffff
}

func GetServiceTypeByMessageID(messageID uint32) uint32 {
	if IsLoginMessage(messageID) {
		return LoginMessage
	} else if IsGatewayMessage(messageID) {
		return GatewayMessage
	} else if IsLogicMessage(messageID) {
		return LogicMessage
	}
	return UnknownMessage
}
