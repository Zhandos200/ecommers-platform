package main

import (
	"consumer-service/internal/nats"
	"consumer-service/logger"
	"fmt"
)

func main() {
	logger.InitLogger()
	logger.Log.Info("🔄 Consumer service started")

	// Connect to NATS
	natsClient, err := nats.NewSubscriber("nats://nats:4222")
	if err != nil {
		logger.Log.Error(fmt.Sprintf("Failed to connect to NATS: %v", err))
		return
	}
	defer natsClient.Close()

	logger.Log.Info("✅ Consumer Service connected to NATS")

	// Subscribe to "order.created" topic
	if err := natsClient.Subscribe("order.created"); err != nil {
		logger.Log.Error(fmt.Sprintf("Failed to subscribe: %v", err))
		return
	}

	logger.Log.Info("🚀 Listening for events...")

	// Block forever
	select {}
}
