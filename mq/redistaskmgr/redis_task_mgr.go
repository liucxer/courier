package redistaskmgr

import (
	"github.com/liucxer/courier/mq"
	"github.com/gomodule/redigo/redis"
	"google.golang.org/protobuf/proto"
)

type RedisController interface {
	// key prefix
	Prefix(key string) string
	// get redis connect
	Get() redis.Conn
}

func RedisOperatorFromPool(pool *redis.Pool) *RedisOperator {
	return &RedisOperator{
		Pool: pool,
	}
}

type RedisOperator struct {
	*redis.Pool
}

func (RedisOperator) Prefix(key string) string {
	return key
}

func NewRedisTaskMgr(redisOp RedisController) *RedisTaskMgr {
	return &RedisTaskMgr{
		c: redisOp,
	}
}

type RedisTaskMgr struct {
	c RedisController
}

var _ mq.TaskMgr = (*RedisTaskMgr)(nil)

func (mgr *RedisTaskMgr) Exec(run func(c redis.Conn) error) error {
	c := mgr.c.Get()
	defer c.Close()

	if err := run(c); err != nil && err != redis.ErrNil {
		return err
	}

	return nil
}

func (mgr *RedisTaskMgr) Destroy(channel string) error {
	return mgr.Exec(func(c redis.Conn) error {
		_, err := QDEL.Do(c, mgr.Key(channel))
		return err
	})
}

func (mgr *RedisTaskMgr) Key(channel string) string {
	return mgr.c.Prefix("mq::" + channel)
}

func (mgr *RedisTaskMgr) Push(channel string, task *mq.Task) error {
	data, err := proto.Marshal(task)
	if err != nil {
		return err
	}
	return mgr.Exec(func(c redis.Conn) error {
		_, err := QPUSH.Do(c, mgr.Key(channel), task.Id, data)
		return err
	})
}

func (mgr *RedisTaskMgr) Shift(channel string) (task *mq.Task, err error) {
	err = mgr.Exec(func(c redis.Conn) error {
		result, err := redis.Bytes(QSHIFT.Do(c, mgr.Key(channel)))
		if len(result) > 0 {
			task = &mq.Task{}
			if err := proto.Unmarshal([]byte(result), task); err != nil {
				return err
			}
		}
		return err
	})
	return
}

func (mgr *RedisTaskMgr) Remove(channel, id string) error {
	return mgr.Exec(func(c redis.Conn) error {
		_, err := QREM.Do(c, mgr.Key(channel), id)
		return err
	})
}

var (
	QPUSH = redis.NewScript(1 /* language=lua */, `
local keySet = KEYS[1] .. '::set'
local keyQueue = KEYS[1] .. '::queue'
local id = ARGV[1]
local data = ARGV[2]
local success = redis.call('HSETNX', keySet, id, data)

if (success == 1) then 
	redis.call('RPUSH', keyQueue, id)
end
`)

	QSHIFT = redis.NewScript(1 /* language=lua */, `
local keySet = KEYS[1] .. '::set'
local keyQueue = KEYS[1] .. '::queue'
local id = redis.call('LPOP', keyQueue)
local data

if (id == nil or id == false) then
	return id
end

data = redis.call('HGET', keySet, id)
redis.call('HDEL', keySet, id)
return data
`)

	QDEL = redis.NewScript(1 /* language=lua */, `
local keySet = KEYS[1] .. '::set'
local keyQueue = KEYS[1] .. '::queue'

redis.call('DEL', keySet)
redis.call('DEL', keyQueue)
`)

	QREM = redis.NewScript(1 /* language=lua */, `
local keySet = KEYS[1] .. '::set'
local id = ARGV[1]

redis.call('HDEL', keySet, id)
`)
)
