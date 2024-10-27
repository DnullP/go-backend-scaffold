# Go-Backend-Scaffold

这是一个用于Go开发的后端服务治理脚手架, 提供了以下功能:
- 服务发现
- 链路追踪
- 接口模板生成

## 使用该工具

为了使用该框架, 请首先安装用于编译go代码的protoc相关工具, 具体请参考:https://grpc.io/docs/languages/go/quickstart/

其中包含了grpc接口编译所需前置工具的全部安装说明.

在完成前置安装后, 安装项目所需模块:
```Go
go get
```

此脚手架通过在`proto`目录下编写`protobuf`文件, 能够自动生成对应的接口和基础代码, 其中集成了自动的服务注册与发现, 链路追踪等.

在完成`proto`的编写后, 通过`go generate`生成代码, 然后你只需要实现在`services`文件夹下的服务业务代码即可.

以下是一个生成的接口示例:
```Go

package services

import (
	"context"
	pb "go-backend-scaffold/proto"
	"go-backend-scaffold/services/generated"
)

type UserService struct {
	generated.UserServiceServer
}


func (u UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	//complete this function
	return &pb.GetUserResponse{}, nil
}

func (u UserService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	//complete this function
	return &pb.LoginResponse{}, nil
}


func StartUserServiceServer(ctx context.Context, port string) error {
	return generated.StartUserServiceServer(ctx, port, &UserService{})
}
```

---

## 外部服务

目前服务注册与发现使用[Consul](https://www.consul.io/)完成.

链路追踪使用[Jaeger All in one](https://www.jaegertracing.io/)完成.

请在`config.yaml`中配置对应的依赖服务的地址.

