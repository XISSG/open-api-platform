package utils

import (
	"github.com/bwmarrin/snowflake"
	"log"
	"math/rand"
	"time"
)

func Snowflake() int64 {
	nodeID := int64(1)

	// 创建一个雪花节点
	node, err := snowflake.NewNode(nodeID)
	if err != nil {
		log.Fatalf("Failed to create snowflake node: %v", err)
	}

	// 生成一个唯一 ID
	id := node.Generate()
	return id.Int64()
}

func RandomExpireTime() time.Duration {
	rand.Seed(time.Now().UnixNano())
	minExpire := time.Minute
	maxExpire := time.Hour
	expire := minExpire + time.Duration(rand.Int63n(int64(maxExpire-minExpire)))
	return expire
}
