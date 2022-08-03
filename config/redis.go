package config

import "github.com/go-redis/redis"

var RS *redis.Client

func RedisConfig() {

	opt, err := redis.ParseURL("redis://@0.0.0.0:6379")
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(opt)

	RS = client
}
