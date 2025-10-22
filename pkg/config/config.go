package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	API        ServerConfig     `mapstructure:"api"`
	Database   DatabaseConfig   `mapstructure:"database"`
	Kafka      KafkaConfig      `mapstructure:"kafka"`
	Redis      RedisConfig      `mapstructure:"redis"`
	Kubernetes KubernetesConfig `mapstructure:"kubernetes"`
	Monitoring MonitoringConfig `mapstructure:"monitoring"`
	Security   SecurityConfig   `mapstructure:"security"`
	ML         MLConfig         `mapstructure:"ml"`
	Logging    LoggingConfig    `mapstructure:"logging"`
}

// ServerConfig contains HTTP server configuration
type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
	TLS          TLSConfig     `mapstructure:"tls"`
}

// TLSConfig contains TLS configuration
type TLSConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	CertFile string `mapstructure:"cert_file"`
	KeyFile  string `mapstructure:"key_file"`
}

// DatabaseConfig contains database configuration
type DatabaseConfig struct {
	InfluxDB      InfluxDBConfig      `mapstructure:"influxdb"`
	Elasticsearch ElasticsearchConfig `mapstructure:"elasticsearch"`
}

// InfluxDBConfig contains InfluxDB configuration
type InfluxDBConfig struct {
	URL    string `mapstructure:"url"`
	Token  string `mapstructure:"token"`
	Org    string `mapstructure:"org"`
	Bucket string `mapstructure:"bucket"`
}

// ElasticsearchConfig contains Elasticsearch configuration
type ElasticsearchConfig struct {
	URLs     []string `mapstructure:"urls"`
	Username string   `mapstructure:"username"`
	Password string   `mapstructure:"password"`
	Index    string   `mapstructure:"index"`
}

// KafkaConfig contains Kafka configuration
type KafkaConfig struct {
	Brokers []string          `mapstructure:"brokers"`
	Topics  map[string]string `mapstructure:"topics"`
	GroupID string            `mapstructure:"group_id"`
	TLS     TLSConfig         `mapstructure:"tls"`
}

// RedisConfig contains Redis configuration
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// KubernetesConfig contains Kubernetes configuration
type KubernetesConfig struct {
	InCluster  bool   `mapstructure:"in_cluster"`
	ConfigPath string `mapstructure:"config_path"`
	Namespace  string `mapstructure:"namespace"`
}

// MonitoringConfig contains monitoring configuration
type MonitoringConfig struct {
	Prometheus PrometheusConfig `mapstructure:"prometheus"`
	Jaeger     JaegerConfig     `mapstructure:"jaeger"`
}

// PrometheusConfig contains Prometheus configuration
type PrometheusConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Path    string `mapstructure:"path"`
	Port    int    `mapstructure:"port"`
}

// JaegerConfig contains Jaeger configuration
type JaegerConfig struct {
	Enabled     bool   `mapstructure:"enabled"`
	ServiceName string `mapstructure:"service_name"`
	Endpoint    string `mapstructure:"endpoint"`
}

// SecurityConfig contains security configuration
type SecurityConfig struct {
	JWTSecret    string        `mapstructure:"jwt_secret"`
	JWTExpiry    time.Duration `mapstructure:"jwt_expiry"`
	RateLimit    RateLimitConfig `mapstructure:"rate_limit"`
	CORS         CORSConfig    `mapstructure:"cors"`
}

// RateLimitConfig contains rate limiting configuration
type RateLimitConfig struct {
	Enabled bool `mapstructure:"enabled"`
	RPS     int  `mapstructure:"rps"`
	Burst   int  `mapstructure:"burst"`
}

// CORSConfig contains CORS configuration
type CORSConfig struct {
	Enabled      bool     `mapstructure:"enabled"`
	Origins      []string `mapstructure:"origins"`
	Methods      []string `mapstructure:"methods"`
	Headers      []string `mapstructure:"headers"`
	Credentials  bool     `mapstructure:"credentials"`
}

// MLConfig contains ML service configuration
type MLConfig struct {
	Enabled     bool   `mapstructure:"enabled"`
	ServiceURL  string `mapstructure:"service_url"`
	ModelPath   string `mapstructure:"model_path"`
	BatchSize   int    `mapstructure:"batch_size"`
	Threshold   float64 `mapstructure:"threshold"`
}

// LoggingConfig contains logging configuration
type LoggingConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Output     string `mapstructure:"output"`
	Structured bool   `mapstructure:"structured"`
}

// Load loads configuration from file and environment variables
func Load(configPath string) (*Config, error) {
	config := &Config{}

	// Set defaults
	setDefaults()

	// Set config file path
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./configs")
		viper.AddConfigPath("../configs")
		viper.AddConfigPath(".")
	}

	// Enable environment variable support
	viper.AutomaticEnv()
	viper.SetEnvPrefix("SENTINEL")

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Unmarshal to struct
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate configuration
	if err := validate(config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	// Server defaults
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")
	viper.SetDefault("server.idle_timeout", "120s")

	// API defaults (same as server)
	viper.SetDefault("api.host", "0.0.0.0")
	viper.SetDefault("api.port", 8080)
	viper.SetDefault("api.read_timeout", "30s")
	viper.SetDefault("api.write_timeout", "30s")
	viper.SetDefault("api.idle_timeout", "120s")

	// Database defaults
	viper.SetDefault("database.influxdb.url", "http://localhost:8086")
	viper.SetDefault("database.influxdb.org", "sentinel-mesh")
	viper.SetDefault("database.influxdb.bucket", "metrics")
	viper.SetDefault("database.elasticsearch.urls", []string{"http://localhost:9200"})
	viper.SetDefault("database.elasticsearch.index", "sentinel-logs")

	// Kafka defaults
	viper.SetDefault("kafka.brokers", []string{"localhost:9092"})
	viper.SetDefault("kafka.group_id", "sentinel-mesh")
	viper.SetDefault("kafka.topics.metrics", "metrics")
	viper.SetDefault("kafka.topics.logs", "logs")
	viper.SetDefault("kafka.topics.alerts", "alerts")

	// Redis defaults
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.db", 0)

	// Kubernetes defaults
	viper.SetDefault("kubernetes.in_cluster", false)
	viper.SetDefault("kubernetes.namespace", "default")

	// Monitoring defaults
	viper.SetDefault("monitoring.prometheus.enabled", true)
	viper.SetDefault("monitoring.prometheus.path", "/metrics")
	viper.SetDefault("monitoring.prometheus.port", 9090)
	viper.SetDefault("monitoring.jaeger.enabled", false)
	viper.SetDefault("monitoring.jaeger.service_name", "sentinel-mesh")

	// Security defaults
	viper.SetDefault("security.jwt_expiry", "24h")
	viper.SetDefault("security.rate_limit.enabled", true)
	viper.SetDefault("security.rate_limit.rps", 100)
	viper.SetDefault("security.rate_limit.burst", 200)
	viper.SetDefault("security.cors.enabled", true)
	viper.SetDefault("security.cors.origins", []string{"*"})
	viper.SetDefault("security.cors.methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	viper.SetDefault("security.cors.headers", []string{"*"})

	// ML defaults
	viper.SetDefault("ml.enabled", true)
	viper.SetDefault("ml.service_url", "http://localhost:5000")
	viper.SetDefault("ml.batch_size", 100)
	viper.SetDefault("ml.threshold", 0.8)

	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
	viper.SetDefault("logging.output", "stdout")
	viper.SetDefault("logging.structured", true)
}

// validate validates the configuration
func validate(config *Config) error {
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", config.Server.Port)
	}

	if config.Database.InfluxDB.URL == "" {
		return fmt.Errorf("InfluxDB URL is required")
	}

	if len(config.Kafka.Brokers) == 0 {
		return fmt.Errorf("at least one Kafka broker is required")
	}

	if config.Security.JWTSecret == "" {
		// Generate a default secret if not provided (not recommended for production)
		config.Security.JWTSecret = os.Getenv("JWT_SECRET")
		if config.Security.JWTSecret == "" {
			config.Security.JWTSecret = "default-secret-change-in-production"
		}
	}

	return nil
}

// GetEnv returns environment variable value or default
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}