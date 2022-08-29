package main

import (
	"github.com/joeshaw/envdecode"
	"github.com/r0bertson/msg-scheduler/common/db"
	"github.com/r0bertson/msg-scheduler/common/messaging"
	"github.com/r0bertson/msg-scheduler/common/utils"
	"github.com/r0bertson/msg-scheduler/scheduler/service"
	"log"
)

func main() {
	cfg := utils.Config{AppName: "github.com/r0bertson/msg-scheduler/scheduler"}

	/*using envdecode to avoid repetition, but the same can be easily
	achieved with multiple os.Getenv(key) and ordinary error handling */
	if err := envdecode.StrictDecode(&cfg); err != nil {
		log.Fatal(err)
	}
	serv := service.Init(db.Init(cfg.DBURL, cfg.Env), messaging.Init(cfg.MsgService, cfg.MsgKey))
	serv.RunCronJobs()
}
