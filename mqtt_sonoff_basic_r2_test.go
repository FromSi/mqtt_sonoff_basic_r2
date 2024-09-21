package mqtt_sonoff_basic_r2

import (
	"fmt"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

const MockCtxCmndResponseTimeoutInSeconds = 1

type MockMQTTServer struct {
	mock.Mock
	subscribeChan chan mqtt.InlineSubFn
}

func NewMockMQTTServer() (*SonoffBasicR2, *MockMQTTServer, error) {
	mockServer := new(MockMQTTServer)
	mockServer.subscribeChan = make(chan mqtt.InlineSubFn, 1)
	sonoffServer, err := NewSonoffBasicR2WithServer(mockServer, 1)

	if err != nil {
		return nil, nil, err
	}

	sonoffServer.SetCtxCmndResponseTimeoutInSeconds(MockCtxCmndResponseTimeoutInSeconds)

	mockServer.On("Subscribe", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockServer.On("Unsubscribe", mock.Anything, mock.Anything).Return(nil)
	mockServer.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	return sonoffServer, mockServer, sonoffServer.Serve()
}

func (m *MockMQTTServer) Serve() error {
	return m.Called().Error(0)
}

func (m *MockMQTTServer) Close() error {
	return m.Called().Error(0)
}

func (m *MockMQTTServer) Subscribe(filter string, subscriptionId int, handler mqtt.InlineSubFn) error {
	fullTeleTopicWithDefaultFormat := fmt.Sprintf("%s/%s/%s", TasmotaPrefixTele, TasmotaTeleTopicLWTValueAll, TasmotaTeleTopicLWT)

	if filter != fullTeleTopicWithDefaultFormat {
		m.subscribeChan <- handler
	}

	return m.Called(filter, subscriptionId, handler).Error(0)
}

func (m *MockMQTTServer) Unsubscribe(filter string, subscriptionId int) error {
	return m.Called(filter, subscriptionId).Error(0)
}

func (m *MockMQTTServer) Publish(topic string, payload []byte, retain bool, qos byte) error {
	return m.Called(topic, payload, retain, qos).Error(0)
}

func TestSonoffBasicR2_Close(t *testing.T) {
	sonoffServer, _, err := NewMockMQTTServer()

	assert.NoError(t, err)

	err = sonoffServer.Close()

	assert.NoError(t, err)

	connectedChan := sonoffServer.TeleConnected()
	disconnectedChan := sonoffServer.TeleDisconnected()

	select {
	case _, ok := <-connectedChan:
		assert.Equal(t, false, ok)
	default:
		t.Fatal("connectedChan is not closed")
	}

	select {
	case _, ok := <-disconnectedChan:
		assert.Equal(t, false, ok)
	default:
		t.Fatal("disconnectedChan is not closed")
	}
}

func TestSonoffBasicR2_CtxCmndResponseTimeoutInSeconds(t *testing.T) {
	sonoffServer, _, err := NewMockMQTTServer()

	assert.NoError(t, err)

	assert.Equal(t, MockCtxCmndResponseTimeoutInSeconds, int(sonoffServer.GetCtxCmndResponseTimeoutInSeconds()))

	sonoffServer.SetCtxCmndResponseTimeoutInSeconds(DefaultCtxCmndResponseTimeoutInSeconds)

	assert.Equal(t, DefaultCtxCmndResponseTimeoutInSeconds, int(sonoffServer.GetCtxCmndResponseTimeoutInSeconds()))

	err = sonoffServer.Close()

	assert.NoError(t, err)
}

func TestSonoffBasicR2_TeleConnected(t *testing.T) {
	sonoffServer, mockServer, err := NewMockMQTTServer()

	assert.NoError(t, err)

	connectedChan := sonoffServer.TeleConnected()

	fullTeleTopicOne := sonoffServer.getFullTeleTopic("1", TasmotaTeleTopicLWT)
	fullTeleTopicTwo := sonoffServer.getFullTeleTopic("2", TasmotaTeleTopicLWT)

	handler := mockServer.Calls[0].Arguments.Get(2).(mqtt.InlineSubFn)
	handler(nil, packets.Subscription{}, packets.Packet{TopicName: fullTeleTopicTwo, Payload: []byte(TasmotaTeleTopicLWTResponseOffline)})
	handler(nil, packets.Subscription{}, packets.Packet{TopicName: fullTeleTopicOne, Payload: []byte(TasmotaTeleTopicLWTResponseOnline)})

	assert.Equal(t, "1", <-connectedChan)

	err = sonoffServer.Close()

	assert.NoError(t, err)
}

func TestSonoffBasicR2_TeleDisconnected(t *testing.T) {
	sonoffServer, mockServer, err := NewMockMQTTServer()

	assert.NoError(t, err)

	disconnectedChan := sonoffServer.TeleDisconnected()

	fullTeleTopicOne := sonoffServer.getFullTeleTopic("1", TasmotaTeleTopicLWT)
	fullTeleTopicTwo := sonoffServer.getFullTeleTopic("2", TasmotaTeleTopicLWT)

	handler := mockServer.Calls[1].Arguments.Get(2).(mqtt.InlineSubFn)
	handler(nil, packets.Subscription{}, packets.Packet{TopicName: fullTeleTopicTwo, Payload: []byte(TasmotaTeleTopicLWTResponseOnline)})
	handler(nil, packets.Subscription{}, packets.Packet{TopicName: fullTeleTopicOne, Payload: []byte(TasmotaTeleTopicLWTResponseOffline)})

	assert.Equal(t, "1", <-disconnectedChan)

	err = sonoffServer.Close()

	assert.NoError(t, err)
}

func TestSonoffBasicR2_Status(t *testing.T) {
	sonoffServer, mockServer, err := NewMockMQTTServer()

	assert.NoError(t, err)

	responseChan := make(chan bool, 1)

	go func() {
		_, _ = sonoffServer.Status("1")

		responseChan <- true
	}()

	fullStatTopic := sonoffServer.getFullStatTopic("1", TasmotaStatTopicStatus)
	fullCmndTopic := sonoffServer.getFullCmndTopic("1", TasmotaCmndTopicStatusAll)

	handler := <-mockServer.subscribeChan
	handler(nil, packets.Subscription{}, packets.Packet{TopicName: fullStatTopic, Payload: []byte("test")})

	<-responseChan

	assert.Equal(t, fullCmndTopic, mockServer.Calls[3].Arguments.Get(0).(string))

	err = sonoffServer.Close()

	assert.NoError(t, err)
}

func TestSonoffBasicR2_StatusOne(t *testing.T) {
	sonoffServer, mockServer, err := NewMockMQTTServer()

	assert.NoError(t, err)

	responseChan := make(chan bool, 1)

	go func() {
		_, _ = sonoffServer.StatusOne("1")

		responseChan <- true
	}()

	fullStatTopic := sonoffServer.getFullStatTopic("1", TasmotaStatTopicStatusOne)
	fullCmndTopic := sonoffServer.getFullCmndTopic("1", TasmotaCmndTopicStatus)

	handler := <-mockServer.subscribeChan
	handler(nil, packets.Subscription{}, packets.Packet{TopicName: fullStatTopic, Payload: []byte("test")})

	<-responseChan

	assert.Equal(t, fullCmndTopic, mockServer.Calls[3].Arguments.Get(0).(string))

	err = sonoffServer.Close()

	assert.NoError(t, err)
}

func TestSonoffBasicR2_StatusTwo(t *testing.T) {
	sonoffServer, mockServer, err := NewMockMQTTServer()

	assert.NoError(t, err)

	responseChan := make(chan bool, 1)

	go func() {
		_, _ = sonoffServer.StatusTwo("1")

		responseChan <- true
	}()

	fullStatTopic := sonoffServer.getFullStatTopic("1", TasmotaStatTopicStatusTwo)
	fullCmndTopic := sonoffServer.getFullCmndTopic("1", TasmotaCmndTopicStatus)

	handler := <-mockServer.subscribeChan
	handler(nil, packets.Subscription{}, packets.Packet{TopicName: fullStatTopic, Payload: []byte("test")})

	<-responseChan

	assert.Equal(t, fullCmndTopic, mockServer.Calls[3].Arguments.Get(0).(string))

	err = sonoffServer.Close()

	assert.NoError(t, err)
}

func TestSonoffBasicR2_StatusThree(t *testing.T) {
	sonoffServer, mockServer, err := NewMockMQTTServer()

	assert.NoError(t, err)

	responseChan := make(chan bool, 1)

	go func() {
		_, _ = sonoffServer.StatusThree("1")

		responseChan <- true
	}()

	fullStatTopic := sonoffServer.getFullStatTopic("1", TasmotaStatTopicStatusThree)
	fullCmndTopic := sonoffServer.getFullCmndTopic("1", TasmotaCmndTopicStatus)

	handler := <-mockServer.subscribeChan
	handler(nil, packets.Subscription{}, packets.Packet{TopicName: fullStatTopic, Payload: []byte("test")})

	<-responseChan

	assert.Equal(t, fullCmndTopic, mockServer.Calls[3].Arguments.Get(0).(string))

	err = sonoffServer.Close()

	assert.NoError(t, err)
}

func TestSonoffBasicR2_StatusFour(t *testing.T) {
	sonoffServer, mockServer, err := NewMockMQTTServer()

	assert.NoError(t, err)

	responseChan := make(chan bool, 1)

	go func() {
		_, _ = sonoffServer.StatusFour("1")

		responseChan <- true
	}()

	fullStatTopic := sonoffServer.getFullStatTopic("1", TasmotaStatTopicStatusFour)
	fullCmndTopic := sonoffServer.getFullCmndTopic("1", TasmotaCmndTopicStatus)

	handler := <-mockServer.subscribeChan
	handler(nil, packets.Subscription{}, packets.Packet{TopicName: fullStatTopic, Payload: []byte("test")})

	<-responseChan

	assert.Equal(t, fullCmndTopic, mockServer.Calls[3].Arguments.Get(0).(string))

	err = sonoffServer.Close()

	assert.NoError(t, err)
}

func TestSonoffBasicR2_StatusFive(t *testing.T) {
	sonoffServer, mockServer, err := NewMockMQTTServer()

	assert.NoError(t, err)

	responseChan := make(chan bool, 1)

	go func() {
		_, _ = sonoffServer.StatusFive("1")

		responseChan <- true
	}()

	fullStatTopic := sonoffServer.getFullStatTopic("1", TasmotaStatTopicStatusFive)
	fullCmndTopic := sonoffServer.getFullCmndTopic("1", TasmotaCmndTopicStatus)

	handler := <-mockServer.subscribeChan
	handler(nil, packets.Subscription{}, packets.Packet{TopicName: fullStatTopic, Payload: []byte("test")})

	<-responseChan

	assert.Equal(t, fullCmndTopic, mockServer.Calls[3].Arguments.Get(0).(string))

	err = sonoffServer.Close()

	assert.NoError(t, err)
}

func TestSonoffBasicR2_StatusSix(t *testing.T) {
	sonoffServer, mockServer, err := NewMockMQTTServer()

	assert.NoError(t, err)

	responseChan := make(chan bool, 1)

	go func() {
		_, _ = sonoffServer.StatusSix("1")

		responseChan <- true
	}()

	fullStatTopic := sonoffServer.getFullStatTopic("1", TasmotaStatTopicStatusSix)
	fullCmndTopic := sonoffServer.getFullCmndTopic("1", TasmotaCmndTopicStatus)

	handler := <-mockServer.subscribeChan
	handler(nil, packets.Subscription{}, packets.Packet{TopicName: fullStatTopic, Payload: []byte("test")})

	<-responseChan

	assert.Equal(t, fullCmndTopic, mockServer.Calls[3].Arguments.Get(0).(string))

	err = sonoffServer.Close()

	assert.NoError(t, err)
}

func TestSonoffBasicR2_StatusSeven(t *testing.T) {
	sonoffServer, mockServer, err := NewMockMQTTServer()

	assert.NoError(t, err)

	responseChan := make(chan bool, 1)

	go func() {
		_, _ = sonoffServer.StatusSeven("1")

		responseChan <- true
	}()

	fullStatTopic := sonoffServer.getFullStatTopic("1", TasmotaStatTopicStatusSeven)
	fullCmndTopic := sonoffServer.getFullCmndTopic("1", TasmotaCmndTopicStatus)

	handler := <-mockServer.subscribeChan
	handler(nil, packets.Subscription{}, packets.Packet{TopicName: fullStatTopic, Payload: []byte("test")})

	<-responseChan

	assert.Equal(t, fullCmndTopic, mockServer.Calls[3].Arguments.Get(0).(string))

	err = sonoffServer.Close()

	assert.NoError(t, err)
}

func TestSonoffBasicR2_StatusEight(t *testing.T) {
	sonoffServer, mockServer, err := NewMockMQTTServer()

	assert.NoError(t, err)

	responseChan := make(chan bool, 1)

	go func() {
		_, _ = sonoffServer.StatusEight("1")

		responseChan <- true
	}()

	fullStatTopic := sonoffServer.getFullStatTopic("1", TasmotaStatTopicStatusEight)
	fullCmndTopic := sonoffServer.getFullCmndTopic("1", TasmotaCmndTopicStatus)

	handler := <-mockServer.subscribeChan
	handler(nil, packets.Subscription{}, packets.Packet{TopicName: fullStatTopic, Payload: []byte("test")})

	<-responseChan

	assert.Equal(t, fullCmndTopic, mockServer.Calls[3].Arguments.Get(0).(string))

	err = sonoffServer.Close()

	assert.NoError(t, err)
}

func TestSonoffBasicR2_StatusEleven(t *testing.T) {
	sonoffServer, mockServer, err := NewMockMQTTServer()

	assert.NoError(t, err)

	responseChan := make(chan bool, 1)

	go func() {
		_, _ = sonoffServer.StatusEleven("1")

		responseChan <- true
	}()

	fullStatTopic := sonoffServer.getFullStatTopic("1", TasmotaStatTopicStatusEleven)
	fullCmndTopic := sonoffServer.getFullCmndTopic("1", TasmotaCmndTopicStatus)

	handler := <-mockServer.subscribeChan
	handler(nil, packets.Subscription{}, packets.Packet{TopicName: fullStatTopic, Payload: []byte("test")})

	<-responseChan

	assert.Equal(t, fullCmndTopic, mockServer.Calls[3].Arguments.Get(0).(string))

	err = sonoffServer.Close()

	assert.NoError(t, err)
}

func TestSonoffBasicR2_StatusPhysicalButton(t *testing.T) {
	sonoffServer, mockServer, err := NewMockMQTTServer()

	assert.NoError(t, err)

	responseChan := make(chan bool, 1)

	go func() {
		_, _ = sonoffServer.StatusPhysicalButton("1")

		responseChan <- true
	}()

	fullStatTopic := sonoffServer.getFullStatTopic("1", TasmotaStatTopicResult)
	fullCmndTopic := sonoffServer.getFullCmndTopic("1", TasmotaCmndTopicPhysicalButton)

	handler := <-mockServer.subscribeChan
	handler(nil, packets.Subscription{}, packets.Packet{TopicName: fullStatTopic, Payload: []byte("test")})

	<-responseChan

	assert.Equal(t, fullCmndTopic, mockServer.Calls[3].Arguments.Get(0).(string))

	err = sonoffServer.Close()

	assert.NoError(t, err)
}

func TestSonoffBasicR2_PhysicalButtonOn(t *testing.T) {
	sonoffServer, mockServer, err := NewMockMQTTServer()

	assert.NoError(t, err)

	sonoffServer.PhysicalButtonOn("1")

	fullCmndTopic := sonoffServer.getFullCmndTopic("1", TasmotaCmndTopicPhysicalButton)

	assert.Equal(t, fullCmndTopic, mockServer.Calls[2].Arguments.Get(0).(string))
	assert.Equal(t, []byte(TasmotaCmndTopicPhysicalButtonValueOn), mockServer.Calls[2].Arguments.Get(1).([]byte))

	err = sonoffServer.Close()

	assert.NoError(t, err)
}

func TestSonoffBasicR2_PhysicalButtonOff(t *testing.T) {
	sonoffServer, mockServer, err := NewMockMQTTServer()

	assert.NoError(t, err)

	sonoffServer.PhysicalButtonOff("1")

	fullCmndTopic := sonoffServer.getFullCmndTopic("1", TasmotaCmndTopicPhysicalButton)

	assert.Equal(t, fullCmndTopic, mockServer.Calls[2].Arguments.Get(0).(string))
	assert.Equal(t, []byte(TasmotaCmndTopicPhysicalButtonValueOff), mockServer.Calls[2].Arguments.Get(1).([]byte))

	err = sonoffServer.Close()

	assert.NoError(t, err)
}

func TestSonoffBasicR2_PowerOn(t *testing.T) {
	sonoffServer, mockServer, err := NewMockMQTTServer()

	assert.NoError(t, err)

	sonoffServer.PowerOn("1")

	fullCmndTopic := sonoffServer.getFullCmndTopic("1", TasmotaCmndTopicPower)

	assert.Equal(t, fullCmndTopic, mockServer.Calls[2].Arguments.Get(0).(string))
	assert.Equal(t, []byte(TasmotaCmndTopicPowerValueOn), mockServer.Calls[2].Arguments.Get(1).([]byte))

	err = sonoffServer.Close()

	assert.NoError(t, err)
}

func TestSonoffBasicR2_PowerOff(t *testing.T) {
	sonoffServer, mockServer, err := NewMockMQTTServer()

	assert.NoError(t, err)

	sonoffServer.PowerOff("1")

	fullCmndTopic := sonoffServer.getFullCmndTopic("1", TasmotaCmndTopicPower)

	assert.Equal(t, fullCmndTopic, mockServer.Calls[2].Arguments.Get(0).(string))
	assert.Equal(t, []byte(TasmotaCmndTopicPowerValueOff), mockServer.Calls[2].Arguments.Get(1).([]byte))

	err = sonoffServer.Close()

	assert.NoError(t, err)
}

func TestSonoffBasicR2_PowerToggle(t *testing.T) {
	sonoffServer, mockServer, err := NewMockMQTTServer()

	assert.NoError(t, err)

	sonoffServer.PowerToggle("1")

	fullCmndTopic := sonoffServer.getFullCmndTopic("1", TasmotaCmndTopicPower)

	assert.Equal(t, fullCmndTopic, mockServer.Calls[2].Arguments.Get(0).(string))
	assert.Equal(t, []byte(TasmotaCmndTopicPowerValueToggle), mockServer.Calls[2].Arguments.Get(1).([]byte))

	err = sonoffServer.Close()

	assert.NoError(t, err)
}

func TestSonoffBasicR2_getFullTopic(t *testing.T) {
	sonoffServer, _, err := NewMockMQTTServer()

	assert.NoError(t, err)

	fullTeleTopicWithDefaultFormat := fmt.Sprintf("%s/%s/%s", TasmotaPrefixTele, "1", TasmotaTeleTopicLWT)
	fullStatTopicWithDefaultFormat := fmt.Sprintf("%s/%s/%s", TasmotaPrefixStat, "2", TasmotaStatTopicResult)
	fullCmndTopicWithDefaultFormat := fmt.Sprintf("%s/%s/%s", TasmotaPrefixCmnd, "3", TasmotaCmndTopicPower)

	assert.Equal(t, fullTeleTopicWithDefaultFormat, sonoffServer.getFullTopic(TasmotaPrefixTele, "1", TasmotaTeleTopicLWT))
	assert.Equal(t, fullStatTopicWithDefaultFormat, sonoffServer.getFullTopic(TasmotaPrefixStat, "2", TasmotaStatTopicResult))
	assert.Equal(t, fullCmndTopicWithDefaultFormat, sonoffServer.getFullTopic(TasmotaPrefixCmnd, "3", TasmotaCmndTopicPower))

	err = sonoffServer.Close()

	assert.NoError(t, err)
}

func TestSonoffBasicR2_getFullStatTopic(t *testing.T) {
	sonoffServer, _, err := NewMockMQTTServer()

	assert.NoError(t, err)

	fullStatTopicWithDefaultFormatOne := fmt.Sprintf("%s/%s/%s", TasmotaPrefixStat, "1", TasmotaStatTopicResult)
	fullStatTopicWithDefaultFormatTwo := fmt.Sprintf("%s/%s/%s", TasmotaPrefixStat, "2", TasmotaStatTopicStatus)

	assert.Equal(t, fullStatTopicWithDefaultFormatOne, sonoffServer.getFullStatTopic("1", TasmotaStatTopicResult))
	assert.Equal(t, fullStatTopicWithDefaultFormatTwo, sonoffServer.getFullStatTopic("2", TasmotaStatTopicStatus))

	err = sonoffServer.Close()

	assert.NoError(t, err)
}

func TestSonoffBasicR2_getFullCmndTopic(t *testing.T) {
	sonoffServer, _, err := NewMockMQTTServer()

	assert.NoError(t, err)

	fullCmndTopicWithDefaultFormatOne := fmt.Sprintf("%s/%s/%s", TasmotaPrefixCmnd, "1", TasmotaCmndTopicPower)
	fullCmndTopicWithDefaultFormatTwo := fmt.Sprintf("%s/%s/%s", TasmotaPrefixCmnd, "2", TasmotaCmndTopicStatus)

	assert.Equal(t, fullCmndTopicWithDefaultFormatOne, sonoffServer.getFullCmndTopic("1", TasmotaCmndTopicPower))
	assert.Equal(t, fullCmndTopicWithDefaultFormatTwo, sonoffServer.getFullCmndTopic("2", TasmotaCmndTopicStatus))

	err = sonoffServer.Close()

	assert.NoError(t, err)
}

func TestSonoffBasicR2_getFullTeleTopic(t *testing.T) {
	sonoffServer, _, err := NewMockMQTTServer()

	assert.NoError(t, err)

	fullTeleTopicWithDefaultFormatOne := fmt.Sprintf("%s/%s/%s", TasmotaPrefixTele, "1", TasmotaTeleTopicLWT)
	fullTeleTopicWithDefaultFormatTwo := fmt.Sprintf("%s/%s/%s", TasmotaPrefixTele, "2", TasmotaTeleTopicLWT)

	assert.Equal(t, fullTeleTopicWithDefaultFormatOne, sonoffServer.getFullTeleTopic("1", TasmotaTeleTopicLWT))
	assert.Equal(t, fullTeleTopicWithDefaultFormatTwo, sonoffServer.getFullTeleTopic("2", TasmotaTeleTopicLWT))

	err = sonoffServer.Close()

	assert.NoError(t, err)
}

func TestSonoffBasicR2_getCmndResponse(t *testing.T) {
	sonoffServer, mockServer, err := NewMockMQTTServer()
	responseChan := make(chan string, 1)

	assert.NoError(t, err)

	go func() {
		response, err := sonoffServer.getCmndResponse("1", "TEST", "TEST1", "TEST2")

		assert.NoError(t, err)

		responseChan <- response
	}()

	fullStatTopic := sonoffServer.getFullStatTopic("1", "TEST1")
	fullCmndTopic := sonoffServer.getFullCmndTopic("1", "TEST")

	handler := <-mockServer.subscribeChan
	handler(nil, packets.Subscription{}, packets.Packet{TopicName: fullStatTopic, Payload: []byte("test")})

	assert.Equal(t, "test", <-responseChan)

	assert.Equal(t, fullStatTopic, mockServer.Calls[2].Arguments.Get(0).(string))

	assert.Equal(t, fullCmndTopic, mockServer.Calls[3].Arguments.Get(0).(string))
	assert.Equal(t, []byte("TEST2"), mockServer.Calls[3].Arguments.Get(1).([]byte))

	assert.Equal(t, fullStatTopic, mockServer.Calls[4].Arguments.Get(0).(string))

	err = sonoffServer.Close()

	assert.NoError(t, err)
}
