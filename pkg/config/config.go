package config

const MessageDelim = '\n'
const messageDelimStr = "\n"

func ParseMessageRead(str string) string {
	if len(str) >= 2 && str[len(str)-2] == '\r' {
		return str[:len(str)-2]
	}
	if len(str) == 0 {
		return ""
	}
	return str[:len(str)-1]
}

func ParseWriteMessage(str string) string {
	return str + messageDelimStr
}

// Not const and different variables to be able to mock it
var ClientPortConnect = 8080
var ClientPortListen = 8080

const ServerHeartbeatPort = 8081
const OK = "OK"
