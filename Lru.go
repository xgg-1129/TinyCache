package TinyCache

import "container/list"

/*
	使用LRU算法的缓存
*/

type Cache struct {

	//cache的最大字节容量
	maxByte int64
	//当前已用字节数
	usedByte int64

	twoList *list.List
	cache map[string]*list.Element

}
type Entry struct {
	key string
	data Value
}
type Value interface {
	Length() int64
}
func NewCache(maxByte int64)*Cache{
	cache:=new(Cache)
	cache.maxByte=maxByte
	cache.twoList=list.New()
	cache.cache=make(map[string]*list.Element)
	return cache
}
//返回链表Elements的个数
func (c *Cache) Capacity() int {
	return c.twoList.Len()
}

func (c *Cache) Get(key string)(Value,bool){
	if element,ok := c.cache[key];ok{
		c.twoList.MoveToFront(element)
		entry := element.Value.(*Entry)
		return entry.data,true
	}
	return nil,false
}
func (c *Cache) RemoveOldest() {
	if c.twoList.Len()== 0{
		return
	}else {
		ele:=c.twoList.Back()
		entry:=ele.Value.(*Entry)
		delete(c.cache,entry.key)
		c.twoList.Remove(ele)
		c.usedByte=c.usedByte-entry.data.Length()
		c.usedByte=c.usedByte-(int64)(len(entry.key))
	}
}
func (c *Cache) Add(key string,data Value){
	if element,ok:=c.cache[key];ok{
		//如果cache本身存在
		entry:=element.Value.(*Entry)
		//修改长度
		c.usedByte=c.usedByte-entry.data.Length()+data.Length()
		//改变在双向链表的位置
		c.twoList.MoveToFront(element)
		//改变值
		entry.data=data
	}else{
		//生成新结点
		newElement:=c.twoList.PushFront(&Entry{key:  key, data: data})
		//在map中注册
		c.cache[key]=newElement
		//修改长度
		c.usedByte+=data.Length()+(int64)(len(key))
	}
	for c.usedByte > c.maxByte{
		c.RemoveOldest()
	}
}
