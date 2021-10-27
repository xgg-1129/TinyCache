package TinyCache

import "sync"

//call代表一次函数调用
type call struct {
	sync.WaitGroup
	res interface{}
	err error
}

type CallLock struct {
	mu sync.Mutex
	callMap map[string]*call
}

func (c *CallLock) Call(name string,f func()(interface{},error))(interface{},error){
	mu.Lock()
	if c.callMap==nil{
		c.callMap=make(map[string]*call)
	}
	if c,ok:=c.callMap[name];ok{
		mu.Unlock()
		c.Wait()
		return c.res,c.err
	}
	//如果没有正在执行的call
	newCall:=new(call)
	newCall.Add(1)
	c.callMap[name]=newCall
	mu.Unlock()

	newCall.res,newCall.err=f()
	newCall.Done()

	c.mu.Lock()
	delete(c.callMap,name)
	c.mu.Unlock()
	return newCall.res,newCall.err
}
