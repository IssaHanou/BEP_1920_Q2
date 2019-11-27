package communication

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestNewCommunicator(t *testing.T) {
	type args struct {
		host             string
		port             string
		topicsOfInterest []string
	}

	optionsLocal := mqtt.NewClientOptions()
	optionsLocal.AddBroker(fmt.Sprintf("%s://%s:%s", "tcp", "localhost", "1883"))
	optionsLocal.SetClientID("back-end")
	optionsLocal.SetConnectionLostHandler(onConnectionLost)

	tests := []struct {
		name string
		args args
		want *Communicator
	}{
		{
			name: "two topics",
			args: args{
				host:             "localhost",
				port:             "1883",
				topicsOfInterest: []string{"test", "back-end"},
			},
			want: &Communicator{
				clientOptions:    *optionsLocal,
				client:           mqtt.NewClient(optionsLocal),
				topicsOfInterest: []string{"test", "back-end"},
			},
		},
		{
			name: "no topics",
			args: args{
				host:             "localhost",
				port:             "1883",
				topicsOfInterest: nil,
			},
			want: &Communicator{
				clientOptions:    *optionsLocal,
				client:           mqtt.NewClient(optionsLocal),
				topicsOfInterest: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCommunicator(tt.args.host, tt.args.port, tt.args.topicsOfInterest)
			assert.Equal(t, got.topicsOfInterest, tt.want.topicsOfInterest)
		})
	}
}

type TokenMock struct {
	mock.Mock
}

func (t TokenMock) Wait() bool {
	return false
}

func (t TokenMock) WaitTimeout(time.Duration) bool {
	return true
}

func (t TokenMock) Error() error {
	return nil
}

type ClientMock struct {
	mock.Mock
}

func (c ClientMock) IsConnected() bool {
	return true
}

func (c ClientMock) IsConnectionOpen() bool {
	return true
}

func (c ClientMock) Connect() mqtt.Token {
	return new(TokenMock)
}

func (c ClientMock) Disconnect(quiesce uint) {
	return
}

func (c ClientMock) Publish(topic string, qos byte, retained bool, payload interface{}) mqtt.Token {
	return new(TokenMock)
}

func (c ClientMock) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) mqtt.Token {
	return new(TokenMock)
}

func (c ClientMock) SubscribeMultiple(filters map[string]byte, callback mqtt.MessageHandler) mqtt.Token {
	return new(TokenMock)
}

func (c ClientMock) Unsubscribe(topics ...string) mqtt.Token {
	return new(TokenMock)
}

func (c ClientMock) AddRoute(topic string, callback mqtt.MessageHandler) {
	return
}

func (c ClientMock) OptionsReader() mqtt.ClientOptionsReader {
	return *new(mqtt.ClientOptionsReader)
}

func TestCommunicator_Start(t *testing.T) {
// Todo: Add test
}

func TestCommunicator_Publish(t *testing.T) {
// Todo: Add test
}
