package main

import (
	"./rproxy"
	"net/http"
)

func main() {
	reverse_proxy := rproxy.NewReverseProxy()

	api_server := http.NewServeMux()
	api_server.HandleFunc("/", reverse_proxy.APIServer())
	go http.ListenAndServe(":8090", api_server)

	http.Handle("/", reverse_proxy)
	http.ListenAndServe(":8080", nil)
}
