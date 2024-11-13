package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Mysql     MysqlConf       `json:"Mysql"`
	RedisConf redis.RedisConf `json:"RedisConf"`
}

// MysqlConf mysql配置
type MysqlConf struct {
	Casino MysqlInsConf `json:"Casino"`
}

type MysqlInsConf struct {
	Host   string `json:"Host"`
	Port   string `json:"Port"`
	User   string `json:"User"`
	Pwd    string `json:"Pwd"`
	DbName string `json:"DbName"`
}
