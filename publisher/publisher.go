package main

import (
	//"io"
	//"os"
	"fmt"
	"log"
	"time"
	"math/rand"

	"wb_level_0/internal/order"
	"github.com/icrowley/fake"
	"github.com/nats-io/nats.go"
)

const orderQuantity = 10

func main() {
	// загружаю заказы из файла
	/*
	file, err := os.OpenFile("orders.json", os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	orders, err := order.UnmarshalOrders(fileData)
	if err != nil {
		log.Fatal(err)
	}
	*/
	nc, err := nats.Connect(nats.DefaultURL) // localhost:4222
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	for i := 0; i < orderQuantity; i++ {
		order := createFakeOrder()
		orderData, err := order.Marshal()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println()
		fmt.Println(string(orderData))
		fmt.Println()
		
		err = nc.Publish("order", orderData)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("publisher: sent a message about the order")
		
		time.Sleep(1 * time.Second)
	}
}

func createFakeOrder() *order.Order {
	odr := order.Order{
		OrderUID:          fake.DigitsN(21),
		TrackNumber:       fake.CharactersN(10),
		Entry:             fake.CharactersN(5),
		Locale:            "en",
		InternalSignature: fake.CharactersN(15),
		CustomerID:        fake.CharactersN(10),
		DeliveryService:   fake.Company(),
		Shardkey:          fake.CharactersN(10),
		SmID:              int64(rand.Intn(100)),
		DateCreated:       "2021-11-26T06:22:19Z",
		OofShard:          fake.CharactersN(5),
	}

	odr.Delivery = order.Delivery{
		Name:    fake.FullName(),
		Phone:   fake.Phone(),
		Zip:     fake.Zip(),
		City:    fake.City(),
		Address: fake.StreetAddress(),
		Region:  fake.State(),
		Email:   fake.EmailAddress(),
	}

	odr.Payment = order.Payment{
		Transaction:  fake.CharactersN(10),
		RequestID:    fake.CharactersN(8),
		Currency:     fake.CurrencyCode(),
		Provider:     fake.Company(),
		Amount:       3,
		PaymentDt:    time.Now().Unix(),
		Bank:         "SBER",
		DeliveryCost: int64(rand.Intn(10000)),
		GoodsTotal:  int64(rand.Intn(100)),
		CustomFee:  int64(rand.Intn(10)),
	}
	for i := 0; i < 3; i++ {
		item := order.Item{
			ChrtID:      int64(i + 1),
			TrackNumber: fake.CharactersN(8),
			Price:       int64(rand.Intn(30000)),
			Rid:         fake.CharactersN(5),
			Name:        fake.ProductName(),
			Sale:        int64(rand.Intn(30)),
			Size:        "0",
			TotalPrice:  int64(rand.Intn(1000)),
			NmID:        int64(rand.Intn(599)),
			Brand:       fake.Brand(),
			Status:      int64(rand.Intn(400)),
		}
		odr.Items = append(odr.Items, item)
	}
	return &odr
}
