package TinyCache

import (
	"sort"
	"strconv"
)
//一致性哈希算法的关键就是，计算的hashcode不是它的位置，而是一个离位置最近的值

type Hash func(data []byte) uint32
type PeerMap struct {
	hashFun Hash
	//虚拟结点的倍数
	replicas int
	//分布式系统上所有的虚拟节点
	keys []int
	hashMap map[int]string

}
//添加结点,这里传递的是每个结点的ip:port地址
func (p *PeerMap) Add(addrs... string) {
	for _, key := range addrs{
		for i:=0;i<p.replicas;i++{
			hashCode:=int(p.hashFun(([]byte)(key+strconv.Itoa(i))))
			p.keys=append(p.keys,hashCode)
			p.hashMap[hashCode]=key
		}
	}
}
//根据缓存的key，定位到某一台具体的主机
func (p *PeerMap) GetOtherPeer(key string)string{
	//首先找到虚拟节点
	if len(p.keys)==0{
		return ""
	}
	code:=int(p.hashFun([]byte(key)))
	//顺时针查找第一个虚拟节点
	idix:=sort.Search(len(p.keys), func(i int) bool {
		return code-p.keys[i]<=0
	})
	return p.hashMap[p.keys[idix]]
}

