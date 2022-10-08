package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"web_app/settings"
	"web_app/utils/netutil"
)

type consul struct {
	client *api.Client
}

func (c *consul) Init(name string) error {
	ip, err := netutil.GetOutboundIP()
	if err != nil {
		return err
	}
	port, err := netutil.GetAvailablePort()
	settings.Conf.Port = port
	if err != nil {
		return err
	}
	return c.RegisterService(name, ip.String(), port)
}

// NewConsul 连接至consul服务返回一个consul对象
func NewConsul(cfg *settings.ConsulConfig) (*consul, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	config := api.DefaultConfig()
	config.Address = addr
	c, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &consul{c}, nil
}

// RegisterService 将gRPC服务注册到consul
func (c *consul) RegisterService(serviceName string, ip string, port int) error {
	settings.Conf.UUID = uuid.NewV4().String()
	url := fmt.Sprintf("%s:%d", ip, port)

	zap.L().Debug("注册中心地址：%s", zap.String("url", url))

	check := &api.AgentServiceCheck{
		TCP:                            url, // 这里一定是外部可以访问的地址
		Timeout:                        "10s",
		Interval:                       "10s",
		DeregisterCriticalServiceAfter: "20s",
	}

	// grpc 健康检查
	/*check := &api.AgentServiceCheck{
		GRPC:                           url, // 这里一定是外部可以访问的地址
		Timeout:                        "10s",
		Interval:                       "10s",
		DeregisterCriticalServiceAfter: "20s",
	}*/
	srv := &api.AgentServiceRegistration{
		ID:      settings.Conf.UUID, // 服务唯一ID
		Name:    serviceName,        // 服务名称
		Tags:    []string{"v1.0.0"}, // 为服务打标签
		Address: ip,
		Port:    port,
		Check:   check,
	}
	return c.client.Agent().ServiceRegister(srv)
}

// Deregister 注销服务
func (c *consul) Deregister(serviceID string) error {
	return c.client.Agent().ServiceDeregister(serviceID)
}

// ListService 服务发现
func (c *consul) ListService() (map[string]*api.AgentService, error) {
	// c.client.Agent().Service("hello-127.0.0.1-8972")
	return c.client.Agent().ServicesWithFilter("Service==`hello`")
}
