package TinyCache

import "log"

//Group负责暴露提服务的接口
var Groups =make(map[string]*Group)

type Group struct {
	name string
	cache *Cache
	get  Getter
	peer PeerPicker
}
type Getter interface {
	Get(key string)(value []byte,err error)
}
type GetterFun func(key string)(value []byte,err error)

func (f GetterFun) Get(key string)(value []byte,err error) {
	return f(key)
}
func NewGroup(name string, maxBytes int64, getter Getter) *Group {
	g:=new(Group)
	g.cache=NewCache(maxBytes)
	g.name=name
	g.get=getter
	Groups[name]=g
	return g
}
func (g *Group) Get(key string) (Data,error){
	if value, ok := g.cache.get(key);ok{
		return value,nil
	}
	return g.Load(key)

}
func (g *Group) LoadLocal(key string)(Data,error){
	v, err := g.get.Get(key)
	if err!=nil{
		return Data{},err
	}
	value:=Data{b: v}
	g.UpdateCache(key,value)
	return value,nil
}
func (g *Group) Load(key string)(Data,error){
	if peer, ok := g.peer.PickPeer(key); ok {
		value, err := g.getFromPeer(peer, key)
		if  err == nil {
			return value, nil
		}
		log.Println("[GeeCache] Failed to get from peer", err)
	}
	return g.LoadLocal(key)
}
func (g *Group) getFromPeer(peer PeerGetter, key string) (Data, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return Data{}, err
	}
	return Data{b: bytes}, nil
}
func (g *Group) UpdateCache(key string,value Data){
	g.cache.add(key,value)
}