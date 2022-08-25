package api

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
	"msg-scheduler/common/messaging"
)

func RegisterRoutes(engine *gin.Engine, db *gorm.DB, msgService messaging.MsgService) {
	h := &handler{
		DB:         db,
		msgService: msgService,
	}

	version1 := engine.Group("/api/v1/")
	{
		group := version1.Group("/users")
		{
			group.POST("/", BasicHandler(h.CreateUser))
			group.GET("/", BasicHandler(h.GetUsers))
			group.GET("/:id", BasicHandler(h.GetUser))
			group.PATCH("/:id", BasicHandler(h.UpdateUser))
			group.DELETE("/:id", BasicHandler(h.DeleteUser))
		}

		msgGroup := version1.Group("/messages")
		{
			msgGroup.POST("/", BasicHandler(h.CreateMessage))
			msgGroup.GET("/", BasicHandler(h.UpdateMessage))
			msgGroup.GET("/:id", BasicHandler(h.GetMessage))
			msgGroup.PATCH("/:id", BasicHandler(h.UpdateMessage))
			msgGroup.DELETE("/:id", BasicHandler(h.DeleteMessage))
		}

		messagingGroup := version1.Group("/messaging")
		{
			messagingGroup.POST("/send", BasicHandler(h.SendTestMessage))
		}
		
		version1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

}