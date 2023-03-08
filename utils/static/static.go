package static

import "sync"

var staticMap sync.Map

func Set(k string, v string) {
	staticMap.Store(k, v)
}

func Get(k string) string {
	v, ok := staticMap.Load(k)
	if ok {
		return v.(string)
	}
	return ""
}

func Del(k string) {
	staticMap.Delete(k)
}

func GetSet(k string, v string) string {
	value, _ := staticMap.LoadOrStore(k, v)

	return value.(string)
}
