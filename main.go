package main

import (
	"wallet/app"
	"wallet/config"
)

func main() {
	config := config.GetConfig()
	app := &app.App{}
	app.InitializeAndRun(config, ":2004")
}
