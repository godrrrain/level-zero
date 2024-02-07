package main

import (
	"context"
	"fmt"

	"zero/handler"
	"zero/model"
	"zero/storage"
	"zero/subscribe"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	stan "github.com/nats-io/stan.go"
)

func main() {
	cache := make(map[string]*model.OrderDb)

	postgresURL := "postgres://program:test@localhost:5432/wb"
	// postgresURL := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s",
	// 	"postgres", 5432, "program", "orders", "test")
	psqlDB, err := storage.NewPgStorage(context.Background(), postgresURL)
	if err != nil {
		fmt.Printf("Postgresql init: %s", err)
	} else {
		fmt.Println("Connected to PostreSQL")
	}
	defer psqlDB.Close()

	subscribe.CacheRestore(context.Background(), psqlDB, cache)

	clusterID := "test-cluster"
	clientID := "publisher-1"
	natsURL := "nats://127.0.0.1:4222"
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL))
	if err != nil {
		fmt.Println(err)
	}
	subscribe.ChannelSubscription(context.Background(), psqlDB, cache, sc)

	handler := handler.NewHandler(psqlDB, cache)

	router := gin.Default()

	router.Use(cors.Default())

	router.GET("/api/v1/orders/:uid/", handler.GetOrderByUid)

	router.Run(":8080")
}
