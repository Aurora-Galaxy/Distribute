package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
)

func RegisterService(r Registration) error {
	heartbeatURL, err := url.Parse(r.HeartbeatURL)
	if err != nil {
		return err
	}
	http.HandleFunc(heartbeatURL.Path, func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})
	serviceUpdateURL, err := url.Parse(r.ServiceUpdateURL)
	if err != nil {
		return err
	}
	http.Handle(serviceUpdateURL.Path, serviceUpdateHandler{})
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err = enc.Encode(r)
	if err != nil {
		return err
	}
	res, err := http.Post(ServicesURL, "application/json", buf)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to register service. "+
			"Registry service responsed with code %v", res.StatusCode)
	}
	return nil
}

type serviceUpdateHandler struct{}

func (suh serviceUpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	dec := json.NewDecoder(r.Body)
	var p patch
	err := dec.Decode(&p)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Printf("Updated receied [Add : %v] [Remove : %v]\n", p.Added, p.Removed)
	prov.Update(p)
}

func ShutDownService(url string) error {
	// http包没有提供 delete 方法，可以自己构建请求
	request, err := http.NewRequest(http.MethodDelete, ServicesURL, bytes.NewBuffer([]byte(url)))
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", "text/plain")
	// 发送构造的请求
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to deregister service. "+
			"Registry service responded with code %v", res.StatusCode)
	}
	return nil
}

// 服务提供方
type providers struct {
	// 每个服务可能有多个 url
	services map[ServiceName][]string
	mutex    *sync.RWMutex
}

// 包内 全局服务提供
var prov = providers{
	services: make(map[ServiceName][]string),
	mutex:    new(sync.RWMutex),
}

func (p *providers) Update(pat patch) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// 增加服务提供方
	for _, patchEntry := range pat.Added {
		if _, ok := p.services[patchEntry.Name]; !ok {
			p.services[patchEntry.Name] = make([]string, 0)
		}
		p.services[patchEntry.Name] = append(p.services[patchEntry.Name], patchEntry.URL)
	}

	// 删除服务提供方
	for _, patchEntry := range pat.Removed {
		if provideUrls, ok := p.services[patchEntry.Name]; ok {
			for i := range provideUrls {
				if provideUrls[i] == patchEntry.URL {
					p.services[patchEntry.Name] = append(provideUrls[:i], provideUrls[i+1:]...)
				}
			}
		}
	}
}

// 代码较简单，每个服务只对应一个 url，所以此处只返回一个，而不是[]string
// 根据服务名称获取其对应的 url
func (p providers) get(name ServiceName) (string, error) {
	providerURLs, ok := p.services[name]
	if !ok {
		return "", fmt.Errorf("No providers available for service %v", name)
	}
	index := int(rand.Float32() * float32(len(providerURLs)))
	return providerURLs[index], nil
}

/**
 * GetProvider
 * @Description: 根据服务名获取其对应的URL
 * @param name
 * @return string url
 * @return error
 */
func GetProvider(name ServiceName) (string, error) {
	return prov.get(name)
}
