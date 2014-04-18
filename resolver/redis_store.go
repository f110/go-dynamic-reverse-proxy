package resolver

import (
    "github.com/garyburd/redigo/redis"
    "log"
    "time"
)

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
