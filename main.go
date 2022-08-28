package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"msg-scheduler/api"
	"msg-scheduler/common/db"
	"msg-scheduler/common/messaging"
	"msg-scheduler/docs"
)

// @title           msg-scheduler API
// @version         1.0
// @description     This is a sample email scheduler.

// @contact.name   Robertson Lima
// @contact.url    http://robertsonlima.com
// @contact.email  email@robertsonlima.com

// @host      localhost:3000
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth
func main() {
	docs.SwaggerInfo.BasePath = "/api/v1"
	viper.SetConfigFile("./dev.env")
	viper.ReadInConfig()

	env := viper.Get("ENV").(string)
	port := viper.Get("PORT").(string)
	dbUrl := viper.Get("DB_URL").(string)

	msgService := viper.Get("MESSAGING_SERVICE").(string)
	msgKey := viper.Get("MESSAGING_KEY").(string)

	engine := gin.Default()

	api.RegisterRoutes(engine, db.Init(dbUrl, env), messaging.Init(msgService, msgKey))

	engine.Run(port)
}
