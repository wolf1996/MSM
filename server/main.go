package main

import (
	"encoding/json"
	"MSM/server/application"
	"MSM/server/logsystem"
	"io/ioutil"
	"os"
)

type config struct {
	Port    string
	DbLogin string
	DbPass  string
	DbUrl   string
}

func main() {
	logsystem.Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	file, err := os.Open("/home/lieroz/go/src/MSM/server/config.json")
	if err != nil {
		logsystem.Error.Printf("%s , %s", "config file open error", err)
		return
	}
	conf := config{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&conf)
	if err != nil {
		logsystem.Error.Printf("%s %s", "config parse error", err)
		return
	}
	logsystem.Info.Printf("Running on port : %s ", conf.Port)
	application.AppStart(conf.Port, conf.DbLogin, conf.DbPass, conf.DbUrl)
}
