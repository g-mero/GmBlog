package config

import (
	"encoding/json"
	"fmt"
	"os"
)

var (
	AppMode  string
	HttpPort string
	Jwtkey   string

	DbFile string

	LocalAdmin bool
)

type config struct {
	AppMode  string `json:"app_mode"`
	HttpPort string `json:"port"`
	Jwtkey   string `json:"jwt_key"`

	DbFile     string `json:"db_file"`
	LocalAdmin bool   `json:"local_admin"`
}

var result config

func init() {
	var err error
	confFile, err := os.ReadFile("data/config.json")
	if err != nil {
		fmt.Println(err.Error())
		useDefault()
		return
	}
	err = json.Unmarshal(confFile, &result)
	if err != nil {
		fmt.Println(err.Error())
		useDefault()
		return
	}

	LoadServer()
	LoadSqlite()

	LocalAdmin = result.LocalAdmin

	fmt.Println(result)
}

func useDefault() {
	AppMode = "debug"
	HttpPort = ":3090"
	Jwtkey = "2h13hdsa"
	DbFile = "data/data.sqlite3"
	LocalAdmin = true
}

func LoadServer() {
	AppMode = result.AppMode
	HttpPort = result.HttpPort
	Jwtkey = result.HttpPort
}

func LoadSqlite() {
	DbFile = result.DbFile
}
