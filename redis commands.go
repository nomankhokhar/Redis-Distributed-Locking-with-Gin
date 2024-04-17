package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func main() {
	// Create a Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",                // No password set
		DB:       0,                 // Use the default DB
	})

	// Example of SET command: Set the string value of a key
	err := client.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Example of GET command: Get the value of a key
	value, err := client.Get(ctx, "key").Result()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Value:", value) // Output: Value: value
	}

	// Example of INCR command: Increment the integer value of a key by one
	err = client.Incr(ctx, "counter").Err()
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Example of DECR command: Decrement the integer value of a key by one
	err = client.Decr(ctx, "counter").Err()
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Example of DEL command: Delete a key
	err = client.Del(ctx, "key").Err()
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Example of EXPIRE command: Set a key's time to live in seconds
	err = client.Expire(ctx, "key", time.Second).Err()
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Example of TTL command: Get the time to live for a key in seconds
	ttl, err := client.TTL(ctx, "key").Result()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("TTL:", ttl) // Output: TTL: 1s
	}

	// Example of HSET command: Set the string value of a hash field
	err = client.HSet(ctx, "hash", "field", "value").Err()
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Example of HGET command: Get the value of a hash field
	hashValue, err := client.HGet(ctx, "hash", "field").Result()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Hash Value:", hashValue) // Output: Hash Value: value
	}

	// Example of HDEL command: Delete one or more hash fields
	err = client.HDel(ctx, "hash", "field").Err()
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Example of HGETALL command: Get all the fields and values in a hash
	hashResult, err := client.HGetAll(ctx, "hash").Result()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Hash Result:", hashResult) // Output: Hash Result: map[]
	}

	// Example of LPUSH command: Insert one or more values at the head of the list
	err = client.LPush(ctx, "list", "value1", "value2").Err()
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Example of RPUSH command: Insert one or more values at the tail of the list
	err = client.RPush(ctx, "list", "value1", "value2").Err()
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Example of LPOP command: Remove and get the first element in a list
	listValue, err := client.LPop(ctx, "list").Result()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("List Value:", listValue) // Output: List Value: value2
	}

	// Example of RPOP command: Remove and get the last element in a list
	listValue, err = client.RPop(ctx, "list").Result()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("List Value:", listValue) // Output: List Value: value1
	}

	// Example of SADD command: Add one or more members to a set
	err = client.SAdd(ctx, "set", "member1", "member2").Err()
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Example of SMEMBERS command: Get all the members of a set
	setMembers, err := client.SMembers(ctx, "set").Result()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Set Members:", setMembers) // Output: Set Members: [member1 member2]
	}

	// Example of SREM command: Remove one or more members from a set
	err = client.SRem(ctx, "set", "member1", "member2").Err()
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Example of SISMEMBER command: Determine if a given value is a member of a set
	isMember, err := client.SIsMember(ctx, "set", "member").Result()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Is Member:", isMember) // Output: Is Member: false
	}

	// Example of ZADD command: Add one or more members to a sorted set
	err = client.ZAdd(ctx, "sorted_set", &redis.Z{Score: 1, Member: "member1"}, &redis.Z{Score: 2, Member: "member2"}).Err()
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Example of ZRANGE command: Return a range of members in a sorted set, by index
	sortedSetRange, err := client.ZRange(ctx, "sorted_set", 0, -1).Result()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Sorted Set Range:", sortedSetRange) // Output: Sorted Set Range: [member1 member2]
	}

	// Example of ZREM command: Remove one or more members from a sorted set
	err = client.ZRem(ctx, "sorted_set", "member1", "member2").Err()
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Example of ZSCORE command: Get the score associated with the given member in a sorted set
	score, err := client.ZScore(ctx, "sorted_set", "member").Result()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Score:", score) // Output: Score: <nil>
	}
}