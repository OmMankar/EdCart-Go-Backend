package models

import "fmt"

type Config struct {
	DbName         string `json:"db_name"`
	UserCollection string `json:"user_collection"`
	CardCollection string `json:"card_collection"`
	DbIp           string `json:"db_ip"`
	DbPort         int    `json:"db_port"`
	DbCluster      string `json:"db_cluster"`
}

func (p *Config) DNS() string {
	return fmt.Sprintf("mongodb://%v:%v/%v", p.DbIp, p.DbPort, p.DbCluster)
}
