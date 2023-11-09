package tools

import (
	"context"
	"devops/common"
	"devops/config"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var RdbClient *redis.Client
var Ctx = context.Background()

func InitRedis() {
	addr := fmt.Sprintf("%s:%s", config.Conf.Redis.Host, config.Conf.Redis.Port)
	RdbClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: config.Conf.Redis.Password, // no password set
		DB:       config.Conf.Redis.DB,       // use default DB
		PoolSize: config.Conf.Redis.PoolSize,
	})
	_, err := RdbClient.Ping(Ctx).Result()
	if err != nil {
		common.Log.Error(err)
	}
}
