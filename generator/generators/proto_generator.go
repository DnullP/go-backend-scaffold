package generators

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func ProtoGen() {
	// 定义包含 .proto 文件的目录
	protoDir := "./proto"

	// 遍历目录，查找所有 .proto 文件
	err := filepath.Walk(protoDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".proto" {
			// 执行 protoc 命令编译该 .proto 文件
			cmd := exec.Command(
				"protoc",
				"--go_out=.",
				"--go-grpc_out=.",
				"--go-grpc_opt=paths=source_relative",
				path,
			)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			fmt.Printf("正在编译 %s...\n", path)
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("编译 %s 失败: %v", path, err)
			}
			fmt.Printf("成功编译 %s\n", path)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("生成 .proto 文件时出错: %v", err)
	}

	fmt.Println("所有 .proto 文件已成功编译。")
}
