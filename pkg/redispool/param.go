package redispool

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/garyburd/redigo/redis"
)

var paramBuilder ParamBuilder

// 根据name 构建 PoolParam
type ParamBuilder func(name string) *PoolParam

func RegisterParamBuilder(builder ParamBuilder) {
	paramBuilder = builder
}

// PoolParam 连接池参数
type PoolParam struct {
	Name                   string         //连接池名字
	Addr                   string         //redis地址 ip:port  @deprecated
	addrList               []string       //负载均衡的地址
	lock                   *sync.RWMutex  //用于加锁addr列表
	state                  []bool         //各台机器的状态
	avls                   int            //当前可用的地址个数
	rrIndex                uint32         //RoundRobin的索引下标
	Policy                 RoutePolicy    //负载均衡策略
	Pass                   string         //密码，没有则为空
	MaxIdle                int            //做多允许多少空闲连接
	MaxActive              int            //最多允许多少并发连接
	IdleTimeoutSecond      int            //单个链接空闲多久被销毁,单位秒
	NetworkTimeoutMsec     int            //网络超时(connect,send,recv)，单位毫秒
	HealthCheckIntervalSec int            //定时检查池子里面的节点是是否健康
	retries                map[string]int //每个地址的失败次数
}

// 以下一些函数为了在运行时动态检测redis的地址列表，并且踢掉不健康的，以及运行时动态设置地址列表
func (p *PoolParam) setup() (err error) {
	if p.Addr == "" {
		err = errors.New("redis addr is empty")
		return
	}
	p.retries = make(map[string]int)

	p.lock = &sync.RWMutex{}

	p.addrList = strings.Split(p.Addr, ",")

	sort.Strings(p.addrList)
	p.state = make([]bool, len(p.addrList))
	for i, addr := range p.addrList {
		if err = p.checkRedisAddr(addr); err != nil {
			return
		}
		p.state[i] = true
	}
	p.avls = len(p.state)
	return
}

func (p *PoolParam) doCheck() {
	p.lock.RLock()
	defer p.lock.RUnlock()
	t := 0
	if len(p.addrList) <= 1 {
		return
	}
	for i, addr := range p.addrList {
		if err := p.checkRedisAddr(addr); err != nil {
			log.Printf("[ERROR] check redis addr:%s fail:%s\n", addr, err.Error())
			retries := p.retries[addr]
			retries++
			p.retries[addr] = retries
			if retries >= 3 {
				p.state[i] = false
				log.Printf("[ERROR] redis addr:%s monitor check failed already reached:%d times,set addr status to DOWN\n", addr, retries)
			}
			continue
		}
		if p.state[i] == false {
			log.Printf("[INFO] redis addr:%s monitor check succ\n", addr)
		}
		p.state[i] = true
		p.retries[addr] = 0
		t++
	}
	//如果经过检测，所有的地址都不可用，那么强制将第一个地址作为本次检测的可用地址
	if t == 0 {
		t = 1
		p.state[0] = true
	}
	p.avls = t
}

func (p *PoolParam) monitorCheck() {
	sec := p.HealthCheckIntervalSec
	if sec <= 0 {
		sec = 1
	}
	tk := time.NewTicker(time.Duration(sec) * time.Second)
	for {
		select {
		case <-tk.C:
			p.doCheck()
		}
	}
}

func (p *PoolParam) dumpStatus() string {
	state := p.state
	addrList := p.addrList
	strArr := make([]string, len(addrList))
	for i := 0; i < len(strArr); i++ {
		strArr[i] = fmt.Sprintf("addr:%s,status:%v", addrList[i], state[i])
	}
	return strings.Join(strArr, "\n")
}

func (p *PoolParam) setAddr(addr string) (err error) {
	if addr == "" {
		err = errors.New("redis addr is empty")
		return
	}
	addrList := strings.Split(addr, ",")
	sort.Strings(addrList)
	state := make([]bool, len(addrList))
	t := 0

	p.lock.Lock()
	defer p.lock.Unlock()

	for i, addr := range addrList {
		if e := p.checkRedisAddr(addr); e != nil {
			state[i] = false
			continue
		}
		state[i] = true
		t++
	}
	if t <= 0 {
		err = errors.New("no reachable servers")
		return
	}
	p.Addr = addr

	p.avls, p.addrList, p.state = t, addrList, state
	return

}

func (p *PoolParam) checkRedisAddr(addr string) (err error) {
	to := 1 * time.Second
	ops := []redis.DialOption{redis.DialConnectTimeout(to), redis.DialReadTimeout(to), redis.DialWriteTimeout(to)}
	conn, err := redis.Dial("tcp", addr, ops...)
	if err != nil {
		log.Printf("redis dial error:%s\n", err.Error())
		return
	}
	defer conn.Close()
	if p.Pass != "" {
		if _, err = redis.String(conn.Do("AUTH", p.Pass)); err != nil {
			log.Printf("redis auth error:%s\n", err.Error())
			return
		}
	}
	if _, err = conn.Do("PING"); err != nil {
		log.Printf("redis ping error:%s", err.Error())
	}

	return
}

func (p *PoolParam) getAddr() string {
	if p.avls <= 0 {
		return ""
	}
	maxLoop := 10
	p.lock.RLock()
	defer p.lock.RUnlock()
	nAddrs := len(p.addrList)
	switch p.Policy {
	case RoutePolicyRandom:
		for maxLoop > 0 {
			rnd := rand.Intn(nAddrs)
			if p.state[rnd] {
				return p.addrList[rnd]
			}
			maxLoop--
		}
	case RoutePolicyRoundRobin:
		for maxLoop > 0 {
			t := atomic.AddUint32(&p.rrIndex, 1) % uint32(nAddrs)
			if p.state[t] {
				return p.addrList[t]
			}
			maxLoop--
		}
	}
	return ""
}
