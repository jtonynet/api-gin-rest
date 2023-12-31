package config

import (
	"github.com/spf13/viper"
)

type API struct {
	Name       string `mapstructure:"API_NAME"`
	Port       string `mapstructure:"API_PORT"`
	TagVersion string `mapstructure:"API_TAG_VERSION"`
	Env        string `mapstructure:"API_ENV"`

	RetryMaxElapsedTimeInMs int `mapstructure:"API_RETRY_MAX_ELAPSED_TIME_IN_MS"`

	FeatureFlags FeatureFlags `mapstructure:",squash"`
}

type Database struct {
	Driver string `mapstructure:"DATABASE_DRIVER"`

	Port string `mapstructure:"DATABASE_PORT"`
	Host string `mapstructure:"DATABASE_HOST"`
	User string `mapstructure:"DATABASE_USER"`
	Pass string `mapstructure:"DATABASE_PASSWORD"`
	DB   string `mapstructure:"DATABASE_DB"`
}

type Cache struct {
	Strategy string `mapstructure:"CACHE_STRATEGY"` // redis

	Pass       string `mapstructure:"REDIS_PASSWORD"`
	Port       string `mapstructure:"REDIS_PORT"`
	Host       string `mapstructure:"REDIS_HOST"`
	DB         int    `mapstructure:"REDIS_DB"`
	Protocol   int    `mapstructure:"REDIS_PROTOCOL"`
	Expiration int    `mapstructure:"REDIS_EXPIRATION_DEFAULT_IN_MS"`
}

type MessageBroker struct {
	Strategy string `mapstructure:"MESSAGE_BROKER_STRATEGY"` // rabbitMQ|kafka in future

	User string `mapstructure:"RABBITMQ_USER"`
	Pass string `mapstructure:"RABBITMQ_PASS"`
	Port string `mapstructure:"RABBITMQ_PORT"`
	Host string `mapstructure:"RABBITMQ_HOST"`

	Exchange     string `mapstructure:"RABBITMQ_EXCHANGE_ALUNO"`
	ExchangeType string `mapstructure:"RABBITMQ_EXCHANGE_ALUNO_TYPE"` //direct|fanout|topic|x-custom
	Queue        string `mapstructure:"RABBITMQ_QUEUE_ALUNO"`
	RoutingKey   string `mapstructure:"RABBITMQ_ROUTINGKEY_ALUNO"`
	ConsumerTag  string `mapstructure:"RABBITMQ_CONSUMER_TAG"`

	// DeadLetter
	ExchangeDL     string `mapstructure:"RABBITMQ_EXCHANGE_DEAD_LETTER"`
	ExchangeTypeDL string `mapstructure:"RABBITMQ_EXCHANGE_DEAD_LETTER_TYPE"` //direct|fanout|topic|x-custom
	QueueDL        string `mapstructure:"RABBITMQ_QUEUE_DEAD_LETTER"`
	RoutingKeyDL   string `mapstructure:"RABBITMQ_ROUTINGKEY_DEAD_LETTER"`
	ConsumerTagDL  string `mapstructure:"RABBITMQ_CONSUMER_TAG_DEAD_LETTER"`

	AutoReconnectEnable              bool `mapstructure:"RABBITMQ_AUTO_RECONNECT_ENABLED"`
	AutoReconnectRetryMaxElapsedInMs int  `mapstructure:"RABBITMQ_AUTO_RECONNECT_RETRY_MAX_ELAPSED_IN_MS"`

	MaxAttempts            int32 `mapstructure:"RABBITMQ_MAX_ATTEMPTS_CONSUME_INT"`
	ReliableMessagesEnable bool  `mapstructure:"RABBITMQ_RELIABLE_MESSAGES_ENABLED"` //Wait for the publisher confirmation before exiting
}

type Config struct {
	API           API           `mapstructure:",squash"`
	Database      Database      `mapstructure:",squash"`
	Cache         Cache         `mapstructure:",squash"`
	MessageBroker MessageBroker `mapstructure:",squash"`
}

func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
