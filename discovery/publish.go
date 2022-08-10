package discovery

import (
	"fmt"
	"os"

	"github.com/leewckk/go-kit-micro-service/discovery/consul"
	"github.com/leewckk/go-kit-micro-service/utils"
)

type Publisher interface {
	Register(instanceId, instanceName string, instanceIP, healthAPI string, instancePort int) error
	UnRegister() error
}

func NewPublisher(host string, port int, regtype string) Publisher {
	return consul.NewPublish(host, port, regtype)
}

func Publish(center string, centerPort int, serviceName string, servicePort int, health string, regtype string) (Publisher, error) {
	pub := NewPublisher(center, centerPort, regtype)
	instance, _ := os.Hostname()
	iplist := utils.GetLocalIpAddress()
	if len(iplist) == 0 {
		panic("failed find valid ip address for register center !!!")
	}
	err := pub.Register(
		fmt.Sprintf("%v.%v.%v", serviceName, instance, regtype),
		fmt.Sprintf("%v.%v", serviceName, regtype),
		iplist[0],
		health,
		servicePort,
	)
	return pub, err
}
