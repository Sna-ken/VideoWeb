package config

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type ModelConfig struct {
	APIKey        string `mapstructure:"api_key"`
	BaseURL       string `mapstructure:"base_url"`
	DecideModel   string `mapstructure:"decide_model"`
	GenerateModel string `mapstructure:"generate_model"`
}

type AgentConfig struct {
	Name    string `mapstructure:"name"`
	UserID  string `mapstructure:"user_id"`
	Persona string `mapstructure:"persona"`
}

type TriggerConfig struct {
	SilenceTimeout       int64   `mapstructure:"silence_timeout"`
	RateLimit            int64   `mapstructure:"rate_limit"`
	EmotionThreshold     float64 `mapstructure:"emotion_threshold"`
	UnknownWordThreshold int64   `mapstructure:"unknown_word_threshold"`
	SilenceTimerKey      string  `mapstructure:"silence_timer_key"`
	RateLimitKey         string  `mapstructure:"rate_limit_key"`
}

type MemoryConfig struct {
	WindowSize int64 `mapstructure:"window_size"`
}

type SearchConfig struct {
	SearchAPIKey string `mapstructure:"search_api_key"`
}

var Model = &ModelConfig{}
var Agent = &AgentConfig{}
var Trigger = &TriggerConfig{}
var Search = &SearchConfig{}
var Memory = &MemoryConfig{}

func InitAIConfig() {
	log.Printf("InitAIConfig 被调用")

	viper.SetConfigFile("./config/ai_config.yml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("read config fail: %v", err)
	}

	viper.WatchConfig()

	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("config file changed: %s", e.Name)

		if err := viper.ReadInConfig(); err != nil {
			log.Printf("read config fail: %v", err)
			return
		}
		log.Printf("config reloaded successfully")
	})

	if err := viper.UnmarshalKey("model", Model); err != nil {
		log.Fatalf("unable to decode model config: %v", err)
	}
	if err := viper.UnmarshalKey("agent", Agent); err != nil {
		log.Fatalf("unable to decode agent config: %v", err)
	}
	if err := viper.UnmarshalKey("trigger", Trigger); err != nil {
		log.Fatalf("unable to decode trigger config: %v", err)
	}
	if err := viper.UnmarshalKey("memory", Memory); err != nil {
		log.Fatalf("unable to decode memory config: %v", err)
	}

	log.Printf("Trigger配置读取结果: %+v", Trigger)
}
