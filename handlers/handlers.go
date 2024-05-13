package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

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

func checkLockTemplateHandler(c *gin.Context, rdb *redis.Client) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	expiration, err := rdb.TTL(ctx, id).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get expiration time"})
		return
	}

	if expiration.Seconds() <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "template is not locked"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "time": expiration.Seconds(), "msg": "template is locked"})
}

func getAllTemplatesHandler(c *gin.Context, rdb *redis.Client) {
	keys, err := rdb.Keys(ctx, "*").Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve templates"})
		return
	}

	templates := make(map[string]string)

	for _, key := range keys {
		value, err := rdb.Get(ctx, key).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get value for key"})
			return
		}
		templates[key] = value
	}

	c.JSON(http.StatusOK, templates)
}

func lockTemplateHandler(c *gin.Context, rdb *redis.Client) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	_, err := rdb.Get(ctx, id).Result()
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"id": id, "error": "template already locked"})
		return
	}

	expiration := time.Minute * 15

	err = rdb.Set(ctx, id, id, expiration).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to lock template"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "msg": "template locked successfully"})
}

func releaseLockTemplateHandler(c *gin.Context, rdb *redis.Client) {
	id := c.Query("paramKey")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	err := rdb.Del(ctx, id).Err()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "failed to release lock"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "msg": "template unlocked"})
}

func increaseLockTemplateHandler(c *gin.Context, rdb *redis.Client) {
	id := c.Query("paramKey")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id not empty"})
		return
	}

	expiration := time.Minute * 15

	exists, err := rdb.Exists(ctx, id).Result()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "failed to check lock status"})
		return
	}

	if exists == 0 {
		c.JSON(http.StatusOK, gin.H{"error": "template not found"})
		return
	}

	err = rdb.Del(ctx, id).Err()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "failed to delete existing key"})
		return
	}

	err = rdb.Set(ctx, id, id, expiration).Err()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "failed to lock template"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "time": expiration.Seconds(), "msg": "template lock time increased successfully"})
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-ClusterQueueLenght, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
