package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var (
	JmeterAgentServicesKey = "jmeterAgentServices"
)

func AddAgentService(ctx context.Context, ip string, port string) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer rdb.Close()
	key := fmt.Sprintf("%s:%s", ip, port)
	// SET 元素重复自动覆盖
	_, err := rdb.SAdd(ctx, JmeterAgentServicesKey, key).Result()
	if err != nil {
		return
	}
}

func GetAgentService(ctx context.Context) ([]string, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer rdb.Close()
	return rdb.SMembers(ctx, JmeterAgentServicesKey).Result()
}
