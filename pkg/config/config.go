package config

import (
	"fmt"
	"parsing-service/pkg/logger"

	"github.com/spf13/viper"
)

type Configuration struct {
	KafkaConfig         KafkaConfig
	RedisConfig         RedisConfig
	KafkaTopicsConfig   KafkaTopicsConfig
	BusinessLogicConfig BusinessLogicConfig
	// AuthConfig AuthConfig
	// APIClientConfig     APIClientConfig
	// AWSCredentials      AWSCredentials
	// PushDataCache       PushDataCache
	Environment Environment
}

// type AuthConfig struct {

// }

type KafkaConfig struct {
	KafkaBootstrapServers string
	KafkaSecurityProtocol string
	KafkaSaslUsername     string
	KafkaSaslPassword     string
	KafkaSaslMechanism    string
}

type RedisConfig struct {
	Host              string
	Port              string
	Password          string
	REDIS_DISABLE_TLS bool
	DefaultDBNumber   int
}

type KafkaTopicsConfig struct {
	Fetch_Data_Kafka_Topic_GROUP_ID string
	Fetch_Data_Kafka_Topic_NAME     string
}

type BusinessLogicConfig struct {
}

// type APIClientConfig struct {
// }

// type AWSCredentials struct {
// }

type Environment struct {
	ENVIRONMENT string
}

// type PushDataCache struct {
// }

func NewConfiguration() (*Configuration, error) {
	cfg, err := SetupConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}
	return cfg, nil
}

func SetupConfig() (*Configuration, error) {
	logger := logger.GetLogger()

	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		logger.Infof("Reading config from env variables")
		viper.AutomaticEnv()
	}

	configuration := &Configuration{
		KafkaConfig: KafkaConfig{
			KafkaBootstrapServers: viper.GetString("KAFKA_BOOTSTRAP_SERVERS"),
			KafkaSecurityProtocol: viper.GetString("KAFKA_SECURITY_PROTOCOL"),
			KafkaSaslUsername:     viper.GetString("KAFKA_SASL_USERNAME"),
			KafkaSaslPassword:     viper.GetString("KAFKA_SASL_PASSWORD"),
			KafkaSaslMechanism:    viper.GetString("KAFKA_SASL_MECHANISM"),
		},

		RedisConfig: RedisConfig{
			Host:                   viper.GetString("REDIS_HOST"),
			Port:                   viper.GetString("REDIS_PORT"),
			Password:               viper.GetString("REDIS_PASSWORD"),
			REDIS_DISABLE_TLS:      viper.GetBool("REDIS_DISABLE_TLS"),
			DefaultDBNumber:        viper.GetInt("REDIS_DEFAULT_DB_NUMBER"),
		},

		KafkaTopicsConfig: KafkaTopicsConfig{
			Fetch_Data_Kafka_Topic_GROUP_ID:  viper.GetString("Fetch_Data_Kafka_Topic_GROUP_ID"),
			Fetch_Data_Kafka_Topic_NAME: viper.GetString("Fetch_Data_Kafka_Topic_NAME"),
		},

		BusinessLogicConfig: BusinessLogicConfig{

		},

		// APIClientConfig: APIClientConfig{
		// },

		// AuthConfig: AuthConfig{
		// },

		// AWSCredentials: AWSCredentials{
		// },

		// PushDataCache: PushDataCache{
		// },

		Environment: Environment{
			ENVIRONMENT: viper.GetString("ENVIRONMENT_TYPE"),
		},
	}

	if err := viper.Unmarshal(configuration); err != nil {
		logger.Errorf("Unmarshal config error: %v", err)
		return nil, err
	}

	return configuration, nil
}