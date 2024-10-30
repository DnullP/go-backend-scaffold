// generators/generator.go
package generators

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/emicklei/proto"
)

type ServiceInfo struct {
	PackageName string
	ServiceName string
	Methods     []MethodInfo
}

type MethodInfo struct {
	MethodName   string
	RequestType  string
	ResponseType string
}

func ServicesGen() {
	protoDir := "./proto"
	protoFilePaths := make([]string, 0)

	err := filepath.Walk(protoDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".proto" {
			protoFilePaths = append(protoFilePaths, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	for _, protoFilePath := range protoFilePaths {
		ServicesGenFile(protoFilePath)
	}
}

// ServicesGen 解析 .proto 文件并生成服务代码
func ServicesGenFile(protoFilePath string) {

	reader, err := os.Open(protoFilePath)
	if err != nil {
		log.Fatalf("无法打开proto文件: %v", err)
	}
	defer reader.Close()

	parser := proto.NewParser(reader)
	definition, err := parser.Parse()
	if err != nil {
		log.Fatalf("解析proto文件失败: %v", err)
	}

	var services []ServiceInfo

	proto.Walk(
		definition,
		proto.WithService(func(s *proto.Service) {
			service := ServiceInfo{
				PackageName: "generated", // 可根据需要动态设置
				ServiceName: s.Name,
				Methods:     []MethodInfo{},
			}

			for _, element := range s.Elements {
				if rpc, ok := element.(*proto.RPC); ok {
					method := MethodInfo{
						MethodName:   rpc.Name,
						RequestType:  rpc.RequestType,
						ResponseType: rpc.ReturnsType,
					}
					service.Methods = append(service.Methods, method)
				}
			}

			services = append(services, service)
		}),
	)

	// 对每个服务生成代码
	for _, service := range services {
		generateServiceCode(service)
		userHandlerCode(service)
		clientCodeGen(service)
		//如有新的服务相关代码, 在此添加
	}

	fmt.Println("所有服务代码生成完成")
}

// generateServiceCode 根据 ServiceInfo 生成对应的 Go 代码
func generateServiceCode(service ServiceInfo) {
	// 定义模板
	tmpl := `
package {{.PackageName}}

import (
    "context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	"go-backend-scaffold/config"
	pb "go-backend-scaffold/proto"
	"go-backend-scaffold/services/discovery"
	"go-backend-scaffold/trace"

	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)


// {{.ServiceName}}Server 是服务接口，用户需要实现这个接口
type {{.ServiceName}}Server interface {
    {{range .Methods}}
    {{.MethodName}}(context.Context, *pb.{{.RequestType}}) (*pb.{{.ResponseType}}, error)
    {{end}}
}

// server 是 gRPC 服务器的实现
type {{.ServiceName}}Implement struct {
    pb.Unimplemented{{.ServiceName}}Server
    Handler {{.ServiceName}}Server
}

{{range .Methods}}
// {{.MethodName}} 实现了 {{$.ServiceName}} 的 {{.MethodName}} 方法
func (s *{{$.ServiceName}}Implement) {{.MethodName}}(ctx context.Context, req *pb.{{.RequestType}}) (*pb.{{.ResponseType}}, error) {
    return s.Handler.{{.MethodName}}(ctx, req)
}
{{end}}


// Start{{.ServiceName}}Server 启动 gRPC 服务器
func Start{{.ServiceName}}Server(ctx context.Context, port string, handler {{.ServiceName}}Server) error {
	//设置监听端口
    lis, err := net.Listen("tcp", ":" + port)
    if err != nil {
        return fmt.Errorf("failed to listen: %v", err)
    }

	trace.SetTraceProvider("{{.ServiceName}}")

	//创建grpc服务器
    grpcServer := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))
    pb.Register{{.ServiceName}}Server(grpcServer, &{{.ServiceName}}Implement{Handler: handler})

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
	err = discovery.RegisterService("{{.ServiceName}}", serverId, iport)
	if err != nil {
		log.Printf("服务注册失败")
		panic(err)
	}
	log.Printf("服务 %s 注册成功", serverId)

	//注销服务
	consulConfig := api.DefaultConfig()
	consulConfig.Address = config.Consul.Address + ":8500"
	client, _ := api.NewClient(consulConfig)
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
		grpcServer.GracefulStop()
	}()

    log.Printf("gRPC server {{.ServiceName}} is running on port :%s", port)
    return grpcServer.Serve(lis)
}
`

	// 解析模板
	t, err := template.New("service").Parse(tmpl)
	if err != nil {
		log.Fatalf("解析模板失败: %v", err)
	}

	// 创建输出目录和文件
	outputDir := "services/" + service.PackageName
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, os.ModePerm)
	}

	outputFile := filepath.Join(outputDir, fmt.Sprintf("generated_%s_server.go", strings.ToLower(service.ServiceName)))
	f, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("创建文件失败: %v", err)
	}
	defer f.Close()

	// 执行模板
	err = t.Execute(f, service)
	if err != nil {
		log.Fatalf("执行模板失败: %v", err)
	}

	fmt.Println("服务代码生成成功:", outputFile)
}

func userHandlerCode(service ServiceInfo) {
	tmpl := `
package services

import (
	"context"
	pb "go-backend-scaffold/proto"
	"go-backend-scaffold/services/generated"
)

type {{.ServiceName}} struct {
	generated.{{.ServiceName}}Server
}

{{range .Methods}}
func (u {{$.ServiceName}}) {{.MethodName}}(ctx context.Context, req *pb.{{.RequestType}}) (*pb.{{.ResponseType}}, error) {
	//complete this function
	return &pb.{{.ResponseType}}{}, nil
}
{{end}}

func Start{{.ServiceName}}Server(ctx context.Context, port string) error {
	return generated.Start{{.ServiceName}}Server(ctx, port, &{{.ServiceName}}{})
}
`
	t, err := template.New("userHandler").Parse(tmpl)
	if err != nil {
		log.Fatalf("解析模板失败: %v", err)
	}

	// 创建输出目录和文件
	outputDir := "services/"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, os.ModePerm)
	}

	outputFile := filepath.Join(outputDir, fmt.Sprintf("%s.go", strings.ToLower(service.ServiceName)))
	f, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("创建文件失败: %v", err)
	}
	defer f.Close()

	// 执行模板
	err = t.Execute(f, service)
	if err != nil {
		log.Fatalf("执行模板失败: %v", err)
	}

	fmt.Println("服务代码生成成功:", outputFile)
}

func clientCodeGen(service ServiceInfo) {
	tmpl := `
package client

import (
	"context"
	pb "go-backend-scaffold/proto"
	"go-backend-scaffold/services/discovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

)


{{range .Methods}}
func {{.MethodName}}(ctx context.Context, req *pb.{{.RequestType}}) (*pb.{{.ResponseType}}, error) {
	serviceAddress := discovery.GetService("{{$.ServiceName}}")

	conn, err := grpc.NewClient(serviceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()),
				 				grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := pb.New{{$.ServiceName}}Client(conn)

	res, err := c.{{.MethodName}}(ctx, req)
	if err != nil {
		panic(err)
	}
	return res, nil
}
{{end}}
`
	t, err := template.New("ClientCode").Parse(tmpl)
	if err != nil {
		log.Fatalf("解析模板失败: %v", err)
	}

	// 创建输出目录和文件
	outputDir := "services/client"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, os.ModePerm)
	}

	outputFile := filepath.Join(outputDir, fmt.Sprintf("%sClient.go", strings.ToLower(service.ServiceName)))
	f, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("创建文件失败: %v", err)
	}
	defer f.Close()

	// 执行模板
	err = t.Execute(f, service)
	if err != nil {
		log.Fatalf("执行模板失败: %v", err)
	}

	fmt.Println("服务代码生成成功:", outputFile)
}
