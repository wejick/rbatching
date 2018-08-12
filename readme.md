# rbatching
One pain point of doing batching in local memory is when the program crash all the batch will be gone. Rbatching provide persistent batching mechanism backed by redis, the persistence is as far as redis definition of persistence. 
The main idea is putting the queue to a redis list, pop and push it to other list before usage then remove when the usage is finished.
All the redis operation is atomic.

**This package is not thread safe, create different rbatching instance for each go routine**

# Usage
```golang
    var redisPool = &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6380")
		},
    }
    // Create new batcher
    // the uniqueness of the id is important to avoid racing condition
	batch := rbatching.NewRBatcher("batch1", 2, redisPool)

    // Enqueue data
    batch.Enqueue("1")
	batch.Enqueue("2")
	batch.Enqueue("3")
    batch.Enqueue("4")
    
    // get the batch batch
    val, err := batch.GetBatch()
    defer batch.CloseBatch() // this is important to close the batch session

    if err != nil {
		log.Println(err)
	}
	for _, val := range val {
		string := string(val.([]byte))
		log.Println(string)
    }
    
    // output:
    // 2
    // 1
    

```
