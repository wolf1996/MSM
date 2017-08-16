package main

import (
	"encoding/json"
	"github.com/wolf1996/MSM/server/application"
	"github.com/wolf1996/MSM/server/logsystem"
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
	pwd, _ := os.Getwd()
	file, err := os.Open(pwd + "/server/config.json")
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
