package http

import (
	"fmt"
	"strings"
)

func GetHttpAPI(version, serverName, moduleName, serviceName, methodName string) string {
	api := fmt.Sprintf("/%v/%v/%v/%v/%v", version, serverName, moduleName, serviceName, methodName)
	return strings.ToLower(api)
}
