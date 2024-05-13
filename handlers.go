package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func RegisterHandlers(router *gin.Engine, rdb *redis.Client) {
	router.POST("/locktemplate/:id", func(c *gin.Context) {
		lockTemplateHandler(c, rdb)
	})
	router.GET("/checklocktemplate/:id", func(c *gin.Context) {
		checkLockTemplateHandler(c, rdb)
	})
	router.POST("/delete/", func(c *gin.Context) {
		releaseLockTemplateHandler(c, rdb)
	})
	router.GET("/alltemplates", func(c *gin.Context) {
		getAllTemplatesHandler(c, rdb)
	})
	router.PUT("/increastime/", func(c *gin.Context) {
		increaseLockTemplateHandler(c, rdb)
	})
}
