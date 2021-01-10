package communication

import (
	"errors"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

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

func (t TokenMockSuccess) Done() <-chan struct {}{
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

func (t TokenMockFailure) Done() <-chan struct {}{
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
	communicator := Communicator{client: client}
	client.On("Connect").Return(new(TokenMockSuccess)).Once()
	communicator.Start()
	client.AssertExpectations(t)
}

func TestCommunicator_Publish(t *testing.T) {
	client := new(ClientMock)
	communicator := Communicator{client: client}
	client.On("Publish", "test", byte(0), false, "json").Return(new(TokenMockSuccess)).Once()
	communicator.setClient(client)
	communicator.Publish("test", "json", 3)
	client.AssertExpectations(t)
}

func TestCommunicator_PublishFailure(t *testing.T) {
	client := new(ClientMock)
	communicator := Communicator{client: client}
	client.On("Publish", "test", byte(0), false, "json").Return(new(TokenMockFailure)).Times(3)
	communicator.setClient(client)
	communicator.Publish("test", "json", 3)
	client.AssertExpectations(t)
}
