package main

import (
	"context"
	"net/http"
	"time"

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
	router.GET("/api/alltemplates", getAllTemplatesHandler)
	router.PUT("/api/increaselocktemplate/:id", increaseLockTemplateHandler) // New API endpoint

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

	// Set expiration time to 1 minute
	expiration := time.Minute

	err = rdb.HSet(ctx, "locked_templates", id, id).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to lock template"})
		return
	}

	err = rdb.Expire(ctx, "locked_templates", expiration).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set expiration time"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "time": expiration.Seconds(), "msg": "template locked successfully"})
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

	// Retrieve expiration time of the key
	expiration, err := rdb.TTL(ctx, "locked_templates").Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get expiration time"})
		return
	}

	if expiration.Seconds() <= 0 {
		// If expiration time is negative or zero, delete the key-value pair
		err := rdb.HDel(ctx, "locked_templates", id).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to release lock"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": id, "msg": "template released"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "time": expiration.Seconds(), "msg": "template locked"})
}

func releaseLockTemplateHandler(c *gin.Context) {
	id := c.Param("id")

	err := rdb.HDel(ctx, "locked_templates", id).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to release lock"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "msg": "template unlocked"})
}

func getAllTemplatesHandler(c *gin.Context) {
	// Get all template keys
	keys, err := rdb.HKeys(ctx, "locked_templates").Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve templates"})
		return
	}

	// Initialize a map to store template IDs and their remaining expiration time
	templates := make(map[string]interface{})

	// Iterate over each key and get its remaining expiration time
	for _, key := range keys {
		expiration, err := rdb.TTL(ctx, key).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get expiration time"})
			return
		}
		// Convert expiration time to seconds without scientific notation
		expirationSeconds := int(expiration.Seconds())
		templates[key] = gin.H{"id": key, "time": expirationSeconds}
	}

	c.JSON(http.StatusOK, templates)
}

func increaseLockTemplateHandler(c *gin.Context) {
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

	// Get current expiration time
	currentExpiration, err := rdb.TTL(ctx, id).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get expiration time"})
		return
	}

	// Increase expiration time by one minute
	newExpiration := currentExpiration + time.Minute
	err = rdb.Expire(ctx, id, newExpiration).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to increase expiration time"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "new_time": newExpiration.Seconds(), "msg": "expiration time increased"})
}
