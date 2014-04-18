package resolver

import (
    "github.com/garyburd/redigo/redis"
    "time"
    "log"
    "os"
    "encoding/json"
)

type ResolverStorage interface {
    Get(string) string
    Set(string, string)
    Delete(string)
    List() []string
}

type RedisStorage struct {
    Pool *redis.Pool
}

func NewRedisStorage(host string) *RedisStorage {
    pool := &redis.Pool{
        MaxIdle: 3,
        IdleTimeout: 240 * time.Second,
        Dial: func () (redis.Conn, error) {
            c, err := redis.Dial("tcp", host)
            if err != nil {
                return nil, err
            }
            return c, err
        },
    }

    return &RedisStorage{Pool: pool}
}

func (this *RedisStorage) Get(key string) (string) {
    conn := this.Pool.Get()
    defer conn.Close()

    result, err := redis.String(conn.Do("GET", key))
    if err != nil {
        log.Printf("%s could not get", key)
        return ""
    }

    return result
}

func (this *RedisStorage) Set(key string, value string) {
    conn := this.Pool.Get()
    defer conn.Close()

    _, err := redis.String(conn.Do("SET", key, value))
    if err != nil {
        log.Printf("%s could not set", key)
    }

    _, err = redis.String(conn.Do("RPUSH", "host-list", key))
    if err != nil {
        log.Printf("%s could not set", key)
    }
}

func (this *RedisStorage) Delete(key string) {
    conn := this.Pool.Get()
    defer conn.Close()

    _, err := conn.Do("DEL", key)
    if err != nil {
        log.Printf(err.Error())
        log.Printf("Could not delete key %s", key)
    }

    _, err = conn.Do("LREM", "host-list", 1, key)
    if err != nil {
        log.Printf(err.Error())
    }
}

func (this *RedisStorage) List() ([]string) {
    conn := this.Pool.Get()
    defer conn.Close()

    result, err := redis.Strings(conn.Do("LRANGE", "host-list", 0, -1))
    if err != nil {
        log.Printf("Could not get list")
        return make([]string, 0)
    }

    return result
}

/*
 *LocalStorage
 */

type LocalStorage struct {
    FilePath string
    Table map[string]string
}

func NewLocalStorage(filePath string) *LocalStorage {
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

    return &LocalStorage{FilePath: filePath, Table: table}
}

func (this *LocalStorage) Get(key string) (string) {
    return this.Table[key]
}

func (this *LocalStorage) Set(key string, value string) {
    this.Table[key] = value

    this.save()
}

func (this *LocalStorage) Delete(key string) {
    delete(this.Table, key)

    this.save()
}

func (this *LocalStorage) List() ([]string) {
    var keys []string

    for k := range this.Table {
        keys = append(keys, k)
    }

    return keys
}

func (this *LocalStorage) save() {
    file, err := os.OpenFile(this.FilePath, os.O_WRONLY|os.O_CREATE, 0644)
    defer file.Close()
    if err != nil {
        log.Fatal("Could not open file ", this.FilePath)
        return
    }

    b, err := json.Marshal(this.Table)
    _, err = file.Write(b)
    if err != nil {
        log.Fatal("Could not write storage file")
        return
    }
}
