package config

import "time"

var C *Config

type Config struct {
	Server       Server  `yaml:"server"`
	Log          Log     `yaml:"log"`
	Mysql        Mysql   `yaml:"mysql"`
	Swagger      Swagger `yaml:"swagger"`
	Video        Video   `yaml:"video"`
	Ffmpeg       Ffmpeg  `yaml:"ffmpeg"`
	Nginx        Nginx   `yaml:"nginx"`
	ServerInfoId uint    `yaml:"-"`
}

type Server struct {
	Profile string `yaml:"Profile"`
	Name    string `yaml:"name"`
	Ip      string `yaml:"ip"`
	Port    string `yaml:"port"`
	Mode    string `yaml:"mode"`
}

type Log struct {
	Format     string `yaml:"format"`
	Dir        string `yaml:"dir"`
	MaxSize    int    `yaml:"maxSize"`
	MaxAge     int    `yaml:"maxAge"`
	MaxBackups int    `yaml:"maxBackups"`
	Compress   bool   `yaml:"compress"`
	Localtime  bool   `yaml:"localtime"`
	ShowLine   bool   `yaml:"showLine"`
}

type Mysql struct {
	Dsn             string        `yaml:"dsn"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
}

type Swagger struct {
	Title    string   `yaml:"title"`
	Desc     string   `yaml:"desc"`
	Host     string   `yaml:"host"`
	Url      string   `yaml:"url"`
	BasePath string   `yaml:"basePath"`
	Schemes  []string `yaml:"schemes"`
}

type Ffmpeg struct {
	LibPath string `yaml:"libPath"`
}

type Nginx struct {
	LibPath string `yaml:"libPath"`
}
