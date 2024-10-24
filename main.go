// cmd/server/main.go
package main

import (
	"fmt"
	"go-backend-scaffold/config"
)

func main() {
	config.LoadConfig()
	fmt.Print()
}
