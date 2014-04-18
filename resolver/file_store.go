package resolver

import (
	"encoding/json"
	"log"
	"os"
)

type FileStorage struct {
	FilePath string
	Table    map[string]string
}

func NewFileStorage(filePath string) *FileStorage {
	table := make(map[string]string)

	file, err := os.Open(filePath)

	if err == nil {
		decoder := json.NewDecoder(file)

		for {
			if err := decoder.Decode(&table); err != nil {
				break
			}
		}
	}

	return &FileStorage{FilePath: filePath, Table: table}
}

func (this *FileStorage) Get(key string) string {
	return this.Table[key]
}

func (this *FileStorage) Set(key string, value string) {
	this.Table[key] = value

	this.save()
}

func (this *FileStorage) Delete(key string) {
	delete(this.Table, key)

	this.save()
}

func (this *FileStorage) List() []string {
	var keys []string

	for k := range this.Table {
		keys = append(keys, k)
	}

	return keys
}

func (this *FileStorage) save() {
	file, err := os.OpenFile(this.FilePath, os.O_WRONLY|os.O_CREATE, 0644)
	defer file.Close()
	if err != nil {
		log.Fatal("Could not open file ", this.FilePath)
		return
	}

	_, err = file.Write(this.serialize())
	if err != nil {
		log.Fatal("Could not write storage file")
		return
	}
}

func (this *FileStorage) serialize() []byte {
	b, err := json.Marshal(this.Table)
	if err != nil {
		log.Fatal("Could not serialize")
		return nil
	}

	return b
}
