package main

import (
	"fmt"
	"picasso/cmd/app"
	"picasso/config"
	_ "picasso/config"
)

func main() {
	cmd := app.NewApiServerCommand()
	fmt.Println("config.SysYamlconfig.Server.Name = ", config.SysYamlconfig.Server.Name)
	cmd.Execute()
}

// go run cmd/apiserver.go apiserver --port=8888
