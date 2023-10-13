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
	Port   string `mapstructure:"DATABASE_PORT"`
	Host   string `mapstructure:"DATABASE_HOST"`
	User   string `mapstructure:"DATABASE_USER"`
	Pass   string `mapstructure:"DATABASE_PASSWORD"`
	DB     string `mapstructure:"DATABASE_DB"`
}

type MessageBroker struct {
	User            string `mapstructure:"RABBITMQ_USER"`
	Pass            string `mapstructure:"RABBITMQ_PASS"`
	Port            string `mapstructure:"RABBITMQ_PORT"`
	Host            string `mapstructure:"RABBITMQ_HOST"`
	ExchangeAluno   string `mapstructure:"RABBITMQ_EXCHANGE_ALUNO"`
	QueueAluno      string `mapstructure:"RABBITMQ_QUEUE_ALUNO"`
	RoutingKeyAluno string `mapstructure:"RABBITMQ_ROUTINGKEY_ALUNO"`
}

type Config struct {
	API           API           `mapstructure:",squash"`
	Database      Database      `mapstructure:",squash"`
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
