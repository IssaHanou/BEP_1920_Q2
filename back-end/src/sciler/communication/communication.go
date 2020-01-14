package communication

import (
	"encoding/json"
	"errors"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	logger "github.com/sirupsen/logrus"
	"time"
)

// Communicator is a type that maintains communication with the front-end and the client computers.
type Communicator struct {
	client mqtt.Client
}

// setClient is a setter for client
// this method is intended for testing only
func (communicator *Communicator) setClient(client mqtt.Client) {
	communicator.client = client
}

// Start is a function that will start the communication by connecting to the broker and subscribing to all topics of interest
func (communicator *Communicator) Start(host string, port int, topicsOfInterest []string, handler mqtt.MessageHandler, onStart func()) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("%s://%s:%d", "tcp", host, port))
	opts.SetClientID("back-end")
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		if err.Error() == "EOF" {
			logger.Warnf("connection lost: broker closed connection, (multiple client ID: %s?)", opts.ClientID)
		} else {
			logger.Warnf("connection lost: %v", err)
		}
	})
	opts.SetConnectRetry(false)
	opts.SetAutoReconnect(true)

	timeout := time.Duration(5) * time.Second
	opts.SetKeepAlive(timeout)            // time before sending a PING request to the broker
	opts.SetPingTimeout(timeout)          // time after sending a PING request to the broker
	opts.SetConnectTimeout(timeout)       // time before timing out and erroring the attempt
	opts.SetMaxReconnectInterval(timeout) // max time before retrying to reconnect

	opts.SetOnConnectHandler(func(client mqtt.Client) {
		// on connect subscribe and execute onStart
		topics := make(map[string]byte)
		for _, topic := range topicsOfInterest {
			topics[topic] = byte(0)
		}
		_ = action(func() mqtt.Token {
			return communicator.client.SubscribeMultiple(topics, handler)
		}, "subscribing", -1)

		logger.Infof("Connected to %s:%d with subscriptions to %s", host, port, topicsOfInterest)
		onStart()
	})
	opts.SetReconnectingHandler(func(client mqtt.Client, options *mqtt.ClientOptions) {
		logger.Info("Trying to reconnect")
	})
	will, _ := json.Marshal(map[string]interface{}{
		"device_id": "back-end",
		"time_sent": time.Now().Format("02-01-2006 15:04:05"),
		"type":      "status",
		"contents": map[string]interface{}{
			"id":         "front-end",
			"status":     map[string]interface{}{},
			"connection": false,
		},
	})
	opts.SetWill("front-end", string(will), 0, false)
	client := mqtt.NewClient(opts)

	communicator.client = client
	_ = action(communicator.client.Connect, "connect", -1)
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
			return err
		}
	}
	logger.Errorf("all retries to %s failed, giving up", actionType)
	return errors.New("action failed")
}
