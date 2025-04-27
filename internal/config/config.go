package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPserver struct {
	Addr string `yaml:"address" env-required:"true"`
}

//  env-default:"production
/*

Host        string `yaml:"host" env-required:"true"`
	Port        string `yaml:"port" env-required:"true"`
	Password    string `yaml:"password" env-required:"true"`
	Dbname      string `yaml:"dbName" env-required:"true"`
*/
type Config struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	DatabaseURL string `yaml:"database_url" env-required:"true"`
	HTTPserver  `yaml:"http_server" env-required:"true"`
}

func MustLoad() *Config {

	var configPath string
	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flags := flag.String("config", "", "path to the configuration file")
		flag.Parse()
		configPath = *flags
		if configPath == "" {
			log.Fatal("config path is not set")
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file not exists..!%s", configPath)
	}

	// _, err := os.Stat(configPath)
	// if os.IsNotExist(err) {
	// 	log.Fatalf("config file doesn't exist %s", configPath)
	// }

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("cant read config file %s", err.Error())
	}

	return &cfg
}
