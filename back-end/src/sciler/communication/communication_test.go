package communication

import (
	"errors"
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
		port             int
		topicsOfInterest []string
	}

	optionsLocal := mqtt.NewClientOptions()
	optionsLocal.AddBroker(fmt.Sprintf("%s://%s:%d", "tcp", "localhost", 1883))
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
				port:             1883,
				topicsOfInterest: []string{"test", "back-end"},
			},
			want: &Communicator{
				client:           mqtt.NewClient(optionsLocal),
				topicsOfInterest: []string{"test", "back-end"},
			},
		},
		{
			name: "no topics",
			args: args{
				host:             "localhost",
				port:             1883,
				topicsOfInterest: nil,
			},
			want: &Communicator{
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

type TokenMockSuccess struct {
	mock.Mock
}

func (t TokenMockSuccess) Wait() bool {
	return false
}

func (t TokenMockSuccess) WaitTimeout(time.Duration) bool {
	return true
}

func (t TokenMockSuccess) Error() error {
	return nil
}

type TokenMockFailure struct {
	mock.Mock
}

func (t TokenMockFailure) Wait() bool {
	return true
}

func (t TokenMockFailure) WaitTimeout(time.Duration) bool {
	return true
}

func (t TokenMockFailure) Error() error {
	return errors.New("testing error of TokenMockFailure")
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
	args := c.Called()
	return args.Get(0).(mqtt.Token)
}

func (c ClientMock) Disconnect(quiesce uint) {
	return
}

func (c ClientMock) Publish(topic string, qos byte, retained bool, payload interface{}) mqtt.Token {
	args := c.Called(topic, qos, retained, payload)
	return args.Get(0).(mqtt.Token)
}

func (c ClientMock) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) mqtt.Token {
	return new(TokenMockSuccess)
}

func (c ClientMock) SubscribeMultiple(filters map[string]byte, callback mqtt.MessageHandler) mqtt.Token {
	return new(TokenMockSuccess)
}

func (c ClientMock) Unsubscribe(topics ...string) mqtt.Token {
	return new(TokenMockSuccess)
}

func (c ClientMock) AddRoute(topic string, callback mqtt.MessageHandler) {
	return
}

func (c ClientMock) OptionsReader() mqtt.ClientOptionsReader {
	return *new(mqtt.ClientOptionsReader)
}

func TestCommunicator_Start(t *testing.T) {
	client := new(ClientMock)
	communicator := Communicator{
		client:           client,
		topicsOfInterest: []string{"back-end", "test"},
	}

	client.On("Connect").Return(new(TokenMockSuccess)).Once()
	communicator.Start(func(client mqtt.Client, message mqtt.Message) {}, func() {})
	client.AssertExpectations(t)
}

func TestCommunicator_Publish(t *testing.T) {
	client := new(ClientMock)
	communicator := Communicator{
		client:           client,
		topicsOfInterest: []string{"back-end", "test"},
	}
	client.On("Publish", "test", byte(0), false, "json").Return(new(TokenMockSuccess)).Once()
	communicator.Publish("test", "json", 3)
	client.AssertExpectations(t)
}

func TestCommunicator_PublishFailure(t *testing.T) {
	client := new(ClientMock)
	communicator := Communicator{
		client:           client,
		topicsOfInterest: []string{"back-end", "test"},
	}
	client.On("Publish", "test", byte(0), false, "json").Return(new(TokenMockFailure)).Times(3)
	communicator.Publish("test", "json", 3)
	client.AssertExpectations(t)
}
