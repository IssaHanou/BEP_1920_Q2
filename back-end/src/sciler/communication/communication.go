package communication

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	logger "github.com/sirupsen/logrus"
	"time"
)

// Communicator is a type that maintains communication with the front-end and the client computers.
type Communicator struct {
	clientOptions    mqtt.ClientOptions
	client           mqtt.Client
	topicsOfInterest []string
}

// NewCommunicator is a constructor for a Communicator
func NewCommunicator(host string, port string, topicsOfInterest []string) *Communicator {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("%s://%s:%s", "tcp", host, port))
	opts.SetClientID("back-end")
	opts.SetConnectionLostHandler(onConnectionLost)
	client := mqtt.NewClient(opts)
	return &Communicator{*opts, client, topicsOfInterest}
}

// Start is a function that will start the communication by connecting to the broker and subscribing to all topics of interest
func (communicator *Communicator) Start() {
	action(communicator.client.Connect, "connect", -1)
	topics := make(map[string]byte)
	for _, topic := range communicator.topicsOfInterest {
		topics[topic] = byte(0)
	}
	action(func() mqtt.Token {
		return communicator.client.SubscribeMultiple(topics, onIncomingDataReceived)
	}, "subscribing", -1)
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

// Publish is a method that will send a message to a specific topic
func (communicator *Communicator) Publish(topic string, message string) {
	action(func() mqtt.Token {
		return communicator.client.Publish(topic, byte(0), false, message)
	}, "publish", 3)
}

// action is a function that will execute a communication action to the broker
// actionType is the description of the action in one word for logging
// retrials is the maximum number of times the action re-executed when failing, when retrials < 0, it is tried forever
func action(action func() mqtt.Token, actionType string, retrials int) {
	for i := 0; i < retrials || retrials < 0; i++ {
		token := action()
		if token.Wait() && token.Error() != nil {
			logger.Warnf("Fail to %s, %v", actionType, token.Error())
			time.Sleep(1 * time.Second)

			logger.Infof("Retry %d to %s", i +1, actionType)
			continue
		} else {
			logger.Infof("%s successful!", actionType)
			return
		}
	}
	logger.Errorf("All retries to %s failed, giving up :(", actionType)
}
