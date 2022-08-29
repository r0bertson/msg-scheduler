package main

import (
	"github.com/joeshaw/envdecode"
	"github.com/r0bertson/msg-scheduler/common/db"
	"github.com/r0bertson/msg-scheduler/scheduler/service"
	"log"
)

type Config struct {
	AppName string
	Env     string `env:"ENV,default=local"`
	DBURL   string `env:"DB_URL,default=./db.sqlite"`
}

func main() {
	cfg := Config{AppName: "github.com/r0bertson/msg-scheduler/scheduler"}

	/*using envdecode to avoid repetition, but the same can be easily
	achieved with multiple os.Getenv(key) and ordinary error handling */
	if err := envdecode.StrictDecode(&cfg); err != nil {
		log.Fatal(err)
	}
	serv := service.Init(db.Init(cfg.DBURL, cfg.Env))
	serv.RunCronJobs()
}
