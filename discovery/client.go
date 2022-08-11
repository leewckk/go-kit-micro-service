package discovery

import (
	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
)

func NewClient(apiUri string) consul.Client {

	cfg := api.DefaultConfig()
	cfg.Address = apiUri

	client, err := api.NewClient(cfg)
	if nil != err {
		panic(err)
	}
	return consul.NewClient(client)
}
