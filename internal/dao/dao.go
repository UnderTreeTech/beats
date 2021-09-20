package dao

import (
	"beats/internal/dao/iface"
	"context"
	"errors"
	"fmt"

	"github.com/UnderTreeTech/waterdrop/pkg/database/es"

	"github.com/UnderTreeTech/waterdrop/pkg/conf"
	"github.com/UnderTreeTech/waterdrop/pkg/database/redis"
	"github.com/UnderTreeTech/waterdrop/pkg/log"
)

// interface Dao
type Dao interface {
	Close() error
	Ping(ctx context.Context) error
	Redis() *redis.Redis

	iface.Es
}

// struct dao
type dao struct {
	redis *redis.Redis
	es    *es.Client
}

// New return a dao that implements interface Dao
func New() Dao {
	redis := NewRedis()
	es := NewEs()
	return &dao{
		redis: redis,
		es:    es,
	}
}

// Close close backend base services
func (d *dao) Close() (err error) {
	d.redis.Close()
	return
}

// Ping ping backend base services, like db, redis, es etc.
func (d *dao) Ping(ctx context.Context) error {
	if alive := d.redis.Ping(ctx); !alive {
		return errors.New("redis has gone")
	}

	return nil
}

func (d *dao) Redis() *redis.Redis {
	return d.redis
}

// NewRedis returns redis instance
func NewRedis() *redis.Redis {
	config := &redis.Config{}
	if err := conf.Unmarshal("redis", config); err != nil {
		panic(fmt.Sprintf("unmarshal redis config fail,err msg %s", err.Error()))
	}
	log.Infof("redis config", log.Any("config", config))

	redis, err := redis.New(config)
	if err != nil {
		panic(fmt.Sprintf("new redis client fail,err msg %s", err.Error()))
	}

	return redis
}

// NewEs returns an instance of es client
func NewEs() *es.Client {
	config := &es.Config{}
	if err := conf.Unmarshal("es", config); err != nil {
		panic(fmt.Sprintf("unmarshal es config fail,err msg %s", err.Error()))
	}
	log.Infof("es config", log.Any("config", config))

	cli := es.NewClient(config)
	return cli
}
