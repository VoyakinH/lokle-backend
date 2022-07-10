package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type PostgresConfig struct {
	User     string
	Password string
	Port     string
	Host     string
	DBName   string
}

type MailerConfig struct {
	Email    string
	Password string
}

type TimeoutsConfig struct {
	WriteTimeout   time.Duration
	ReadTimeout    time.Duration
	ContextTimeout time.Duration
}

var (
	Lokle    ServerConfig
	Redis    RedisConfig
	Postgres PostgresConfig
	Mailer   MailerConfig
	Timeouts TimeoutsConfig
)

func SetConfig() {
	viper.SetConfigFile("config.json")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	Lokle = ServerConfig{
		Port: viper.GetString(`lokle.port`),
	}

	Redis = RedisConfig{
		Addr:     viper.GetString(`redis.address`),
		Password: viper.GetString(`redis.password`),
		DB:       viper.GetInt(`redis.db`),
	}

	Postgres = PostgresConfig{
		Port:     viper.GetString(`postgres.port`),
		Host:     viper.GetString(`postgres.host`),
		User:     viper.GetString(`postgres.user`),
		Password: viper.GetString(`postgres.pass`),
		DBName:   viper.GetString(`postgres.name`),
	}

	Mailer = MailerConfig{
		Email:    viper.GetString(`mailer.email`),
		Password: viper.GetString(`mailer.password`),
	}

	Timeouts = TimeoutsConfig{
		WriteTimeout:   5 * time.Second,
		ReadTimeout:    5 * time.Second,
		ContextTimeout: time.Second * 2,
	}
}
