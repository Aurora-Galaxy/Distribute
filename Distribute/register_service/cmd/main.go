package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"register_service/registry"
	"syscall"
)

func main() {
	// 心跳检测
	registry.SetHeartbeatService()
	http.Handle("/services", registry.RegistryService{})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var srv http.Server
	srv.Addr = registry.ServerPort

	// 信号处理
	sigChan := make(chan os.Signal, 1)
	// 将 SIGINT，SIGTERM与sigChan关联起来，这些信号发生时，会自动将信号值发送到这个channel
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println(srv.ListenAndServe())
		cancel()
	}()

	go func() {
		fmt.Println("Registry Service started.")
		sig := <-sigChan
		fmt.Println("\nReceived signal:", sig)
		_ = srv.Shutdown(ctx)
		cancel()
	}()
	<-ctx.Done()
	fmt.Println("Shutting down registry service")
}
