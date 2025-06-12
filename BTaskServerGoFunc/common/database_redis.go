package common

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"strconv"
)

var (
	RDS *redis.Client
	ctx = context.Background()
)

func InitRedis() bool {
	addr := viper.GetString("datasouce_redis.host")
	port, _ := strconv.Atoi(viper.GetString("datasouce_redis.port"))
	password := viper.GetString("datasouce_redis.password")
	defultDB, _ := strconv.Atoi(viper.GetString("datasouce_redis.DB"))

	rd := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", addr, port),
		Password: password, // no password set
		DB:       defultDB, // use default DB
	})

	if _, err := rd.Ping(ctx).Result(); err != nil {
		panic("redis连接失败" + err.Error())
		return false
	}

	RDS = rd

	fmt.Println("redis连接成功")
	return true
}

func GetRedis() *redis.Client {
	return RDS
}
