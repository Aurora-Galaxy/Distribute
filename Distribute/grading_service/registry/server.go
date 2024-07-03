package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	ServerPort = ":3000"
	// 通过该地址可以查看哪些服务已经在此注册
	ServicesURL = "http://register_service" + ServerPort + "/services"
)

type registry struct {
	// 保存已经注册的服务
	registrations []Registration
	mutex         *sync.RWMutex
}

// 全局 registry 实例,用于管理所有注册的服务
var reg = registry{
	registrations: make([]Registration, 0),
	mutex:         new(sync.RWMutex),
}

// 添加服务，在添加该服务时直接将该服务所依赖的服务给他
func (r *registry) add(reg Registration) error {
	r.mutex.Lock()
	r.registrations = append(r.registrations, reg)
	r.mutex.Unlock()
	// 在注册服务时，通知需要该服务的service
	r.notify(patch{
		Added: []patchEntry{
			{
				Name: reg.ServiceName,
				URL:  reg.ServiceURL,
			},
		},
	})
	err := r.sendRequiredServices(reg)
	return err
}

func (r registry) notify(fullPatch patch) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	for _, registration := range r.registrations {
		go func(reg Registration) {
			for _, requireServiceName := range reg.RequiredServices {
				p := patch{
					Added:   []patchEntry{},
					Removed: []patchEntry{},
				}
				sendUpdate := false
				for _, added := range fullPatch.Added {
					if added.Name == requireServiceName {
						p.Added = append(p.Added, added)
						sendUpdate = true
					}
				}
				for _, removed := range fullPatch.Removed {
					if removed.Name == requireServiceName {
						p.Removed = append(p.Removed, removed)
						sendUpdate = true
					}
				}
				if sendUpdate {
					err := r.sendPatch(p, reg.ServiceUpdateURL)
					if err != nil {
						log.Println(err)
						return
					}
				}
			}
		}(registration)
	}
}

func (r *registry) sendRequiredServices(reg Registration) error {
	var p patch
	// 查找是否有当前服务需要的服务
	for _, existService := range r.registrations {
		for _, needService := range reg.RequiredServices {
			if existService.ServiceName == needService {
				// 匹配到后，将其加入保存需要添加服务的结构体中
				p.Added = append(p.Added, patchEntry{
					Name: existService.ServiceName,
					URL:  existService.ServiceURL,
				})
			}
		}
	}
	err := r.sendPatch(p, reg.ServiceUpdateURL)
	if err != nil {
		return err
	}
	return nil
}

func (r *registry) sendPatch(p patch, url string) error {
	jsonData, err := json.Marshal(p)
	if err != nil {
		return err
	}
	_, err = http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	return nil
}

// 取消服务
func (r *registry) remove(url string) error {
	for index, registration := range r.registrations {
		if registration.ServiceURL == url {
			r.mutex.Lock()
			r.registrations = append(r.registrations[:index], r.registrations[index+1:]...)
			r.mutex.Unlock()
			r.notify(patch{
				Removed: []patchEntry{
					{
						Name: registration.ServiceName,
						URL:  registration.ServiceUpdateURL,
					},
				},
			})
			return nil
		}
	}
	return fmt.Errorf("Service at URL %s not found", url)
}

/**
 * HeartBeat
 * @Description: 心跳检测机制
 * @receiver r
 * @param freq 进行心跳检测的间隔
 */
func (r *registry) HeartBeat(freq time.Duration) {
	for {
		var wg sync.WaitGroup
		for _, registration := range r.registrations {
			wg.Add(1)
			go func(reg Registration) {
				defer wg.Done()
				success := true
				// 心跳检查失败，重试 3 次
				for attempts := 0; attempts < 3; attempts++ {
					res, err := http.Get(reg.HeartbeatURL)
					if err != nil {
						log.Println(err)
					} else if res.StatusCode == http.StatusOK {
						log.Printf("Heartbeat check passed for %v\n", reg.ServiceName)
						if !success {
							_ = r.add(reg)
						}
						break
					}
					log.Printf("Heartbeat check failed for %v\n", reg.ServiceName)
					if success {
						success = false
						_ = r.remove(reg.ServiceURL)
					}
					time.Sleep(1 * time.Second)
				}
			}(registration)
		}
		wg.Wait()
		time.Sleep(freq)
	}
}

var once sync.Once

func SetHeartbeatService() {
	once.Do(func() {
		go reg.HeartBeat(3 * time.Second)
	})
}

type RegistryService struct{}

func (rs RegistryService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received")
	switch r.Method {
	// post 注册
	case http.MethodPost:
		dec := json.NewDecoder(r.Body)
		var register Registration
		err := dec.Decode(&register)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("Adding service: %v with URL:%v \n", register.ServiceName, register.ServiceURL)
		// 添加服务
		err = reg.add(register)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	//	delete 取消服务
	case http.MethodDelete:
		payload, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		url := string(payload)
		log.Printf("Removing service at URL:%s \n", url)
		// 添加服务
		err = reg.remove(url)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
