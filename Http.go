package TinyCache

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

const (
	defaultPath = "/cache/"
)

type HTTPPool struct {
	ip string
	path string

	mu sync.Mutex
	peers *PeerMap
	//这个map中存的是其他节点的ip地址
	httpGetter map[string]*Client
}
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}
// 根据key，获取其他结点的ip地址和方法
func (p *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	//如果远程结点存在，并且不等于自己
	if peer := p.peers.GetOtherPeer(key); peer != "" && peer != p.ip{
		return p.httpGetter[peer], true
	}
	return nil, false
}

func NewServer()*HTTPPool {
	return &HTTPPool{
		path: defaultPath,
	}
}

func (s *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	if !strings.HasPrefix(r.URL.Path,defaultPath){
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	parts := strings.Split(r.URL.Path[len(s.path):], "/")
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	groupName:=parts[0]
	key:=parts[1]

	group:=Groups[groupName]
	if group==nil{
		http.Error(w, "404 Not Found the group:"+groupName, http.StatusNotFound)
		return
	}
	data,err:=group.Get(key);
	if err!=nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	w.Write(data.b)
}

func (s *HTTPPool) Run(addr string) {
	http.ListenAndServe(addr,s)
}

type Client struct {
	addrPath string

}

func (c *Client) Get(groupName string,key string) ([]byte,error) {
	u := fmt.Sprintf("%v%v%v/%v", c.addrPath,  groupName, key)
	get, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer get.Body.Close()
	if get.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", get.Status)
	}
	res, err := io.ReadAll(get.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}
	return res, nil
}
