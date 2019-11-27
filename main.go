package main

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/hughluo/go-tiny-url-api/controllers"
	"github.com/hughluo/go-tiny-url-api/models"
	"github.com/hughluo/go-tiny-url/pb"
	UTILS "github.com/hughluo/go-tiny-url/utils"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
)

func main() {
	UTILS.ConfigureLog()
	// Configure redis client
	redisClient := CreateClient()
	controllers.SetRedisClient(redisClient)

	// Configure gRPC client

	conn, err := grpc.Dial("kgs-service:50052", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	kgsClient := pb.NewKGSServiceClient(conn)
	models.SetKGSClient(kgsClient)

	// Configure REST API router
	router := httprouter.New()
	router.POST("/gotinyurl/", controllers.CreateTinyURL)
	router.GET("/gotinyurl/:tinyurl", controllers.RetrieveLongURL)

	log.Fatal(http.ListenAndServe(":8080", router))

}

func CreateClient() *redis.Client {
	REDIS_MAIN_PASSWORD := os.Getenv("REDIS_MAIN_PASSWORD")

	client := redis.NewClient(&redis.Options{
		Addr:     "redis-main-service:6379",
		Password: REDIS_MAIN_PASSWORD,
		DB:       0, // use default DB
	})

	pong, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(pong, err)
	// Output: PONG <nil>
	return client
}
