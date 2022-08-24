package api

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func RegisterRoutes(engine *gin.Engine, db *gorm.DB) {
	h := &handler{
		DB: db,
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
		version1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

}
