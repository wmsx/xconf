package config

import (
	"container/list"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/micro/go-micro/util/log"
	"github.com/wmsx/xconf/agent-api/config/cache"
	"github.com/wmsx/xconf/proto/config"
)

type Config struct {
	sync.RWMutex
	configServiceClient config.ConfigService
	cache               cache.Cache
	watchers            map[string]*list.List
}

var defaultConfig *Config

func Init(client config.ConfigService, cacheSize int) {
	defaultConfig = newConfig(client, cacheSize)
	go defaultConfig.run()
}

func ReadConfig(appName, clusterName, namespaceName string) (*config.ConfigResponse, error) {
	return defaultConfig.ReadConfig(appName, clusterName, namespaceName)
}

func Watch(appName, clusterName, namespaceName string) *Watcher {
	return defaultConfig.Watch(appName, clusterName, namespaceName)
}

func newConfig(client config.ConfigService, cacheSize int) *Config {
	return &Config{
		configServiceClient: client,
		cache:               cache.New(cacheSize),
		watchers:            make(map[string]*list.List),
	}
}

func (c *Config) run() {
	for {
		c.cache.Clear()

		stream, err := c.configServiceClient.Watch(context.Background(), &config.Request{})
		if err != nil {
			log.Errorf("stream watch error : %s (will sleep 1 second)", err.Error())
			time.Sleep(time.Second)
			continue
		}
		log.Info("config watcher reconnected")

		for {
			value, err := stream.Recv()
			if err != nil {
				log.Info("stream recv error :", err)
				break
			}

			log.Info("get release config :", value)

			watchers := c.copyWatchers(getKey(value.AppName, value.ClusterName, value.NamespaceName))
			for e := watchers.Front(); e != nil; e = e.Next() {
				w := e.Value.(*Watcher)
				select {
				case w.updates <- value:
				default:
				}
			}

			if err := c.cache.Set(value); err != nil {
				log.Error("update cache error:", err)
			}
		}

		log.Info("config watcher reconnecting")
	}
}

func (c *Config) copyWatchers(key string) *list.List {
	watcherList := list.New()
	c.RLock()
	watchers := c.watchers[key]
	if watchers != nil {
		watcherList.PushBackList(watchers)
	}
	c.RUnlock()

	return watcherList
}

func (c *Config) Watch(appName, clusterName, namespaceName string) *Watcher {
	w := &Watcher{
		exit:    make(chan interface{}),
		updates: make(chan *config.ConfigResponse, 1),
	}

	key := getKey(appName, clusterName, namespaceName)
	c.Lock()
	watchers := c.watchers[key]
	if watchers == nil {
		watchers = list.New()
		c.watchers[key] = watchers
	}

	e := watchers.PushBack(w)
	c.Unlock()

	go func() {
		<-w.exit
		c.Lock()
		watchers.Remove(e)
		c.Unlock()
	}()

	return w
}

func (c *Config) ReadConfig(appName, clusterName, namespaceName string) (*config.ConfigResponse, error) {
	reqConf := &config.QueryConfigRequest{
		AppName:       appName,
		ClusterName:   clusterName,
		NamespaceName: namespaceName,
	}

	value, ok := c.cache.Get(reqConf)
	if ok { // 命中缓存
		log.Info("命中缓存")
		return value, nil
	} else {
		log.Info("未能命中缓存")

		conf, err := c.configServiceClient.Read(context.Background(), reqConf)
		if err != nil {
			return nil, err
		}

		// 更新缓存
		if err := c.cache.Set(conf); err != nil {
			log.Error("update cache error:", err)
			return nil, err
		}
		return conf, nil
	}
}

func getKey(appName, clusterName, namespaceName string) string {
	return fmt.Sprintf("%s/%s/%s", appName, clusterName, namespaceName)
}
