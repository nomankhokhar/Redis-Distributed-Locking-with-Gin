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
	router.PUT("/api/increaselocktemplate/:id", increaseLockTemplateHandler)

	router.Run(":8080")
}



func checkLockTemplateHandler(c *gin.Context) {
    id := c.Param("id")

    // Get expiration time of the locked template
    expiration, err := rdb.TTL(ctx, id).Result()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get expiration time"})
        return
    }

    if expiration.Seconds() <= 0 {
        // If expiration time is negative or zero, the template is not locked
        c.JSON(http.StatusNotFound, gin.H{"error": "template is not locked"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"id": id, "time": expiration.Seconds(), "msg": "template locked"})
}

func getAllTemplatesHandler(c *gin.Context) {
	// Get all template keys
	keys, err := rdb.Keys(ctx, "*").Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve templates"})
		return
	}

	// Initialize a map to store template key-value pairs
	templates := make(map[string]string)

	// Iterate over each key and get its value
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

func lockTemplateHandler(c *gin.Context) {
    id := c.Param("id")

    // Check if the key already exists
    _, err := rdb.Get(ctx, id).Result()
    if err == nil {
        c.JSON(http.StatusConflict, gin.H{"error": "template already locked"})
        return
    }

	expiration := time.Minute

    // Store the ID as both key and value
    err = rdb.Set(ctx, id, id, expiration).Err() // 0 means no expiration
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to lock template"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"id": id, "msg": "template locked successfully"})
}

func releaseLockTemplateHandler(c *gin.Context) {
	id := c.Param("id")

	err := rdb.Del(ctx, id).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to release lock"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "msg": "template unlocked"})
}

func increaseLockTemplateHandler(c *gin.Context) {
    id := c.Param("id")

    // Set expiration time to 1 minute
    expiration := time.Minute

    // Check if the template with the provided ID exists
    exists, err := rdb.Exists(ctx, id).Result()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check lock status"})
        return
    }

    if exists == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
        return
    }

    // Delete the existing key-value pair
    err = rdb.Del(ctx, id).Err()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete existing key"})
        return
    }

    // Store a new key-value pair with the provided ID and set the expiration time
    err = rdb.Set(ctx, id, id, expiration).Err()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to lock template"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"id": id, "time": expiration.Seconds(), "msg": "template lock time increased successfully"})
}
