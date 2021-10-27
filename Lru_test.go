package TinyCache

import (
	"fmt"
	"log"
	"testing"
)

type str string

func (s str) Length() int64 {
	return int64(len(s))
}


func TestLRU(t *testing.T) {
	t.Run("TestGet", func(t *testing.T) {
		Get(t)
	})
	t.Run("TestGet", func(t *testing.T) {
		Remove(t)
	})
}
func Get(t *testing.T){
	cache := NewLRUCache(1024)
	cache.Add("k1",str("aaa"))
	cache.Add("k2",str("bbb"))
	cache.Add("k2",str("ccc"))
	if cache.Capacity()!=2{
		t.Error("the capacity of cache is not 2")
	}
	value, ok:= cache.Get("k2")
	if !ok{
		t.Error("can not find k2")
	}
	if 	s := value.(str);s!="ccc"{
		t.Error("can not update k2")
	}
}
func Remove(t *testing.T){
	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"
	Cap := len(k1 + k2 + v1 + v2)
	lru := NewLRUCache(int64(Cap))
	lru.Add(k1, str(v1))
	lru.Add(k2, str(v2))
	lru.Add(k3, str(v3))

	if _, ok := lru.Get("key1"); ok || lru.Capacity() != 2 {
		t.Fatalf("Removeoldest key1 failed")
	}
}
var 	db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}
func TestGroup(t *testing.T) {
	loadCounts := make(map[string]int, len(db))
	gee := NewGroup("scores", 2<<10, GetterFun(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				if _, ok := loadCounts[key]; !ok {
					loadCounts[key] = 0
				}
				loadCounts[key] += 1
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	for k, v := range db {
		if view, err := gee.Get(k); err != nil || view.toString() != v {
			t.Fatal("failed to get value of Tom")
		} // load from callback function
		if _, err := gee.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		} // cache hit
	}

	if view, err := gee.Get("unknown"); err == nil {
		t.Fatalf("the value of unknow should be empty, but %s got", view)
	}
}