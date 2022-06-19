package redis

import "os"

var RedisAddress = os.Getenv("REDIS_DSN")
