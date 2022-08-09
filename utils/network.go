package utils

import (
	"net"

	"github.com/sirupsen/logrus"
)

func GetLocalIpAddress() []string {
	list := make([]string, 0, 0)

	/// 获取IP地址
	addrs, err := net.InterfaceAddrs()
	if nil != err {
		logrus.Errorf("error get interface address: %v", err)
		return list
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			list = append(list, ipnet.IP.String())
		}
	}
	return list
}
