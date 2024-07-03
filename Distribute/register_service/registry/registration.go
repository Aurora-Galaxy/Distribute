package registry

type ServiceName string

type Registration struct {
	ServiceName ServiceName
	ServiceURL  string
	// 指定该服务所依赖的服务
	RequiredServices []ServiceName
	// 服务注册中心通过该URL通知当前服务是否存在其需要的服务，服务的更新也会使用该url进行通知
	ServiceUpdateURL string
	// 用于心跳检测的URL
	HeartbeatURL string
}

// 目前存在的服务类型
const (
	LogService     = ServiceName("LogService")
	GradingService = ServiceName("GradingService")
	PortalService  = ServiceName("PortalService")
)

type patchEntry struct {
	Name ServiceName
	URL  string
}

// 每次增加和删除的服务
type patch struct {
	Added   []patchEntry
	Removed []patchEntry
}
