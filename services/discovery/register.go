// register.go
package discovery

import (
	"fmt"
	"go-backend-scaffold/config"
	"net"

	"github.com/hashicorp/consul/api"
)

func RegisterService(serviceName string, serviceID string, servicePort int) error {
	client, err := api.NewClient(config.Consul)
	if err != nil {
		return err
	}

	// 获取本机 IP
	hostIP, err := getLocalIP()
	if err != nil {
		return err
	}

	// 定义服务检查
	check := &api.AgentServiceCheck{
		HTTP:     fmt.Sprintf("http://%s:%d/health", hostIP, 8000),
		Interval: "10s",
		Timeout:  "1s",
	}

	// 定义服务
	registration := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Address: hostIP,
		Port:    servicePort,
		Check:   check,
	}

	// 注册服务
	return client.Agent().ServiceRegister(registration)
}

// 获取本地 IP 地址
func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("无法获取本地 IP 地址")
}

// 服务注销函数
func DeregisterService(serviceID string, client *api.Client) error {
	return client.Agent().ServiceDeregister(serviceID)
}
