package rproxy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	//"log"
	"../resolver"
)

type RProxy struct {
	ReverseProxy *httputil.ReverseProxy
	Resolver     *resolver.Resolver
}

func NewReverseProxy() *RProxy {
	host_resolver := resolver.NewRedisResolver(":6379")

	director := func(request *http.Request) {
		request.URL.Scheme = "http"
		request.URL.Host = host_resolver.Resolve(request.Host)
	}
	reverse_proxy := &httputil.ReverseProxy{Director: director}

	return &RProxy{ReverseProxy: reverse_proxy, Resolver: host_resolver}
}

func (this *RProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	this.ReverseProxy.ServeHTTP(rw, req)
}

func (this *RProxy) APIServer() func(http.ResponseWriter, *http.Request) {
	server := func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" && req.URL.Path == "/delete" {
			_ = req.ParseForm()

			this.Resolver.Delete(req.PostForm.Get("from"))
			fmt.Fprint(w, "ok")
			return
		}

		if req.Method == "GET" && req.URL.Path == "/" {
			encoder := json.NewEncoder(w)
			encoder.Encode(this.Resolver.List())
			return
		}

		if req.Method == "POST" && req.URL.Path == "/" {
			_ = req.ParseForm()

			this.Resolver.Set(req.PostForm.Get("from"), req.PostForm.Get("to"))
			fmt.Fprint(w, "ok")
			return
		}

		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	return server
}
