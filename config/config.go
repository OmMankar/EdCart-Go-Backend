package config

import (
	"encoding/json"
	"fmt"
	"main/models"
	"os"
)

type Config struct {
}

func (p *Config) Load(confPath string) models.Config {
	data, err := os.ReadFile(confPath)
	if err != nil {
		fmt.Println("err", err.Error())
	}
	var ret models.Config
	if err := json.Unmarshal(data, &ret); err != nil {
		fmt.Println("err", err.Error())
	}
	return ret
}
