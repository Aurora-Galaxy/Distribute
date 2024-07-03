package log

import (
	"grading_service/registry"
	"bytes"
	"fmt"
	stlog "log"
	"net/http"
)

func SetClientLogger(ServiceURL string, clientService registry.ServiceName) {
	stlog.SetPrefix(fmt.Sprintf("[%v] - ", clientService))
	// 不设置 flag
	stlog.SetFlags(0)
	stlog.SetOutput(clientLogger{url: ServiceURL})
}

type clientLogger struct {
	url string
}

func (cl clientLogger) Write(data []byte) (int, error) {
	b := bytes.NewBuffer(data)
	res, err := http.Post(cl.url+"/log", "text/plain", b)
	if err != nil {
		return 0, err
	}
	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Failed to send log message , Service Responded wrong")
	}
	return len(data), nil
}
