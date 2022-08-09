package discovery

import "github.com/leewckk/go-kit-micro-service/discovery/consul"

type Publisher interface {
	Register(instanceId, instanceName string, instanceIP, healthAPI string, instancePort int) error
	UnRegister() error
}

func NewPublisher(host string, port int, regtype string) Publisher {
	return consul.NewPublish(host, port, regtype)
}
