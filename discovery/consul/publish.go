package consul

import (
	"fmt"

	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
)

const (
	RegisterTypeHTTP = "http"
	RegisterTypeGRPC = "grpc"
)

type Publish struct {
	Host     string
	Port     int
	RegType  string
	Registed map[string]struct{}
	client   consul.Client
}

func NewPublish(host string, port int, regtype string) *Publish {

	pub := &Publish{
		Host:     host,
		Port:     port,
		RegType:  regtype,
		Registed: make(map[string]struct{}),
	}

	if regtype != RegisterTypeGRPC && regtype != RegisterTypeHTTP {
		panic(fmt.Sprintf("unknown register type : %v , support type: [%v, %v]", regtype, RegisterTypeHTTP, RegisterTypeGRPC))
	}

	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%v:%v", host, port)
	apicli, err := api.NewClient(config)
	if nil != err {
		panic(fmt.Sprintf("Error create consul client : %v", err))
	}
	pub.client = consul.NewClient(apicli)
	return pub
}

func (this *Publish) NewRegistration(instanceId, instanceName string, instanceIP, healthAPI string, port int, regtype string) *api.AgentServiceRegistration {

	check := api.AgentServiceCheck{
		DeregisterCriticalServiceAfter: "15s",
		Interval:                       "10s",
		Timeout:                        "1s",
		Notes:                          fmt.Sprintf("%v %v service", instanceName, regtype),
	}

	if regtype == RegisterTypeGRPC {
		check.GRPC = fmt.Sprintf("%v:%v/%v", instanceIP, port, healthAPI)
	} else if regtype == RegisterTypeHTTP {
		check.HTTP = fmt.Sprintf("http://%v:%v%v", instanceIP, port, healthAPI)
	} else {
		panic(fmt.Sprintf("unknown register type : %v , support type: [%v, %v]", regtype, RegisterTypeHTTP, RegisterTypeGRPC))
	}

	registration := api.AgentServiceRegistration{
		ID:      instanceId,
		Name:    instanceName,
		Address: instanceIP,
		Port:    port,
		Check:   &check,
	}
	return &registration
}

func (this *Publish) Register(instanceId, instanceName string, instanceIP, healthAPI string, instancePort int) error {
	registration := this.NewRegistration(instanceId, instanceName, instanceIP, healthAPI, instancePort, this.RegType)
	this.Registed[registration.ID] = struct{}{}
	return this.client.Register(registration)
}

func (this *Publish) UnRegister() error {
	for instanceId, _ := range this.Registed {
		registration := api.AgentServiceRegistration{
			ID: instanceId,
		}
		this.client.Deregister(&registration)
	}
	return nil
}
