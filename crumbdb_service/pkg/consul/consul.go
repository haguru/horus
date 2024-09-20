package consul

import (
	"fmt"
	"os"

	"github.com/haguru/horus/crumbdb/config"
	consulapi "github.com/hashicorp/consul/api"
)

const (
	CHECK_INTERVAL   = "5s"
	CHECK_TIMEOUT    = "30s"
	DEREGISTER_AFTER = "10s"
)

type Consul struct {
	Config *consulapi.Config
	client *consulapi.Client
}

func NewConsul(config *config.Consul) (*Consul, error) {
	address := fmt.Sprintf("%v:%v", config.Host, config.Port)
	consulConfig := &consulapi.Config{
		Address: address,
		// Transport: cleanhttp.DefaultPooledTransport(),
	}

	client, err := consulapi.NewClient(consulConfig)
	if err != nil {
		return nil, err
	}

	return &Consul{
		client: client,
		Config: consulConfig,
	}, nil
}

func (c *Consul) RegisterService(serviceName string, port int) error {
	address, err := os.Hostname()
	if err != nil {
		return err
	}

	registration := &consulapi.AgentServiceRegistration{
		ID:      serviceName,
		Name:    serviceName,
		Port:    port,
		Address: address,

		Check: &consulapi.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%v:%v/%v", address, port, serviceName),
			Interval:                       CHECK_INTERVAL,
			Timeout:                        CHECK_TIMEOUT,
			DeregisterCriticalServiceAfter: DEREGISTER_AFTER,
		},
	}

	err = c.client.Agent().ServiceRegister(registration)
	if err != nil {
		return err
	}

	return nil
}
