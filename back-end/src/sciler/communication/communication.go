package communication

import (
	"errors"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	logger "github.com/sirupsen/logrus"
	"time"
)

// Communicator is a type that maintains communication with the front-end and the client computers.
type Communicator struct {
	client           mqtt.Client
	topicsOfInterest []string
}

// NewCommunicator is a constructor for a Communicator
func NewCommunicator(host string, port int, topicsOfInterest []string) *Communicator {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("%s://%s:%d", "tcp", host, port))
	opts.SetClientID("back-end")
	opts.SetConnectionLostHandler(onConnectionLost)
	client := mqtt.NewClient(opts)
	return &Communicator{client, topicsOfInterest}
}

// Start is a function that will start the communication by connecting to the broker and subscribing to all topics of interest
func (communicator *Communicator) Start(handler mqtt.MessageHandler) {
	_ = action(communicator.client.Connect, "connect", -1)
	topics := make(map[string]byte)
	for _, topic := range communicator.topicsOfInterest {
		topics[topic] = byte(0)
	}
	_ = action(func() mqtt.Token {
		return communicator.client.SubscribeMultiple(topics, handler)
	}, "subscribing", -1)
}

func onConnectionLost(client mqtt.Client, e error) {
	logger.Warn(fmt.Sprintf("connection lost: %v", e))
	if client.IsConnected() {
		client.Disconnect(500)
	}
	// TODO reconnect
}

// Publish is a method that will send a message to a specific topic
func (communicator *Communicator) Publish(topic string, message string, retrials int) {
	err := action(func() mqtt.Token {
		return communicator.client.Publish(topic, byte(0), false, message)
	}, "publish", retrials)
	if err == errors.New("action failed") {
		// TODO reconnect
	}
}

// action is a function that will execute a communication action to the broker
// actionType is the description of the action in one word for logging
// retrials is the maximum number of times the action re-executed when failing, when retrials < 0, it is tried forever
func action(action func() mqtt.Token, actionType string, retrials int) error {
	for i := 0; i < retrials || retrials < 0; i++ {
		token := action()
		var err error = nil
		if token.Wait() && token.Error() != nil {
			logger.Warnf("fail to %s, %v", actionType, token.Error())
			time.Sleep(1 * time.Second)

			logger.Infof("retry %d to %s", i+1, actionType)
			err = errors.New("action eventually successful")
			continue
		} else {
			logger.Infof("back-end: %s successful!", actionType)
			return err
		}
	}
	logger.Errorf("all retries to %s failed, giving up", actionType)
	//TODO reconnect
	return errors.New("action failed")
}
