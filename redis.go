package rbatching

import (
	"github.com/gomodule/redigo/redis"
)

const luaGetBatch = `
for i=0,ARGV[3] do
 redis.call("RPOPLPUSH",ARGV[1],ARGV[2])
end
local element = redis.call("LRANGE",ARGV[2],0,-1)

return element`

func (r *RedisBatcher) initScript() (err error) {
	con := r.redispool.Get()
	defer con.Close()
	r.getBatchScript = redis.NewScript(0, luaGetBatch)
	err = r.getBatchScript.Load(con)

	return
}

func (r *RedisBatcher) redisEnqueue(element interface{}) (err error) {
	con := r.redispool.Get()
	defer con.Close()
	_, err = con.Do("LPUSH", r.queueName, element)
	return
}

func (r *RedisBatcher) redisGetBatch() (batch []interface{}, err error) {
	con := r.redispool.Get()
	defer con.Close()

	reply, err := r.getBatchScript.Do(con, r.queueName, r.batchName, r.maxBatch-1)

	batch = reply.([]interface{})

	return
}

func (r *RedisBatcher) redisCloseBatch() (err error) {
	con := r.redispool.Get()
	defer con.Close()

	_, err = con.Do("del", r.batchName)

	return
}
