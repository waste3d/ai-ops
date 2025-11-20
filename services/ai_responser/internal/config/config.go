// services/ai_responser/internal/infrastructure/config/config.go
package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Kafka KafkaConfig
	LLM   LLMConfig
}

type KafkaConfig struct {
	Brokers []string
	Topics  KafkaTopics
	GroupID string `mapstructure:"group_id"`
}

type KafkaTopics struct {
	Input  string
	Output string
}

type LLMConfig struct {
	Provider string
	APIKey   string `mapstructure:"api_key"`
}

func LoadConfig() (*Config, error) {
	// Viper будет искать файл config.yaml в корневой директории проекта
	viper.AddConfigPath(".")      // Указываем, где искать конфиг
	viper.SetConfigName("config") // Имя файла без расширения
	viper.SetConfigType("yaml")

	// Настраиваем чтение переменных окружения
	viper.AutomaticEnv()
	// Позволяет переменной окружения KAFKA_BROKERS соответствовать kafka.brokers
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// Приоритет переменной окружения для API ключа
	// Имя переменной будет AI_OPS_LLM_API_KEY
	if apiKey := viper.GetString("llm.api_key"); apiKey != "" {
		cfg.LLM.APIKey = apiKey
	}

	return &cfg, nil
}
