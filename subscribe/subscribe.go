package subscribe

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"zero/model"
	"zero/storage"

	"github.com/nats-io/stan.go"
)

func ChannelSubscription(ctx context.Context, storage storage.Storage, cache map[string]*model.OrderDb, sc stan.Conn) {
	_, err := sc.Subscribe("foo", func(msg *stan.Msg) {
		err := insertingOrder(ctx, storage, cache, msg)
		if err != nil {
			log.Printf("err: %v", err)
			return
		}
	}, stan.DurableName("my-durable"))

	if err != nil {
		println(err)
	}
}

func insertingOrder(ctx context.Context, storage storage.Storage, cache map[string]*model.OrderDb, msg *stan.Msg) error {
	var order model.Order
	err := json.Unmarshal(msg.Data, &order)
	if err != nil {
		fmt.Printf("cannot unmarshal: %v", err)
	}
	dataJson, err := json.Marshal(order)
	if err != nil {
		fmt.Printf("cannot marshal: %v", err)
	}

	orderDb := model.OrderDb{
		OrderUID: order.OrderUID,
		DataJson: dataJson,
	}
	//запись в бд
	err = storage.CreateOrder(ctx, &orderDb)
	if err != nil {
		fmt.Printf("cannot create order in database: %v", err)
	}
	// запись в кеш
	cache[order.OrderUID] = &orderDb
	return nil
}

func CacheRestore(ctx context.Context, storage storage.Storage, cache map[string]*model.OrderDb) {
	ordersDb, err := storage.GetOrdersDb(ctx)
	if err != nil {
		fmt.Printf("cannot get orders from db: %v", err)
	}
	if len(ordersDb) == 0 {
		return
	}
	for _, orderDb := range ordersDb {
		cache[orderDb.OrderUID] = &orderDb
	}
}
