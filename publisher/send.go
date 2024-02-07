package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"zero/model"

	"github.com/google/uuid"
	stan "github.com/nats-io/stan.go"
)

func generateOrder() model.Order {
	orderUid := uuid.New().String()
	fmt.Println(orderUid)
	order := model.Order{
		OrderUID:    orderUid,
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: model.Delivery{
			Name:   "Test Testov",
			Phone:  "+9720000000",
			Zip:    "2639809",
			City:   "Kiryat Mozkin",
			Adress: "Ploshad Mira 15",
			Region: "Kraiot",
			Email:  "test@gmail.com",
		},
		Payment: model.Payment{
			Transaction:  "b563feb7b2b84b6test",
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDT:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: []model.Item{
			{
				ChrtID:      9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
				Rid:         "ab4219087a764ae0btest",
				Name:        "Mascaras",
				Sale:        30,
				Size:        "5",
				TotalPrice:  317,
				NmID:        2389212,
				Brand:       "VivienneSabo",
				Status:      202,
			},
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test",
		DeliveryService:   "meest",
		Shardkey:          "9",
		SmID:              99,
		DateCreated:       time.Now().Format("2006-01-02T15:04:05Z"),
		OofShard:          "1",
	}
	return order
}

func queueOrders(lenQueue, delaySec int, channel string, sc stan.Conn) {
	for i := 0; i < lenQueue; i++ {
		order := generateOrder()
		data, _ := json.Marshal(order)
		// Sending a message to the channel
		if err := sc.Publish(channel, data); err != nil {
			log.Fatalf("Error sending the message: %v", err)
		}

		time.Sleep(time.Second * time.Duration(delaySec))
	}

}

func main() {

	clusterID := "test-cluster"
	clientID := "consumer-1"
	natsURL := "nats://127.0.0.1:4222"

	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL))
	if err != nil {
		println(err)
	}

	queueOrders(1, 0, "foo", sc)

	// // Simple Async Subscriber
	// sub, _ := sc.Subscribe("foo", func(m *stan.Msg) {
	// 	fmt.Printf("Received a message: %s\n", string(m.Data))
	// })

	// Unsubscribe
	// sub.Unsubscribe()

	// // Close connection
	// sc.Close()

}
