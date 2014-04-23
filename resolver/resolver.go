package resolver

import (
	"log"
)

type Resolver struct {
	Storage ResolverStorage
}

func (r *Resolver) Resolve(host string) string {
	result := r.Storage.Get(host)

	log.Println("Resolved: from", host, "to", result)
	return result
}

func (r *Resolver) Set(from string, to string) {
	r.Storage.Set(from, to)
}

func (r *Resolver) Delete(from string) {
	r.Storage.Delete(from)
}

func (r *Resolver) List() []string {
	return r.Storage.List()
}

func NewRedisResolver(redisHost string) *Resolver {
	resolverStorage := NewRedisStorage(redisHost)

	resolver := &Resolver{Storage: resolverStorage}

	return resolver
}

func NewFileResolver(filePath string) *Resolver {
	localStorage := NewFileStorage(filePath)

	return &Resolver{Storage: localStorage}
}

func NewMemoryResolver() *Resolver {
	localStorage := NewMemStorage()

	return &Resolver{Storage: localStorage}
}
