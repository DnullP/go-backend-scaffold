package discovery

import (
	"fmt"
	"go-backend-scaffold/config"
	"log"
	"strconv"

	"github.com/hashicorp/consul/api"
)

func GetService(name string) string {
	// 初始化 Consul 客户端
	client, err := api.NewClient(config.Consul)
	if err != nil {
		log.Fatalf("无法创建 Consul 客户端: %v", err)
	}

	// 使用 Agent API 获取已注册的服务
	services, _, err := client.Health().Service(name, "", true, &api.QueryOptions{})
	if err != nil {
		log.Fatalf("无法获取服务列表: %v", err)
	}

	fmt.Println("已注册的服务列表:")
	//TODO: 服务选择策略
	for _, service := range services {
		fmt.Printf("服务ID: %s\n", service.Service.ID)
		fmt.Printf("服务地址: %s\n", service.Service.Address)
		fmt.Printf("服务端口: %d\n", service.Service.Port)
		fmt.Println("---------------------------")
		return service.Service.Address + ":" + strconv.Itoa(service.Service.Port)
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
	return "no service"
}
