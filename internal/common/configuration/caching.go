package configuration

import (
	"strings"

	"git.larsa.io/mahdawi/microservices-commons.git/cache"
)

var Rdb cache.Cache

func InitRedis() {
	host := EnvironmentVariables.RedisHost
	port := EnvironmentVariables.RedisPort
	password := EnvironmentVariables.RedisPassword
	valKeySentinelHost := EnvironmentVariables.ValKeySentinelHost
	valKeyMasterName := EnvironmentVariables.ValKeyMasterName
	db := EnvironmentVariables.RedisDB
	uri := host + ":" + port
	var redisCache cache.Cache
	cacheCluserMode := EnvironmentVariables.RedisCluster
	if cacheCluserMode {
		redisCache = cache.NewRedisClusterCache()
		errCache := redisCache.Connect(uri, password, db, []string{})
		if errCache != nil {
			return
		}
	} else {
		var host []string = []string{}
		redisCache = cache.NewRedisCache()
		if valKeySentinelHost != "" {
			uri = valKeyMasterName
			host = strings.Split(valKeySentinelHost, ", ")
		}
		errCache := redisCache.Connect(uri, password, db, host)
		if errCache != nil {
			return
		}
	}

	Rdb = redisCache

}
