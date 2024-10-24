package discovery

import (
	"fmt"
	"go-backend-scaffold/config"
	"log"

	"github.com/hashicorp/consul/api"
)

func GetService(name string) *api.AgentService {
	// 创建默认的配置
	consulConfig := api.DefaultConfig()
	consulConfig.Address = config.Consul.Address + ":8500"

	// 初始化 Consul 客户端
	client, err := api.NewClient(consulConfig)
	if err != nil {
		log.Fatalf("无法创建 Consul 客户端: %v", err)
	}

	// 使用 Agent API 获取已注册的服务
	services, err := client.Agent().Services()
	if err != nil {
		log.Fatalf("无法获取服务列表: %v", err)
	}

	fmt.Println("已注册的服务列表:")
	for serviceName, service := range services {
		fmt.Printf("服务名称: %s\n", serviceName)
		fmt.Printf("服务ID: %s\n", service.ID)
		fmt.Printf("服务地址: %s\n", service.Address)
		fmt.Printf("服务端口: %d\n", service.Port)
		fmt.Printf("服务标签: %v\n", service.Tags)
		fmt.Println("---------------------------")
	}

	// 如果您想获取 Consul 目录中的所有服务，可以使用 Catalog API
	catalogServices, _, err := client.Catalog().Services(nil)
	if err != nil {
		log.Fatalf("无法获取目录中的服务列表: %v", err)
	}

	fmt.Println("目录中的服务列表:")
	for serviceName, tags := range catalogServices {
		fmt.Printf("服务名称: %s\n", serviceName)
		fmt.Printf("标签: %v\n", tags)
		fmt.Println("---------------------------")
	}
	return nil
}
