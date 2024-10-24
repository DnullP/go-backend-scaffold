package main

import (
	"go-backend-scaffold/generator/generators"
)

func main() {
	generators.ProtoGen()
	generators.ServicesGen()
}
