package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joeshaw/envdecode"
	"log"
	"msg-scheduler/api"
	"msg-scheduler/common/db"
	"msg-scheduler/common/messaging"
	"msg-scheduler/docs"
)

type Config struct {
	AppName    string
	Env        string `env:"ENV,default=local"`
	DBURL      string `env:"DB_URL,default=./db.sqlite"`
	Port       string `env:"PORT,default=3000"`
	MsgService string `env:"MESSAGING_SERVICE,required"`
	MsgKey     string `env:"MESSAGING_KEY,required"`
}

// @title           msg-scheduler API
// @version         1.0
// @description     This is a sample email scheduler.

// @contact.name   Robertson Lima
// @contact.url    http://robertsonlima.com
// @contact.email  email@robertsonlima.com

// @host      localhost:3000
// @BasePath  /api/v1
func main() {
	docs.SwaggerInfo.BasePath = "/api/v1"
	cfg := Config{AppName: "msg-scheduler"}

	/*using envdecode to avoid repetition, but the same can be easily
	achieved with multiple os.Getenv(key) and ordinary error handling */
	if err := envdecode.StrictDecode(&cfg); err != nil {
		log.Fatal(err)
	}

	engine := gin.Default()

	api.RegisterRoutes(engine, db.Init(cfg.DBURL, cfg.Env), messaging.Init(cfg.MsgService, cfg.MsgKey))

	engine.Run(cfg.Port)
}
