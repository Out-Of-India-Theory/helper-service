package main

import (
	"context"
	"fmt"
	"github.com/Out-Of-India-Theory/helper-service/config"
	"github.com/Out-Of-India-Theory/helper-service/server"
	"github.com/Out-Of-India-Theory/oit-go-commons/app"
)

func main() {
	configuration := config.GetConfig()
	ctx := context.Background()
	App, err := app.NewApp(ctx, configuration.ServerConfig)
	if err != nil {
		panic(fmt.Sprintf("Unable to initialize the app : %v", err))
	}
	server.InitServer(ctx, App, configuration)
}
