package pkg

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

func GetWeek(rdb *redis.Client) float64 {
	now := time.Now()
	timestamp, err := rdb.Get(rdb.Context(), "福州大学").Result()
	fmt.Println("timestamp", timestamp)
	if err != nil {
		return -1
	}
	//当前学期开学时间
	timestamp2, _ := strconv.Atoi(timestamp)
	tbegin := time.Unix(int64(timestamp2), 0)
	dist := now.Sub(tbegin)
	//判断日期--在当前学期内
	weeknow := math.Ceil(dist.Hours() / 24 / 7)
	return weeknow

}
