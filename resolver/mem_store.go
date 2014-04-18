package resolver

type MemStorage struct {
	Table map[string]string
}

func NewMemStorage() *MemStorage {
	return &MemStorage{Table: make(map[string]string)}
}

func (this *MemStorage) Get(key string) string {
	return this.Table[key]
}

func (this *MemStorage) Set(key string, value string) {
	this.Table[key] = value
}

func (this *MemStorage) List() []string {
	var keys []string

	for k := range this.Table {
		keys = append(keys, k)
	}

	return keys
}
