package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var rdb *redis.Client

type templateInfo struct {
	ID   string `json:"id"`
	Msg  string `json:"msg"`
	Time int    `json:"time"`
}

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

	template := templateInfo{
		ID:   id,
		Msg:  "template locked successfully",
		Time: int(expiration.Seconds()),
	}

	// Marshal templateInfo struct to JSON
	templateJSON, err := json.Marshal(template)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal template data"})
		return
	}

	err = rdb.HSet(ctx, "locked_templates", id, templateJSON).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to lock template"})
		return
	}

	err = rdb.Expire(ctx, "locked_templates", expiration).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set expiration time"})
		return
	}

	c.Data(http.StatusOK, "application/json", templateJSON)
}

func checkLockTemplateHandler(c *gin.Context) {
	id := c.Param("id")

	templateJSON, err := rdb.HGet(ctx, "locked_templates", id).Bytes()
	if err != nil {
		if err == redis.Nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "template not locked"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve template"})
		return
	}

	var template templateInfo
	err = json.Unmarshal(templateJSON, &template)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unmarshal template data"})
		return
	}

	c.Data(http.StatusOK, "application/json", templateJSON)
}

func releaseLockTemplateHandler(c *gin.Context) {
	id := c.Param("id")

	err := rdb.HDel(ctx, "locked_templates", id).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to release lock"})
		return
	}

	template := templateInfo{
		ID:  id,
		Msg: "template unlocked",
	}

	// Marshal templateInfo struct to JSON
	templateJSON, err := json.Marshal(template)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal template data"})
		return
	}

	c.Data(http.StatusOK, "application/json", templateJSON)
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
		templateJSON, err := rdb.HGet(ctx, "locked_templates", key).Bytes()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve template"})
			return
		}

		var template templateInfo
		err = json.Unmarshal(templateJSON, &template)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unmarshal template data"})
			return
		}

		templates[key] = template
	}

	c.JSON(http.StatusOK, templates)
}

func increaseLockTemplateHandler(c *gin.Context) {
	id := c.Param("id")

	// Check if the template exists
	templateJSON, err := rdb.HGet(ctx, "locked_templates", id).Bytes()
	if err != nil {
		if err == redis.Nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "template not locked"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve template"})
		return
	}

	var template templateInfo
	err = json.Unmarshal(templateJSON, &template)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unmarshal template data"})
		return
	}

	// Increase expiration time by one minute
	template.Time += 60

	// Marshal templateInfo struct to JSON
	templateJSON, err = json.Marshal(template)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal template data"})
		return
	}

	err = rdb.HSet(ctx, "locked_templates", id, templateJSON).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update expiration time"})
		return
	}

	c.Data(http.StatusOK, "application/json", templateJSON)
}
