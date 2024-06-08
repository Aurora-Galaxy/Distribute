package service

import (
	"Distribute/registry"
	"context"
	"fmt"
	"log"
	"net/http"
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
		// 用户可以输入任意内容，然后停止服务
		fmt.Printf("%v started. Press any key to stop. \n", serviceName)
		var s string
		fmt.Scanln(&s)
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
