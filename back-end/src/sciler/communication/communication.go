package communication

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	logger "github.com/sirupsen/logrus"
	"time"
)

type communicator struct {
	clientOptions    mqtt.ClientOptions
	topicsOfInterest []string
}

// NewCommunicator is a constructor for a communicator
func NewCommunicator(host string, port string, topicsOfInterest []string) *communicator {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("%s://%s:%s", "tcp", host, port))
	opts.SetClientID("back-end")
	opts.SetConnectionLostHandler(onConnectionLost)
	return &communicator{*opts, topicsOfInterest}
}

func (communicator *communicator) Start() {
	client := mqtt.NewClient(&communicator.clientOptions)
	client.Connect()
	topics := make(map[string]byte)
	for _, topic := range communicator.topicsOfInterest {
		topics[topic] = byte(0)
	}
	for {
		token := client.SubscribeMultiple(topics, onIncomingDataReceived)
		if token.Wait() && token.Error() != nil {
			logger.Error("Fail to sub... ", token.Error())
			time.Sleep(1 * time.Second)

			logger.Errorf("Retry to subscribe")
			continue
		} else {
			logger.Info("Subscribe successful!")
			break
		}
	}
}

func onConnectionLost(client mqtt.Client, e error) {
	logger.Warn(fmt.Sprintf("Connection lost : %v", e))
	if client.IsConnected() {
		client.Disconnect(500)
	}

	// Todo try to reconnect
}

func onIncomingDataReceived(client mqtt.Client, message mqtt.Message) {
	logger.Info(string(message.Payload()))
}
