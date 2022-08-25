package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"msg-scheduler/api"
	"msg-scheduler/common/db"
	"msg-scheduler/common/messaging"
	"msg-scheduler/docs"
)

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
