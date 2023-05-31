package pkg

import (
	"fmt"
	"math"
	"strconv"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

func TestXxx(t *testing.T) {
	RDB05 := redis.NewClient(&redis.Options{
		Addr:     "47.113.216.236:6379",
		Password: "147258", // no password set
		DB:       5,        // use default DB
	})
	now := time.Now()
	if RDB05 == nil {
		fmt.Println("------")
	}
	t2 := time.Date(2023, 2, 20, 0, 0, 0, 0, time.Local)
	i := t2.Unix()
	fmt.Println(i)
	fmt.Println("----------------------------:", RDB05)
	timestamp, err := RDB05.Get(RDB05.Context(), "福州大学").Result()
	fmt.Println("timestamp", timestamp)
	if err != nil {
		fmt.Println("err", err)
		return
	}

	fmt.Println("====1=====================")
	//当前学期开学时间
	timestamp2, _ := strconv.Atoi(timestamp)
	tbegin := time.Unix(int64(timestamp2), 0)
	dist := now.Sub(tbegin)
	//判断日期--在当前学期内
	fmt.Println("====2=====================")
	weeknow := math.Ceil(dist.Hours() / 24 / 7)
	fmt.Println("weeknow:", weeknow)
}
