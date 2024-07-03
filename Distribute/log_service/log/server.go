package log

import (
	"io"
	stlog "log"
	"net/http"
	"os"
)

// 接收 post 请求，将其内容写入日志文件

var log *stlog.Logger

type fileLog string

// fileLog 写入文件的路径, 该方法目的为实现 io.Writer接口
func (fl fileLog) Write(data []byte) (int, error) {
	f, err := os.OpenFile(string(fl), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return 0, err
	}
	defer func() { _ = f.Close() }()
	return f.Write(data)
}

// 服务启动时，指定固定地址写 log 文件
func Run(dest string) {
	// flag 定义日志属性
	log = stlog.New(fileLog(dest), "[go] - ", stlog.LstdFlags)
}

func RegisterHandlers() {
	http.HandleFunc("/log", func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodPost:
			msg, err := io.ReadAll(request.Body)
			defer func() { _ = request.Body.Close() }()
			if err != nil || len(msg) == 0 {
				writer.WriteHeader(http.StatusBadRequest)
				return
			}
			write(string(msg))
		default:
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	})
}

func write(msg string) {
	log.Printf("%v\n", msg)
}
