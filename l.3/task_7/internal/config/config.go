package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/wb-go/wbf/config"
)

var Cfg = initConfig(getConfigPath())

func initConfig(path string) *Config {
	wbfConfig := config.New()

	err := wbfConfig.LoadConfigFiles(path)
	if err != nil {
		log.Fatal("could not read config file: ", err)
	}

	var cfg Config
	if err := wbfConfig.Unmarshal(&cfg); err != nil {
		log.Fatal("could not parse config file: ", err)
	}

	err = godotenv.Load(getEnvPath())
	if err != nil {
		log.Fatal("could not load .env file: ", err)
	}

	value, _ := os.LookupEnv("DB_PASSWORD")
	cfg.Postgres.Password = value

	secret, _ := os.LookupEnv("SECRET")
	cfg.HttpServer.Secret = secret

	return &cfg
}

func getConfigPath() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	root := filepath.Dir(filepath.Dir(dir))
	return filepath.Join(root, "config", "config.yaml")
}

func getEnvPath() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	root := filepath.Dir(filepath.Dir(dir))
	return filepath.Join(root, "/", ".env")
}
