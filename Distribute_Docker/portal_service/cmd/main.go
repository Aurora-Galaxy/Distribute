package main

import (
	"portal_service/log"
	"portal_service/portal"
	"portal_service/registry"
	"portal_service/service"
	"context"
	"fmt"
	stlog "log"
	// "os"
)

func main() {
	// dir, _:= os.Getwd()
    // fmt.Println("当前工作目录:", dir)
	err := portal.ImportTemplates()
	if err != nil {
		stlog.Fatalln(err)
	}
	host, port := "portal_service", "10000"
	serviceAddress := fmt.Sprintf("http://%s:%s", host, port)
	r := registry.Registration{
		ServiceName: registry.PortalService,
		ServiceURL:  serviceAddress,
		RequiredServices: []registry.ServiceName{
			registry.GradingService,
			registry.LogService,
		},
		ServiceUpdateURL: serviceAddress + "/services",
		HeartbeatURL:     serviceAddress + "/heartbeat",
	}
	ctx, err := service.Start(
		context.Background(),
		host,
		port,
		r,
		portal.RegisterHandlers)
	if err != nil {
		stlog.Fatalln(err)
	}
	if logProvider, err := registry.GetProvider(registry.LogService); err == nil {
		fmt.Printf("Log Service found at : %s\n", logProvider)
		log.SetClientLogger(logProvider, r.ServiceName)
	}
	<-ctx.Done()
	fmt.Println("Shutting down Portal service")
}
