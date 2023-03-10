package cache

import "github.com/VictoriaMetrics/fastcache"

var cache = fastcache.New(1) // 这里会申请一个32MB的cache空间

// 保存数据到缓存
func Set(key string, value []byte) {
	cache.Set([]byte(key), value)
}

// 获取值，不存在则返回nil
func Get(key string) []byte {
	if v, exist := cache.HasGet(nil, []byte(key)); exist && len(v) != 0 {
		return v
	} else {
		return nil
	}
}

// 获取或者创建
func GetSet(key string, value []byte) []byte {
	var res []byte
	if res = Get(key); res != nil {
		return res
	} else {
		Set(key, value)
		return value
	}
}

// 删除键值对
func Del(key string) {
	cache.Del([]byte(key))
}

// 清空缓存
func Reset() {
	cache.Reset()
}
