package eureka

import (
	"fmt"
	"strings"
	"sync"

	"github.com/chengjianxi/goc/haoxin/balancer"
	eureka_client "github.com/xuanbo/eureka-client"
)

type Eureka struct {
	eureka     *eureka_client.Client
	serverPool *balancer.NamesServerPool
	mutex      sync.RWMutex
}

// 启动服务注册发现
func Start(zone string, appName string, port int) *Eureka {
	// 服务注册 github.com/xuanbo/eureka-client
	// 创建 eureka client
	// zone := fmt.Sprintf("http://%s:%d/eureka/", c.Eureka.Host, c.Eureka.Port)
	eureka := eureka_client.NewClient(&eureka_client.Config{
		DefaultZone:           zone,
		App:                   appName,
		Port:                  port,
		RenewalIntervalInSecs: 30,
		DurationInSecs:        30,
		InstanceID:            fmt.Sprintf("%s:%s:%d", strings.ToLower(appName), eureka_client.GetLocalIP(), port),
	})
	// 启动 eureka client, register、heartbeat、refresh
	eureka.Start()
	return &Eureka{
		eureka:     eureka,
		serverPool: balancer.NewNamesServerPool(),
	}
}

func (c *Eureka) GetServiceAddress(name string) string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	instances := c.eureka.GetApplicationInstance(name)
	urls := make([]string, 0)
	for _, instance := range instances {
		if instance.Status == "UP" {
			urls = append(urls, instance.HomePageURL)
		}
	}
	c.serverPool.SetServerAddrs(name, urls)
	return c.serverPool.GetServerAddr(name)
}
