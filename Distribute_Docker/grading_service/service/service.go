package service

import (
	"grading_service/registry"
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"os"
	"syscall"
)

func Start(ctx context.Context, host, port string, reg registry.Registration, registerHandlers func()) (context.Context, error) {
	registerHandlers()
	ctx = StartService(ctx, reg.ServiceName, host, port)
	err := registry.RegisterService(reg)
	if err != nil {
		return ctx, err
	}
	return ctx, nil
}

func StartService(ctx context.Context, serviceName registry.ServiceName, host, port string) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	var srv http.Server
	// 本地运行，只需指定端口号
	srv.Addr = ":" + port
	// 信号处理
	sigChan := make(chan os.Signal, 1)
	// 将 SIGINT，SIGTERM与sigChan关联起来，这些信号发生时，会自动将信号值发送到这个channel
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		// 协程 监听服务端口，出现错误时打印错误并发出取消信号
		log.Println(srv.ListenAndServe())
		// 监听发生错误时，注册请求已经发送，所以需要取消注册
		err := registry.ShutDownService(fmt.Sprintf("http://%s:%s", host, port))
		if err != nil {
			log.Println(err)
		}
		cancel()
	}()

	go func() {
		fmt.Println("Grading Service started.")
		sig := <-sigChan
		fmt.Println("\nReceived signal:", sig)
		// 用户取消服务也需要取消注册
		err := registry.ShutDownService(fmt.Sprintf("http://%s:%s", host, port))
		if err != nil {
			log.Println(err)
		}
		_ = srv.Shutdown(ctx)
		cancel()
	}()
	return ctx
}


