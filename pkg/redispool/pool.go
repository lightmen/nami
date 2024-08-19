package redispool

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
)

var lock sync.Mutex
var poolMap sync.Map

// Pool implements the connection pool of redis
type Pool struct {
	pool  *redis.Pool
	param *PoolParam
}

// RoutePolicy redis各个地址负载均衡策略
type RoutePolicy int

const (
	// RoutePolicyRoundRobin 轮转
	RoutePolicyRoundRobin = iota
	// RoutePolicyRandom 随机
	RoutePolicyRandom
)

// DumpStatus redis连接池对外暴露的状态接口
func (p *Pool) DumpStatus() string {
	return p.param.dumpStatus()
}

// SetAddrList 支持运行时动态设置redis服务器地址
func (p *Pool) SetAddrList(addrList string) {
	p.param.setAddr(addrList)
}

// GetClient alloc a connection from pool
func (p *Pool) GetClient() (cli *Client, err error) {
	conn := p.pool.Get()
	if conn.Err() != nil {
		err = conn.Err()
		return
	}

	cli = &Client{Conn: conn}
	return
}

func newPool(pparam *PoolParam) (pool *Pool, err error) {
	param := new(PoolParam)
	*param = *pparam
	opts := make([]redis.DialOption, 0, 3)
	if param.NetworkTimeoutMsec != 0 {
		to := time.Duration(param.NetworkTimeoutMsec) * time.Millisecond
		opts = append(opts, redis.DialConnectTimeout(to), redis.DialReadTimeout(to), redis.DialWriteTimeout(to))
	}

	if err = param.setup(); err != nil {
		return
	}
	p := &redis.Pool{
		MaxIdle:     param.MaxIdle,
		Wait:        true,
		MaxActive:   param.MaxActive,
		IdleTimeout: time.Duration(param.IdleTimeoutSecond) * time.Second,

		Dial: func() (c redis.Conn, err error) {
			addr := param.getAddr()

			if addr == "" {
				err = errors.New("no server available")
				return
			}
			c, err = redis.Dial("tcp", addr, opts...)
			if err != nil {
				err = errors.New(err.Error())
				return
			}
			if param.Pass != "" {
				if _, err = c.Do("AUTH", param.Pass); err != nil {
					err = errors.New(err.Error())
					c.Close()
					c = nil
				}
			}
			return
		},
		// TestOnBorrow: func(c redis.Conn, t time.Time) error {
		// 	_, err := c.Do("PING")
		// 	return err
		// },
	}
	//做链接检查
	conn := p.Get()
	defer conn.Close()
	if _, err = conn.Do("PING"); err != nil {
		err = errors.New(err.Error())
		return
	}

	pool = &Pool{pool: p, param: param}
	go param.monitorCheck()
	return
}

// NewPool 新建Redis连接池,默认网络超时时间为0
func NewPool(addr, pass string, maxIdle, maxActive, idleTimeoutSecond int) (pool *Pool, err error) {
	param := &PoolParam{
		Addr:              addr,
		Pass:              pass,
		MaxIdle:           maxIdle,
		MaxActive:         maxActive,
		IdleTimeoutSecond: idleTimeoutSecond,
	}
	pool, err = newPool(param)
	return
}

// NewPoolWithParam 使用param参数创建连接池
func NewPoolWithParam(param *PoolParam) (pool *Pool, err error) {
	pool, err = newPool(param)
	if err != nil {
		return
	}

	return
}

// GetPool 根据name，获取该name对应的pool
func GetPool(name string) (pool *Pool, err error) {
	val, ok := poolMap.Load(name)
	if ok {
		return val.(*Pool), nil
	}

	if paramBuilder == nil {
		err = fmt.Errorf("can not found redis pool for %s", name)
		return
	}

	param := paramBuilder(name)

	lock.Lock()
	defer lock.Unlock()

	val, ok = poolMap.Load(name)
	if ok {
		return val.(*Pool), nil
	}

	pool, err = NewPoolWithParam(param)
	if err != nil {
		return
	}

	poolMap.Store(param.Name, pool)
	return
}
