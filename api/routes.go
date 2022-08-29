package api

import (
	"github.com/gin-gonic/gin"
	"github.com/msg-scheduler/common/messaging"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func RegisterRoutes(engine *gin.Engine, db *gorm.DB, msgService messaging.MsgService) {
	h := &handler{
		DB:         db,
		msgService: msgService,
	}

	version1 := engine.Group("/api/v1/")
	{
		authGroup := version1.Group("/auth")
		{
			authGroup.POST("/login", BasicHandler(h.Login))
		}
		group := version1.Group("/users")
		{
			group.GET("/me", BasicHandler(h.GetMe))
			group.POST("/", BasicHandler(h.CreateUser))
			group.GET("/", BasicHandler(h.GetUsers))
			group.GET("/:id", BasicHandler(h.GetUser))
			group.DELETE("/:id", BasicHandler(h.DeleteUser))
		}

		msgGroup := version1.Group("/messages")
		{
			msgGroup.POST("/", BasicHandler(h.CreateMessage))
			msgGroup.GET("/", BasicHandler(h.GetMessages))
			msgGroup.GET("/:id", BasicHandler(h.GetMessage))
			msgGroup.POST("/:id", BasicHandler(h.UpdateMessage))
			msgGroup.DELETE("/:id", BasicHandler(h.DeleteMessage))
		}

		messagingGroup := version1.Group("/messaging")
		{
			messagingGroup.POST("/send", BasicHandler(h.SendTestMessage))
		}

		version1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

}
