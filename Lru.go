package TinyCache

import (
	"container/list"
	"log"
	"os"
	"sync"
)

/*
	使用LRU算法的缓存
*/

type LRUCache struct {

	//cache的最大字节容量
	maxByte int64
	//当前已用字节数
	usedByte int64

	twoList *list.List
	cache map[string]*list.Element
}
//*Entry是双向链表结点的data值
type Entry struct {
	key string
	data Value
}
type Value interface {
	Length() int64
}
//Data表示存储的数据
type Data struct {
	b []byte
}
func (d Data) Length()int64{
	return int64(len(d.b))
}
func (d Data) toString() string{
	return string(d.b)
}

func NewLRUCache(maxByte int64)*LRUCache {
	cache:=new(LRUCache)
	cache.maxByte=maxByte
	cache.twoList=list.New()
	cache.cache=make(map[string]*list.Element)
	return cache
}
//返回链表Elements的个数
func (c *LRUCache) Capacity() int {
	return c.twoList.Len()
}

func (c *LRUCache) Get(key string)(Value,bool){
	if element,ok := c.cache[key];ok{
		c.twoList.MoveToFront(element)
		entry := element.Value.(*Entry)
		return entry.data,true
	}
	return nil,false
}
func (c *LRUCache) RemoveOldest() {
	if c.twoList.Len()== 0{
		return
	}else {
		ele:=c.twoList.Back()
		entry:=ele.Value.(*Entry)
		LRUInfof("Remove element which key is %s ",entry.key,entry.data)
		delete(c.cache,entry.key)
		c.twoList.Remove(ele)
		c.usedByte=c.usedByte-entry.data.Length()
		c.usedByte=c.usedByte-(int64)(len(entry.key))
	}
}
func (c *LRUCache) Add(key string,data Value){
	if element,ok:=c.cache[key];ok{
		//如果cache本身存在
		entry:=element.Value.(*Entry)
		LRUInfof("the element which key is %s  existed",key)
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
		LRUInfof("add a element ,key is %s",key)
		//修改长度
		c.usedByte+=data.Length()+(int64)(len(key))
	}
	for c.usedByte > c.maxByte{
		c.RemoveOldest()
	}
}

/*为LRUcache添加个日志*/

var (
	errorLog = log.New(os.Stdout, "\033[31m[error]\033[0m ", log.LstdFlags|log.Lshortfile)
	infoLog  = log.New(os.Stdout, "\033[34m[info]\033[0m ", log.LstdFlags|log.Lshortfile)
	LRUinfoLog  = log.New(os.Stdout, "\033[34m[LRU_info]\033[0m ", log.LstdFlags|log.Lshortfile)
	loggers  = []*log.Logger{errorLog, infoLog}
	mu       sync.Mutex
)
var (
	Error  = errorLog.Println
	Errorf = errorLog.Printf
	Info   = infoLog.Println
	Infof  = infoLog.Printf
	LRUInfof =LRUinfoLog.Printf
)
