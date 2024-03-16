/* trunk-ignore-all(gofmt) */
package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)
	fileContent, err := os.ReadFile("./orders_test.json")
	if err != nil {
		log.Fatalf("Ошибка чтения файла: %v", err)
	}
	var orders []map[string]interface{}
	if err := json.Unmarshal(fileContent, &orders); err != nil {
		log.Fatalf("Ошибка декодирования JSON: %v", err)
	}
	for _, order := range orders {
		orderJSON, err := json.Marshal(order)
		if err != nil {
			log.Fatalf("Ошибка кодирования JSON: %v", err)
		}
		var response string
		nc.PublishRequest("orders", response, orderJSON)
	}
	defer nc.Close()
}