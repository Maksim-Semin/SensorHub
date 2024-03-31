package resource

import (
	"github.com/BurntSushi/toml"
	"log"
)

type Config struct {
	Port       string
	Baud       int
	DBUser     string
	DBPassword string
	DBName     string
}

var Cfg Config

func init() {
	_, err := toml.DecodeFile("pkg/resource/config.toml", &Cfg)
	if err != nil {
		log.Fatal(err)
	}
}
