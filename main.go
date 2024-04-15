package main

import (
	"context"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var rdb *redis.Client

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "redis-12200.c330.asia-south1-1.gce.cloud.redislabs.com:12200",
		Password: "td0g9FgW67e7BZx1RMx5UVNceFSvVkKa",
		DB:       0,
	})

	err := rdb.Ping(ctx).Err()
	if err != nil {
		panic("failed to connect to Redis: " + err.Error())
	}
}

func main() {
	router := gin.Default()

	router.POST("/api/locktemplate/:id", lockTemplateHandler)
	router.GET("/api/checklocktemplate/:id", checkLockTemplateHandler)
	router.DELETE("/api/releaselocktemplate/:id", releaseLockTemplateHandler)

	router.Run(":8080")
}

func lockTemplateHandler(c *gin.Context) {
	id := c.Param("id")

	exists, err := rdb.HExists(ctx, "locked_templates", id).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check lock status"})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "template already locked"})
		return
	}

	err = rdb.HSet(ctx, "locked_templates", id, id).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to lock template"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "template locked successfully", "id": id})
}

func checkLockTemplateHandler(c *gin.Context) {
	id := c.Param("id")

	_, err := rdb.HExists(ctx, "locked_templates", id).Result()
	if err != nil {
		if err == redis.Nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "template not locked"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check lock status"})
		return
	}

	val, err := rdb.HGet(ctx, "locked_templates", id).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get lock value"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "template locked", "id": id, "value": val})
}

func releaseLockTemplateHandler(c *gin.Context) {
	id := c.Param("id")

	err := rdb.HDel(ctx, "locked_templates", id).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to release lock"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "template lock released", "id": id})
}