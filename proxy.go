package main

import (
	"./resolver"
	"./rproxy"
	"flag"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type Options struct{
	Store string
	ListenPort int
	ApiPort int
}

func main() {
	options := new(Options)
	flag.StringVar(&options.Store, "store", "redis://:6379", "store")
	flag.IntVar(&options.ListenPort, "listen", 8080, "listen port")
	flag.IntVar(&options.ApiPort, "api_listen", 8090, "api server listen port")
	flag.Parse()

	var resolve *resolver.Resolver

	store_url, err := url.Parse(options.Store)
	if err != nil {
		log.Fatal("Could not parse store url")
	}
	switch store_url.Scheme {
	case "redis":
		log.Print("using redis as resolver storage \"", store_url.Host, "\"")
		resolve = resolver.NewRedisResolver(store_url.Host)
	case "memory":
		log.Print("using memory as resolver storage")
		log.Print("Be careful. all data has gone if you stop this process.")
		resolve = resolver.NewMemoryResolver()
	case "file":
		file_path := store_url.Host + store_url.Path
		log.Print("using file as resolver storage \"", file_path, "\"")
		resolve = resolver.NewFileResolver(file_path)
	}
	reverse_proxy := rproxy.NewReverseProxy(resolve)

	api_server := http.NewServeMux()
	api_server.HandleFunc("/", reverse_proxy.APIServer())
	api_server_listen := ":" + strconv.Itoa(options.ApiPort)
	go http.ListenAndServe(api_server_listen, api_server)
	log.Print("API Server listen ", api_server_listen)

	http.Handle("/", reverse_proxy)
	listen := ":" + strconv.Itoa(options.ListenPort)
	log.Print("Reverse Proxy Listen ", listen)
	http.ListenAndServe(listen, nil)
}
