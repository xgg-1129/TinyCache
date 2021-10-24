package TinyCache

import "testing"

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
	cache := NewCache(1024)
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
	lru := NewCache(int64(Cap))
	lru.Add(k1, str(v1))
	lru.Add(k2, str(v2))
	lru.Add(k3, str(v3))

	if _, ok := lru.Get("key1"); ok || lru.Capacity() != 2 {
		t.Fatalf("Removeoldest key1 failed")
	}
}