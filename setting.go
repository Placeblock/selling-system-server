package main

import (
	"log"

	"github.com/go-ini/ini"
)

type Database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
}

var DatabaseSetting = &Database{}

type App struct {
	SellPassword string
}

var AppSetting = &App{}

type Server struct {
	RunMode  string
	HttpPort int
}

var ServerSetting = &Server{}

var cfg *ini.File

func SetupSetting() {
	var err error
	cfg, err = ini.Load("config.ini")
	if err != nil {
		log.Fatalf("Failed to parse config.ini: %v", err)
	}

	mapTo("database", DatabaseSetting)
	mapTo("server", ServerSetting)
	mapTo("app", AppSetting)
}

func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}
