package main

import {
	"net/http/httputils"
	"net/url"
	"net/http"
	"fmt"
	"os"

}

type servers interface{
	Address() string
	IsAlive() bool
	Serve(rw http.ResponseWriter,r *http.Request)

}

type simpleServer struct{
	addr  string
	proxy httputil.ReverseProxy
}

type LoadBalancer struct{
	port 			 string
	roundRobinCount  int
	server           []servers
}
func newSimpleServer (addr string) *simpleServer{
	serverUrl, err := url.Parse()
	handleErr(err)

	return simpleServer{
		addr :addr,
		proxy:httputil.NewSingleHostReverseProxy(serverUrl)

	}
}

func NewLoadBalancer(port string,server []servers) *LoadBalancer{
	return LoadBalancer{
		port:port,
		roundRobinCount:0
		server:server
	}


} 

func handleErr(err error){
	if err!=nil{
		fmt.Printf("Error:%v\n",err)
		os.exit(1)
	
	}

}

func (s *simpleServer) Address() string {
	return s.addr
}

func (s *simpleServer) IsAlive() bool {
	return true
}

func (s *simpleServer) Serve(rw http.ResponseWriter,r *http.Request) {
	s.proxy.ServeHttp(rw,r)
} 

func (*LoadBalancer) getNextAvailableServer() servers{
	server:=lb.server[lb.roundRobinCount%len(server)]
	for server.IsAlive(){
		lb.roundRobinCount++
		server:=lb.server[lb.roundRobinCount%len(server)]
		
	}
	lb.roundRobinCount++
	return server
}

func(*LoadBalancer) serverProxy(rw http.ResponseWriter,r *http.Request){
	targetServer := lb.getNextAvailableServerr()
	fmt.Printf("forwaring request to %q\n",targerServer.Address())

	targerServer.Serve(rw,r)


}

func main(){

	server := []servers{
		newSimpleServer("http://www.google.com")
		newSimpleServer("http://www.youtube.com")
		newSimpleServer("http://www.bing.com")
	}

	lb := NewLoadBalancer("8000",server)
	handleRedirect := func (rw http.ResponseWriter,r *http.Request){
		lb.serverProxy(rw,r)
	}
	http.HandleFunc("/",handleRedirect)

	fmt.Printf("server serving at host %s\n",lb.port)
	httpListenAndServe(":"lb.port,nil)
}