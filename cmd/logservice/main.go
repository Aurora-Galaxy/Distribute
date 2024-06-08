package main

import (
	"Distribute/log"
	"Distribute/registry"
	"Distribute/service"
	"context"
	"fmt"
	stlog "log"
)

func main() {
	log.Run("./distribute.log")
	host, port := "localhost", "4000"
	serviceAddress := fmt.Sprintf("http://%s:%s", host, port)
	r := registry.Registration{
		ServiceName:      registry.LogService,
		ServiceURL:       serviceAddress,
		RequiredServices: make([]registry.ServiceName, 0),
		ServiceUpdateURL: serviceAddress + "/services",
		HeartbeatURL:     serviceAddress + "/heartbeat",
	}
	ctx, err := service.Start(
		context.Background(),
		host,
		port,
		r,
		log.RegisterHandlers)
	if err != nil {
		// 本身的日志服务启动出错，使用标准库写入日志
		stlog.Fatalln(err)
	}
	<-ctx.Done()
	fmt.Println("Shutting down log service")
}
