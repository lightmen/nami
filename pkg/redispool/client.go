package redispool

import (
	"crypto/md5"
	"encoding/hex"
	"errors"

	"github.com/garyburd/redigo/redis"
)

// Client Wrapper the Redis conn
type Client struct {
	redis.Conn
}

// GetClient gets a redis client by business key
func GetClient(name string) (cli *Client, err error) {
	pool, err := GetPool(name)
	if err != nil {
		return
	}

	cli, err = pool.GetClient()
	if err != nil {
		return
	}

	return
}

// kv start
type KeyType string

const (
	KeyTypeNone   KeyType = "none"
	KeyTypeString KeyType = "string"
	KeyTypeList   KeyType = "list"
	KeyTypeSet    KeyType = "set"
	KeyTypeZSet   KeyType = "zset"
	KeyTypeHash   KeyType = "hash"
)

// 获取指定类型, 返回值err表示redis错误, ok 表示key是否存在, t为具体的类型
func (c *Client) Type(key string) (t KeyType, ok bool, err error) {
	var ret string

	ok = true
	ret, err = redis.String(c.doCmd("TYPE", key))
	t = KeyType(ret)

	if err == nil { //如有错误, 判断是否type none
		if t == KeyTypeNone {
			ok = false
		}

		return
	}

	if err == redis.ErrNil { // 如果错误为ErrNil, 正常来说不会出现
		ok = false
		err = nil
	}
	return
}

// Get from Redis
func (c *Client) Get(key string) (val []byte, ok bool, err error) {
	ok = true
	val, err = redis.Bytes(c.doCmd("GET", key))
	if err == redis.ErrNil {
		ok = false
		err = nil
	}
	return
}

// Get from Redis
func (c *Client) GetInt64(key string) (val int64, ok bool, err error) {
	ok = true
	val, err = redis.Int64(c.doCmd("GET", key))
	if err == redis.ErrNil {
		ok = false
		err = nil
	}
	return
}

func (c *Client) doCmd(cmd string, args ...any) (reply any, err error) {
	reply, err = c.Conn.Do(cmd, args...)
	if err == nil {
		return
	}
	if err != redis.ErrNil {
		err = errors.New(err.Error())
	}
	return
}

// Set RedisSet
func (c *Client) Set(key string, val any) (err error) {
	_, err = c.doCmd("SET", key, val)
	return
}

// SetEX RedisSetEX
func (c *Client) SetEX(key string, ttl int, val any) (err error) {
	_, err = c.doCmd("SETEX", key, ttl, val)
	return
}

// SetNX RedisSetNX
func (c *Client) SetNX(key string, val any) (ok bool, err error) {
	ok, err = redis.Bool(c.doCmd("SETNX", key, val))
	return
}

// SetEX SetEXNX 当ok == true新设置一个key, ok == false key已存在, 设置未生效
func (c *Client) SetEXNX(key string, ttl int, val any) (ok bool, err error) {
	out, err := c.doCmd("SET", key, val, "EX", ttl, "NX")
	if err != nil {
		return false, err
	}

	if out == nil {
		return false, nil
	}
	return true, nil
}

// SetPXNX 当ok == true新设置一个key, ok == false key已存在, 设置未生效
func (c *Client) SetPXNX(key string, ttl int, val any) (ok bool, err error) {
	out, err := c.doCmd("SET", key, val, "PX", ttl, "NX")
	if err != nil {
		return false, err
	}

	if out == nil {
		return false, nil
	}
	return true, nil
}

// IncrBy Redis incrby
func (c *Client) IncrBy(key string, delta int) (ret int, err error) {
	ret, err = redis.Int(c.doCmd("INCRBY", key, delta))
	return
}

// MSet Redis Mset
func (c *Client) MSet(kvs ...any) (err error) {
	_, err = c.doCmd("MSET", kvs...)
	return
}

// MGet Redis MGet
func (c *Client) MGet(keys ...string) (ret [][]byte, err error) {
	args := redis.Args{}.AddFlat(keys)
	ret, err = redis.ByteSlices(c.doCmd("MGET", args...))
	return
}

// GetMD5 get value from redis by key and return the md5sum of the value
func (c *Client) GetMD5(key string) (val []byte, md string, ok bool, err error) {
	val, ok, err = c.Get(key)
	if err != nil || ok == false {
		return
	}
	binMd := md5.Sum(val)
	md = hex.EncodeToString(binMd[:])
	return
}

// SetMD5 set the value to the key, success if the md equals to the original md5sum of the value
func (c *Client) SetMD5(key string, val any, md string) (err error) {
	if md == "" {
		md = "dummymd5"
	}
	_, err = c.doCmd("SETMD", key, val, md)
	return
}

//kv end

//keys start

// Exists Redis check whether key exists
func (c *Client) Exists(key string) (ok bool, err error) {
	ok, err = redis.Bool(c.doCmd("EXISTS", key))
	return
}

// Expire set expire of a key(in seconds)
func (c *Client) Expire(key string, ttl int) (err error) {
	_, err = c.doCmd("EXPIRE", key, ttl)
	return
}

// ExpireAt a key at a certain timestamp
func (c *Client) ExpireAt(key string, expire int) (err error) {
	_, err = c.doCmd("EXPIREAT", key, expire)
	return
}

// TTL 返回过期时间，-2表示key不存在，-1表示未设置过期时间
func (c *Client) TTL(key string) (int, error) {
	return redis.Int(c.doCmd("TTL", key))
}

// Del Del multiple keys in  Redis
func (c *Client) Del(keys ...string) (err error) {
	args := redis.Args{}.AddFlat(keys)
	_, err = c.doCmd("DEL", args...)
	return
}

//keys end

//hash start

// HGetAll Do a Redis HGETALL command
func (c *Client) HGetAll(key string) (val map[string]string, err error) {
	val, err = redis.StringMap(c.doCmd("HGETALL", key))

	return
}

// HMSet do a HMSET cmd in Redis
func (c *Client) HMSet(key string, val ...any) (err error) {
	args := make([]any, 1, len(val)+1)
	args[0] = key
	args = append(args, val...)
	_, err = c.doCmd("HMSET", args...)
	return
}

// HMGet do a HMGET
func (c *Client) HMGet(key string, fields ...any) (ret [][]byte, err error) {
	args := make([]any, 1, len(fields)+1)
	args[0] = key
	args = append(args, fields...)
	ret, err = redis.ByteSlices(c.doCmd("HMGET", args...))
	return
}

func (c *Client) HMGetInt64(key string, fields ...any) (ret []int64, err error) {
	args := make([]any, 1, len(fields)+1)
	args[0] = key
	args = append(args, fields...)
	ret, err = redis.Int64s(c.doCmd("HMGET", args...))
	return
}

// HMSetMap do a HMSET in Redis, the key-value pairs are stored in map
func (c *Client) HMSetMap(key string, val map[string]any) (err error) {
	args := redis.Args{}.Add(key).AddFlat(val)
	_, err = c.doCmd("HMSET", args...)
	return
}

func (c *Client) HMSetMapInt64(key string, val map[string]int64) (err error) {
	args := redis.Args{}.Add(key).AddFlat(val)
	_, err = c.doCmd("HMSET", args...)
	return
}

// HGet do a HGET cmd in Redis
func (c *Client) HGet(key, field string) (val []byte, err error) {
	val, err = redis.Bytes(c.doCmd("HGET", key, field))
	return
}

func (c *Client) HGetInt64(key, field string) (val int64, err error) {
	val, err = redis.Int64(c.doCmd("HGET", key, field))
	return
}

// HExists Redis check whether hash field exists
func (c *Client) HExists(key, field string) (ok bool, err error) {
	ok, err = redis.Bool(c.doCmd("HEXISTS", key, field))
	return
}

// HSet do a HSET cmd in Redis
func (c *Client) HSet(key, field string, val any) (err error) {
	_, err = c.doCmd("HSET", key, field, val)
	return
}

// HSetNX Redis HSETNX
func (c *Client) HSetNX(key, field string, val any) (ok bool, err error) {
	ok, err = redis.Bool(c.doCmd("HSETNX", key, field, val))
	return
}

// HIncrBy do a HINCRBY cmd in Redis
func (c *Client) HIncrBy(key, field string, delta int) (val int, err error) {
	val, err = redis.Int(c.doCmd("HINCRBY", key, field, delta))
	return
}

// HLen do a HLen cmd in Redis
func (c *Client) HLen(key string) (nLen int, err error) {
	nLen, err = redis.Int(c.doCmd("HLEN", key))
	return
}

// HDel hash delete
func (c *Client) HDel(key string, fields ...any) error {
	args := make([]any, 1, len(fields)+1)
	args[0] = key
	args = append(args, fields...)
	_, err := c.doCmd("HDEL", args...)
	return err
}

// HKeys 获取hash所有fields
func (c *Client) HKeys(key string) (keys []string, err error) {
	keys, err = redis.Strings(c.doCmd("HKEYS", key))
	return
}

// HScan hash scan
func (c *Client) HScan(key string, cursor, count int, pattern string) (newCursor int, pairs map[string]string, err error) {
	args := make([]any, 2)
	args[0] = key
	args[1] = cursor

	if pattern != "" {
		args = append(args, "MATCH")
		args = append(args, pattern)
	}

	if count > 0 {
		args = append(args, "COUNT")
		args = append(args, count)
	}

	bulk, err := redis.Values(c.doCmd("HSCAN", args...))
	if err == redis.ErrNil {
		err = nil
		newCursor = 0
		return
	}
	if err != nil {
		return
	}

	newCursor, err = redis.Int(bulk[0], nil)
	if err != nil {
		return
	}
	pairs, err = redis.StringMap(bulk[1], nil)
	return
}

// Hash Multi HGET
func (c *Client) MultiHGet(keys []string, field string) (vals [][]byte, err error) {
	sz := len(keys)
	for i := 0; i < sz; i++ {
		key := keys[i]
		if err = c.Send("HGET", key, field); err != nil {
			return
		}
	}
	if err = c.Flush(); err != nil {
		return
	}

	vals = make([][]byte, sz)
	for i := 0; i < sz; i++ {
		vals[i], err = redis.Bytes(c.Receive())
		if err == redis.ErrNil {
			err = nil
			continue
		}
		if err != nil {
			return
		}
	}
	return
}

//hash end

// Scan scan
func (c *Client) Scan(cursor, count int, patterns ...string) (newCursor int, keys []string, err error) {
	args := make([]any, 0, 3)
	args = append(args, cursor)
	if count > 0 {
		args = append(args, "COUNT", count)
	}
	if len(patterns) > 0 {
		for _, pattern := range patterns {
			args = append(args, "MATCH", pattern)
		}
	}

	bulk, err := redis.Values(c.doCmd("SCAN", args...))
	if err == redis.ErrNil {
		err = nil
		newCursor = 0
		return
	}
	if err != nil {
		return
	}

	newCursor, err = redis.Int(bulk[0], nil)
	if err != nil {
		return
	}
	keys, err = redis.Strings(bulk[1], nil)
	return
}

// Put returns the connection to the redis pool
func (c *Client) Put() {
	c.Close()
}

// SAdd 向集合中添加一个或者多个成员
func (c *Client) SAdd(key string, val ...any) error {
	args := make([]any, 1, len(val)+1)
	args[0] = key
	args = append(args, val...)
	_, err := c.doCmd("SADD", args...)
	return err
}

// SRem 移除集合中的一个或多个成员
func (c *Client) SRem(key string, val ...any) error {
	args := make([]any, 1, len(val)+1)
	args[0] = key
	args = append(args, val...)
	_, err := c.doCmd("SREM", args...)
	return err
}

// SPop 移除并返回集合中的一个随机元素
func (c *Client) SPop(key string) (val string, err error) {
	val, err = redis.String(c.doCmd("SPOP", key))
	if err == redis.ErrNil {
		err = nil
	}
	return
}

// LTrim 对一个列表进行修剪, 从0开始的闭区间
func (c *Client) LTrim(key string, start int, end int) (err error) {
	_, err = c.doCmd("LTRIM", key, start, end)
	return
}

// LPush 插入到list头
func (c *Client) LPush(key string, val ...any) error {
	args := make([]any, 1, len(val)+1)
	args[0] = key
	args = append(args, val...)
	_, err := c.doCmd("LPUSH", args...)
	return err
}

// LPop 移除并返回list头元素
func (c *Client) LPop(key string) (val []byte, err error) {
	val, err = redis.Bytes(c.doCmd("LPOP", key))
	if err == redis.ErrNil {
		err = nil
	}
	return
}

// RPush 插入到list尾
func (c *Client) RPush(key string, val ...any) error {
	args := make([]any, 1, len(val)+1)
	args[0] = key
	args = append(args, val...)
	_, err := c.doCmd("RPUSH", args...)
	return err
}

// LLen 返回list长度
func (c *Client) LLen(key string) (n int, err error) {
	n, err = redis.Int(c.doCmd("LLEN", key))
	if err == redis.ErrNil {
		err = nil
	}
	return
}

// LRange 返回list元素
func (c *Client) LRange(key string, start, stop int) ([]string, error) {
	return redis.Strings(c.doCmd("LRANGE", key, start, stop))
}

// SIsMember 判断member元素是否是集合key的成员
func (c *Client) SIsMember(key string, member any) (ok bool, err error) {
	ok, err = redis.Bool(c.doCmd("SISMEMBER", key, member))
	return
}

// SMembers 返回集合中的所有成员
// Deprecated: 不要使用这个函数，对于大的集合，会卡爆
func (c *Client) SMembers(key string) ([]string, error) {
	return redis.Strings(c.doCmd("SMEMBERS", key))
}

func (c *Client) SScan(key string, cursor, count int, pattern string) (nextCursor int, vals []string, err error) {
	args := make([]any, 2)
	args[0] = key
	args[1] = cursor

	if pattern != "" {
		args = append(args, "MATCH")
		args = append(args, pattern)
	}

	if count > 0 {
		args = append(args, "COUNT")
		args = append(args, count)
	}

	values, err := redis.Values(c.doCmd("SSCAN", args...))
	if err == redis.ErrNil {
		err = nil
		nextCursor = 0
		return
	} else if err != nil {
		return
	}

	nextCursor, err = redis.Int(values[0], nil)
	if err != nil {
		return
	}

	vals, err = redis.Strings(values[1], nil)
	return
}

// SRandMembers   随机返回集合中的成员
func (c *Client) SRandMembers(key string, count ...int) ([]string, error) {
	if len(count) == 1 {
		return redis.Strings(c.doCmd("SRANDMEMBER", key, count[0]))
	}

	buf, err := redis.Bytes(c.doCmd("SRANDMEMBER", key))
	if err != nil {
		return nil, err
	}

	return []string{string(buf)}, nil
}

// SRandMembers   返回集合中的总数
func (c *Client) SCard(key string) (int, error) {
	return redis.Int(c.doCmd("SCARD", key))
}

// ZAdd 向有序集合添加一个或多个成员，或者更新已存在成员的分数
func (c *Client) ZAdd(key string, val ...any) error {
	args := make([]any, 1, len(val)+1)
	args[0] = key
	args = append(args, val...)
	_, err := c.doCmd("ZADD", args...)
	return err
}

// ZCard 获取有序集合的成员数
func (c *Client) ZCard(key string) (int, error) {
	return redis.Int(c.doCmd("ZCARD", key))
}

// ZRange 按分数值递增返回有序集合成指定区间[start, end]内的成员(不包括score)
func (c *Client) ZRange(key string, start, end int) ([]string, error) {
	return redis.Strings(c.doCmd("ZRANGE", key, start, end))
}

// ZRangeOfSlices 按分数值递增返回有序集合成指定区间[start, end]内的成员(不包括score)
func (c *Client) ZRangeOfSlices(key string, start, end int) ([][]byte, error) {
	return redis.ByteSlices(c.doCmd("ZRANGE", key, start, end))
}

// ZRangeWithScores 按分数值递减返回有序集合成指定区间[start, end]内的成员
// 返回的[]string中，偶数下标的元素为member，该member对应的score值存储在后一个元素里面
func (c *Client) ZRangeWithScores(key string, start, end int) ([]string, error) {
	return redis.Strings(c.doCmd("ZRANGE", key, start, end, "WITHSCORES"))
}

// ZRevRangeWithScores 通过索引区间返回有序集合成指定区间[start, end]内的成员
// 返回的[]string中，偶数下标的元素为member，该member对应的score值存储在后一个元素里面
func (c *Client) ZRevRangeWithScores(key string, start, end int) ([]string, error) {
	return redis.Strings(c.doCmd("ZREVRANGE", key, start, end, "WITHSCORES"))
}

func (c *Client) ZRangeByScore(key string, start, end int) ([]string, error) {
	return redis.Strings(c.doCmd("ZRANGEBYSCORE", key, start, end))
}

// ZRevRange 通过索引区间返回有序集合成指定区间[start, end]内的成员
func (c *Client) ZRevRange(key string, start, end int) ([]string, error) {
	return redis.Strings(c.doCmd("ZREVRANGE", key, start, end))
}

// ZRem 移除有序集合中的一个或多个成员
func (c *Client) ZRem(key string, members ...any) error {
	args := make([]any, 1, len(members)+1)
	args[0] = key
	args = append(args, members...)
	_, err := c.doCmd("ZREM", args...)
	return err
}

func (c *Client) ZRemRangeByRank(key string, start, stop int64) (n int64, err error) {
	n, err = redis.Int64(c.doCmd("ZREMRANGEBYRANK", key, start, stop))
	return
}

// ZRank 返回有序集合中指定成员的索引
func (c *Client) ZRank(key, member string) (int, error) {
	return redis.Int(c.doCmd("ZRANK", key, member))
}

// ZRevRank 返回有序集合中指定成员的排名(从0开始)，有序集成员按分数值递减(从大到小)排序
func (c *Client) ZRevRank(key, member string) (int, error) {
	rs, err := redis.Int(c.doCmd("ZREVRANK", key, member))

	return rs, err
}

// ZScore 返回有序集中，成员的分数值
func (c *Client) ZScore(key, member string) (int, error) {
	rs, err := redis.Int(c.doCmd("ZSCORE", key, member))

	return rs, err
}

// ZScore 返回有序集中，成员的分数值
func (c *Client) ZIncrBy(key, member string, score int) (int, error) {
	return redis.Int(c.doCmd("ZINCRBY", key, score, member))
}

// Time 返回当前服务器时间,一个包含两个整型的列表： 第一个整数是当前时间(以 UNIX 时间戳格式表示)，而第二个整数是当前这一秒钟已经逝去的微秒数。
func (c *Client) Time() ([]int, error) {
	return redis.Ints(c.doCmd("TIME"))
}

// PutRedisClient puts a redis client into pool
func PutClient(cli *Client) {
	cli.Put()
}
