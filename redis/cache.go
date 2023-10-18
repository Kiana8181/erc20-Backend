package cacheService

import (
	"context"
	"crypto/sha256"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func LoginTokenInsert(phoneNumber string, token string) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	input := token

	key := phoneNumber
	hashKey := sha256.New()
	hashKey.Write([]byte(key))
	hashKeyStr := string(hashKey.Sum(nil))

	// SET key value EX 10 NX
	_, err := rdb.SetNX(ctx, hashKeyStr, input, 180*time.Second).Result()
	if err != nil {
		return err
	}

	return nil
}

func LoginTokenFetch(phoneNumber string, otpCode string) (string, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	key := phoneNumber
	hashKey := sha256.New()
	hashKey.Write([]byte(key))
	hashKeyStr := string(hashKey.Sum(nil))

	raw_TokenCode, _ := rdb.Get(ctx, hashKeyStr).Result()
	credentials := strings.Split(raw_TokenCode, "/(+)/")

	if credentials[0] != otpCode {
		return "", errors.New("otp code mismatched")
	}
	token := credentials[1]

	return token, nil
}
