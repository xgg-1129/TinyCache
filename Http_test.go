package TinyCache

import (
	"fmt"
	"log"
	"testing"
)


func TestHTTP(t *testing.T) {
	NewGroup("scores",2<<10, GetterFun(func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	server := NewServer()
	server.Run(":9999")
}
