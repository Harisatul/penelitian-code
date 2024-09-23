package shared

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"log/slog"
	"os"
)

type pgConfig struct {
	Host     string `yaml:"host" json:"host"`
	Port     uint   `yaml:"port" json:"port"`
	User     string `yaml:"user" json:"user"`
	Password string `yaml:"password" json:"password"`
	DBName   string `yaml:"db_name" json:"db_name"`
	SslMode  string `yaml:"ssl_mode" json:"ssl_mode"`
	MinConn  uint   `yaml:"min_conn" json:"min_conn"`
	MaxConn  uint   `yaml:"max_conn" json:"max_conn"`
}

func (p pgConfig) ConnStr() string {
	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%d database=%s sslmode=%s pool_min_conns=%d pool_max_conns=%d pool_max_conn_lifetime=5s",
		p.User, p.Password, p.Host, p.Port, p.DBName, p.SslMode, p.MinConn, p.MaxConn,
	)
}

type serverConfig struct {
	Port         uint `yaml:"port" json:"port"`
	ReadTimeout  uint `yaml:"read_timeout" json:"read_timeout"`
	WriteTimeout uint `yaml:"write_timeout" json:"write_timeout"`
}

func (l serverConfig) Addr() string {
	return fmt.Sprintf(":%d", l.Port)
}

type redisConfig struct {
	Addr    string `yaml:"addr" json:"addr"`
	MinConn uint   `yaml:"min_conn" json:"min_conn"`
}

type kafkaConfig struct {
	Broker string `yaml:"broker" json:"broker"`
}

type kafkaTopicConfig struct {
	CompleteOrder string `yaml:"complete_order" json:"complete_order"`
}

type kafkaGroupConfig struct {
	CompleteOrder string `yaml:"complete_order" json:"complete_order"`
}

type config struct {
	Server      serverConfig     `yaml:"server" json:"server"`
	DBConfig    pgConfig         `yaml:"db" json:"db"`
	RedisConfig redisConfig      `yaml:"redis" json:"redis"`
	Kafka       kafkaConfig      `yaml:"kafka" json:"kafka"`
	KafkaTopic  kafkaTopicConfig `yaml:"kafka_topic" json:"kafka_topic"`
	KafkaGroup  kafkaGroupConfig `yaml:"kafka_group" json:"kafka_group"`
}

func loadConfigFromReader(r io.Reader, c *config) error {
	return yaml.NewDecoder(r).Decode(c)
}

func loadConfigFromFile(fn string, c *config) error {
	_, err := os.Stat(fn)

	if err != nil {
		return err
	}

	f, err := os.Open(fn)

	if err != nil {
		return err
	}

	defer f.Close()

	return loadConfigFromReader(f, c)
}

func LoadConfig(fn string) config {
	cfg := config{}
	err := loadConfigFromFile(fn, &cfg)
	if err != nil {
		panic(err)
	}

	slog.Debug("config loaded", slog.Any("config", cfg))
	return cfg
}
