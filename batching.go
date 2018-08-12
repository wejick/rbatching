package rbatching

import (
	"errors"
	"log"

	"github.com/gomodule/redigo/redis"
)

// RedisBatcher provide persistence batching mechanism
// backed by redis, the persistence is as far as redis
// definition of persistance
type RedisBatcher struct {
	redispool *redis.Pool
	uniqueID  string
	queueName string
	batchName string
	maxBatch  int

	openSession bool

	getBatchScript *redis.Script
}

const (
	queueNamePrefix = "rbatcher_queue_" //main queue
	wipNamePrefix   = "rbatcher_wip_"   //inprogress batch
)

// NewRBatcher create redis batcher instance
// the uniqueness of the id is important to avoid racing condition with other batcher
func NewRBatcher(uniqueID string, maxBatch int, redispool *redis.Pool) (batcher *RedisBatcher) {
	batcher = &RedisBatcher{
		redispool: redispool,
		uniqueID:  uniqueID,
		maxBatch:  maxBatch,
		queueName: queueNamePrefix + uniqueID,
		batchName: wipNamePrefix + uniqueID,
	}

	err := batcher.initScript()
	if err != nil {
		log.Println(err)
	}

	return
}

// Enqueue add element to the queue
func (r *RedisBatcher) Enqueue(element interface{}) (err error) {
	err = r.redisEnqueue(element)
	if err != nil {
		err = errors.New("[rbatching] couldn't put element into queue" + err.Error())
	}
	return
}

// GetBatch returns elements from the queue as many as maxBatch.
// This will open a batch session which need to be closed before getting another batch.
// When the app crash before calling CloseBatch, the elements will be included on the
// next GetBatch call
func (r *RedisBatcher) GetBatch() (elements []interface{}, err error) {
	if r.openSession {
		return nil, errors.New("[rbatching] Couldn't open batch, there is another open session")
	}
	elements, err = r.redisGetBatch()
	if err == nil {
		r.openSession = true
	}

	return
}

// CloseBatch close a batch session
func (r *RedisBatcher) CloseBatch() {
	err := r.redisCloseBatch()
	if err != nil {
		log.Println("[rbatching] Couldn't close batch ", err)
		return
	}
	r.openSession = false
}
