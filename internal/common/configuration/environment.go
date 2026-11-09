package configuration

import (
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

type EnvironmentVariablesI struct {
	Port                string
	RedisCluster        bool
	RedisHost           string
	RedisPort           string
	RedisPassword       string
	RedisDB             int
	DBFullHost          string
	DBUser              string
	DBPassword          string
	DBHost              string
	DBPort              string
	NatsUrl             string
	DBClient            string
	Mode                string
	NatsJetStream       string
	ValKeySentinelHost  string
	ValKeyMasterName    string
	DefaultCurrencyCode string
	GeocodingEnabled    bool
	GeocodingTimeoutSec int
}

func GetEnvironmentVariable(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func LoadEnvironmentVariables() EnvironmentVariablesI {
	port := GetEnvironmentVariable("PORT", "3259")
	clusterMode := GetEnvironmentVariable("REDIS_CLUSTER", "false")
	hostRedis := GetEnvironmentVariable("REDIS_HOST", "localhost")
	portRedis := GetEnvironmentVariable("REDIS_PORT", "6379")
	passwordRedis := GetEnvironmentVariable("REDIS_PASSWORD", "")
	dbRedis, err := strconv.Atoi(GetEnvironmentVariable("REDIS_DB", "0"))
	if err != nil {
		log.Fatal(err)
	}
	dbFullHost := GetEnvironmentVariable("DB_FULL_HOST", "")
	dbPass := GetEnvironmentVariable("DB_PASSWORD", "")
	dbUser := GetEnvironmentVariable("DB_USER", "")
	dbHost := GetEnvironmentVariable("DB_HOST", "localhost")
	dbPort := GetEnvironmentVariable("DB_PORT", "27017")
	client := GetEnvironmentVariable("CLIENT", "27017")
	natsHost := GetEnvironmentVariable("NATS_HOST", "localhost")
	natsPort := GetEnvironmentVariable("NATS_PORT", "4223")
	natsUrl := fmt.Sprintf(`http://%v:%v`, natsHost, natsPort)
	mode := GetEnvironmentVariable("MODE", "")
	valKeySentinelHost := GetEnvironmentVariable("VALKEY_SENTINEL_HOST", "")
	valKeyMasterName := GetEnvironmentVariable("VALKEY_MASTER_NAME", "")
	natsJetStream := GetEnvironmentVariable("NATS_JETSTREAM", "false")
	defaultCurrencyCode := GetEnvironmentVariable("DEFAULT_CURRENCY_CODE", "AED")
	geocodingEnabled := GetEnvironmentVariable("GEOCODING_ENABLED", "true") == "true"
	geocodingTimeoutSec, _ := strconv.Atoi(GetEnvironmentVariable("GEOCODING_TIMEOUT_SEC", "5"))

	return EnvironmentVariablesI{
		Port:                port,
		RedisCluster:        clusterMode == "true",
		RedisHost:           hostRedis,
		RedisPort:           portRedis,
		RedisPassword:       passwordRedis,
		RedisDB:             dbRedis,
		DBFullHost:          dbFullHost,
		DBUser:              dbUser,
		DBPassword:          dbPass,
		DBHost:              dbHost,
		DBPort:              dbPort,
		NatsUrl:             natsUrl,
		DBClient:            client,
		Mode:                mode,
		NatsJetStream:       natsJetStream,
		ValKeySentinelHost:  valKeySentinelHost,
		ValKeyMasterName:    valKeyMasterName,
		DefaultCurrencyCode: defaultCurrencyCode,
		GeocodingEnabled:    geocodingEnabled,
		GeocodingTimeoutSec: geocodingTimeoutSec,
	}
}

var EnvironmentVariables EnvironmentVariablesI

func init() {
	EnvironmentVariables = LoadEnvironmentVariables()
}
