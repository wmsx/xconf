package cache

import (
	"encoding/json"
	"fmt"
	"runtime/debug"

	"github.com/Allenxuxu/toolkit/convert"
	"github.com/coocood/freecache"
	"github.com/micro/go-micro/v2/util/log"
	"github.com/wmsx/xconf/proto/config"
)

var _ Cache = &freeCache{}

type freeCache struct {
	cache *freecache.Cache
}

func newFreeCache(size int) *freeCache {
	c := &freeCache{
		cache: freecache.NewCache(size),
	}
	debug.SetGCPercent(10)

	return c
}

func (f *freeCache) Set(c *config.ConfigResponse) error {
	value, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return f.cache.Set(getKey(c.AppName, c.ClusterName, c.NamespaceName), value, -1)
}

func (f *freeCache) Get(c *config.QueryConfigRequest) (*config.ConfigResponse, bool) {
	v, err := f.cache.Get(getKey(c.AppName, c.ClusterName, c.NamespaceName))
	if err != nil {
		return nil, false
	}

	var value config.ConfigResponse
	err = json.Unmarshal(v, &value)
	if err != nil {
		log.Error("json unmarshal err")
		return nil, false
	}
	return &value, true
}
func (f *freeCache) Clear() {
	f.cache.Clear()
}

func getKey(appName, clusterName, namespaceName string) []byte {
	return convert.StringToBytes(fmt.Sprintf("%s/%s/%s", appName, clusterName, namespaceName))
}
