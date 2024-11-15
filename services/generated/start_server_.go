package generated

import (
	"context"
	"fmt"
	"go-backend-scaffold/config"
	"go-backend-scaffold/services/generated/discovery"
	"go-backend-scaffold/trace"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
)

func startServer(ctx context.Context, port string, server *grpc.Server, serviceName string) error {
	//设置监听端口
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	trace.SetTraceProvider(serviceName)

	//创建健康检查服务
	go func() {
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			log.Println("收到健康检查请求")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})
		if err := http.ListenAndServe(":8000", nil); err != nil {
			log.Printf("无法启动健康检查 HTTP 服务器: %v", err)
		}
	}()

	//注册服务
	serverId := uuid.New().String()
	iport, _ := strconv.Atoi(port)
	err = discovery.RegisterService(serviceName, serverId, iport)
	if err != nil {
		log.Printf("服务注册失败")
		panic(err)
	}
	log.Printf("服务 %s 注册成功", serverId)

	//注销服务
	client, _ := api.NewClient(config.Consul)
	defer func() {
		err := discovery.DeregisterService(serverId, client)
		if err != nil {
			log.Printf("服务注销失败: %v", err)
		} else {
			log.Printf("服务 %s 已注销", serverId)
		}
	}()

	//关闭服务
	go func() {
		<-ctx.Done()
		log.Println("接收到关闭信号，正在优雅地停止 gRPC 服务器...")
		server.GracefulStop()
	}()

	log.Printf("gRPC server %s is running on port :%s", serviceName, port)
	return server.Serve(lis)
}
