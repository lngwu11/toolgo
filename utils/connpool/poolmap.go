package connpool

import (
	"sync"
)

// 适用于读多写少场景
var poolMap sync.Map

// GetPool get a pool with the specified key from the global map.
// newPool will be executed and stored if not found.
func GetPool(key string, config *Config) Pool {
	p, ok := poolMap.Load(key)
	if !ok {
		p = newPool(key, config)
		poolMap.Store(key, p)
	}
	return p.(Pool)
}
