package models

import (
	"errors"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/hughluo/go-tiny-url/pb"
	"golang.org/x/net/context"
)

var KGS_CLIENT pb.KGSServiceClient

func SetKGSClient(client pb.KGSServiceClient) {
	KGS_CLIENT = client
}

func RetrieveLongURL(client *redis.Client, tinyURL string) (bool, string, string) {

	ok := existTinyURL(client, tinyURL)
	message := "Not Found"
	longURL := "NONE"

	if ok {
		longURL = getURLMapping(client, tinyURL)
		message = "Success"
	}
	return ok, message, longURL
}

func CreateTinyURL(client *redis.Client, longURL string, duration time.Duration) (bool, string, string) {
	got, tinyURL := getFreeTinyURL(client)
	if !got {
		panic(errors.New("not get free tinyurl"))
	}
	message := "Internal Error"
	ok := false
	if existTinyURL(client, tinyURL) {
		panic(errors.New("Free tinyURL from KGS already exists in DB"))
	} else {
		setURLMapping(client, tinyURL, longURL, duration)
		message = "Success"
		ok = true
	}
	return ok, message, tinyURL
}

func existTinyURL(client *redis.Client, tinyURL string) bool {
	exist, err := client.Exists(tinyURL).Result()
	if err != nil {
		panic(err)
	}
	return exist == 1
}

func getFreeTinyURL(client *redis.Client) (bool, string) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Minute)
	defer cancel()

	ok := false
	tinyURL := ""
	req := &pb.KGSRequest{Request: "REQUST"}
	if resp, err := KGS_CLIENT.GetFreeGoTinyURL(ctx, req); err == nil {
		ok = true
		tinyURL = resp.Result
	} else {

	}
	return ok, tinyURL
}

func setURLMapping(client *redis.Client, tinyURL string, longURL string, duration time.Duration) {
	err := client.Set(tinyURL, longURL, duration).Err()
	if err != nil {
		panic(err)
	}
}

func getURLMapping(client *redis.Client, tinyURL string) string {
	longURL, err := client.Get(tinyURL).Result()
	if err != nil {
		panic(err)
	}
	return longURL
}
