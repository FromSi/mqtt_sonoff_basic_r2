package mqtt_sonoff_basic_r2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	mqtt "github.com/mochi-mqtt/server/v2"
	mqttauth "github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/mochi-mqtt/server/v2/packets"
	"math"
	"math/rand"
	"strings"
	"time"
)

// SonoffBasicR2 is a struct that manages MQTT connections to a Sonoff Basic R2 device using the Tasmota firmware.
// It handles device commands, status checks, and power control over MQTT.
type SonoffBasicR2 struct {
	server                          *mqtt.Server
	qos                             byte
	isOwnServer                     bool
	connected                       chan string
	disconnected                    chan string
	ctxCmndResponseTimeoutInSeconds uint
	mainContext                     context.Context
	mainContextCancel               context.CancelFunc
}

// NewSonoffBasicR2 initializes a new instance of SonoffBasicR2 and sets up an internal MQTT server.
// It listens for TCP connections on the provided IP and port and allows connections from MQTT clients.
func NewSonoffBasicR2(ip string, port uint16, qos byte) (*SonoffBasicR2, error) {
	server := mqtt.New(
		&mqtt.Options{
			InlineClient: true,
		},
	)

	id := uuid.New().String()
	address := fmt.Sprintf("%s:%d", ip, port)

	tcp := listeners.NewTCP(listeners.Config{ID: id, Address: address})
	err := server.AddListener(tcp)

	if err != nil {
		return nil, err
	}

	// Allow all clients to connect with no authentication
	err = server.AddHook(new(mqttauth.AllowHook), nil)

	if err != nil {
		return nil, err
	}

	mainContext, mainContextCancel := context.WithCancel(context.Background())

	return &SonoffBasicR2{
		server:                          server,
		qos:                             qos,
		isOwnServer:                     true,
		connected:                       make(chan string),
		disconnected:                    make(chan string),
		ctxCmndResponseTimeoutInSeconds: 10,
		mainContext:                     mainContext,
		mainContextCancel:               mainContextCancel,
	}, nil
}

// NewSonoffBasicR2WithServer initializes a SonoffBasicR2 instance with an external MQTT server.
// The server must have the InlineClient option enabled.
func NewSonoffBasicR2WithServer(server *mqtt.Server, qos byte) (*SonoffBasicR2, error) {
	if !server.Options.InlineClient {
		return nil, errors.New("inline_client must be true")
	}

	mainContext, mainContextCancel := context.WithCancel(context.Background())

	return &SonoffBasicR2{
		server:                          server,
		qos:                             qos,
		isOwnServer:                     false,
		connected:                       make(chan string),
		disconnected:                    make(chan string),
		ctxCmndResponseTimeoutInSeconds: 10,
		mainContext:                     mainContext,
		mainContextCancel:               mainContextCancel,
	}, nil
}

// GetCtxCmndResponseTimeoutInSeconds returns the command response timeout duration in seconds.
func (sonoffBasicR2 SonoffBasicR2) GetCtxCmndResponseTimeoutInSeconds() uint {
	return sonoffBasicR2.ctxCmndResponseTimeoutInSeconds
}

// SetCtxCmndResponseTimeoutInSeconds sets the command response timeout duration in seconds.
func (sonoffBasicR2 *SonoffBasicR2) SetCtxCmndResponseTimeoutInSeconds(value uint) {
	sonoffBasicR2.ctxCmndResponseTimeoutInSeconds = value
}

// TeleConnected returns a channel that emits the ID of a device when it is connected to the MQTT broker.
func (sonoffBasicR2 SonoffBasicR2) TeleConnected() <-chan string {
	return sonoffBasicR2.connected
}

// TeleDisconnected returns a channel that emits the ID of a device when it is disconnected from the MQTT broker.
func (sonoffBasicR2 SonoffBasicR2) TeleDisconnected() <-chan string {
	return sonoffBasicR2.disconnected
}

// Serve starts the MQTT server and subscribes to connection status topics for devices.
// It handles the telemetric connection status (`LWT` - Last Will and Testament) from Tasmota devices.
func (sonoffBasicR2 SonoffBasicR2) Serve() error {
	// Subscribe to telemetric messages for connection status (Online/Offline)
	topicTeleConnected := sonoffBasicR2.getFullTeleTopic("+", "LWT")
	subscribeConnected := func(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet) {
		// If the device is online, send the ID to the connected channel
		if string(pk.Payload) == "Online" {
			select {
			case sonoffBasicR2.connected <- strings.Split(pk.TopicName, "/")[1]:
			case <-sonoffBasicR2.mainContext.Done():
			}
		}
	}

	err := sonoffBasicR2.server.Subscribe(topicTeleConnected, sonoffBasicR2.generateSubscriptionId(), subscribeConnected)

	if err != nil {
		return err
	}

	topicTeleDisconnected := sonoffBasicR2.getFullTeleTopic("+", "LWT")
	subscribeDisconnected := func(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet) {
		// If the device is offline, send the ID to the disconnected channel
		if string(pk.Payload) == "Offline" {
			select {
			case sonoffBasicR2.disconnected <- strings.Split(pk.TopicName, "/")[1]:
			case <-sonoffBasicR2.mainContext.Done():
			}
		}
	}

	err = sonoffBasicR2.server.Subscribe(topicTeleDisconnected, sonoffBasicR2.generateSubscriptionId(), subscribeDisconnected)

	if err != nil {
		return err
	}

	// Start the MQTT server if SonoffBasicR2 manages its own server
	if sonoffBasicR2.isOwnServer {
		return sonoffBasicR2.server.Serve()
	}

	return nil
}

// Close closes the MQTT server and stops the internal channels.
func (sonoffBasicR2 SonoffBasicR2) Close() error {
	close(sonoffBasicR2.connected)
	close(sonoffBasicR2.disconnected)

	sonoffBasicR2.mainContextCancel()

	// Close the MQTT server if SonoffBasicR2 manages its own server
	if sonoffBasicR2.isOwnServer {
		return sonoffBasicR2.server.Close()
	}

	return nil
}

// StatusZero retrieves the complete status (STATUS 0) of the Sonoff device via MQTT.
// This includes the device's overall configuration and current state.
// The response is unmarshaled into the AutoGenerated structure.
func (sonoffBasicR2 SonoffBasicR2) StatusZero(id string) (*AutoGenerated, error) {
	response, err := sonoffBasicR2.getCmndResponse(id, "STATUS0", "STATUS0", "")

	if err != nil {
		return nil, err
	}

	return UnmarshalAutoGenerated([]byte(response))
}

// StatusOne retrieves specific system-related information (STATUS 1) from the Sonoff device.
// This includes details like uptime, boot count, and other system parameters.
func (sonoffBasicR2 SonoffBasicR2) StatusOne(id string) (*StatusOne, error) {
	response, err := sonoffBasicR2.getCmndResponse(id, "STATUS", "STATUS1", "1")

	if err != nil {
		return nil, err
	}

	return UnmarshalStatusOne([]byte(response))
}

// StatusTwo retrieves firmware-related information (STATUS 2) from the Sonoff device.
// This includes firmware version, build date, and other firmware-specific data.
func (sonoffBasicR2 SonoffBasicR2) StatusTwo(id string) (*StatusTwo, error) {
	response, err := sonoffBasicR2.getCmndResponse(id, "STATUS", "STATUS2", "2")

	if err != nil {
		return nil, err
	}

	return UnmarshalStatusTwo([]byte(response))
}

// StatusThree retrieves logging-related settings (STATUS 3) from the Sonoff device.
// This includes serial, web, and MQTT log configurations.
func (sonoffBasicR2 SonoffBasicR2) StatusThree(id string) (*StatusThree, error) {
	response, err := sonoffBasicR2.getCmndResponse(id, "STATUS", "STATUS3", "3")

	if err != nil {
		return nil, err
	}

	return UnmarshalStatusThree([]byte(response))
}

// StatusFour retrieves memory and storage-related information (STATUS 4) from the Sonoff device.
// This includes program size, free heap space, flash size, and other memory metrics.
func (sonoffBasicR2 SonoffBasicR2) StatusFour(id string) (*StatusFour, error) {
	response, err := sonoffBasicR2.getCmndResponse(id, "STATUS", "STATUS4", "4")

	if err != nil {
		return nil, err
	}

	return UnmarshalStatusFour([]byte(response))
}

// StatusFive retrieves network configuration details (STATUS 5) from the Sonoff device.
// This includes IP address, gateway, subnet mask, and DNS server information.
func (sonoffBasicR2 SonoffBasicR2) StatusFive(id string) (*StatusFive, error) {
	response, err := sonoffBasicR2.getCmndResponse(id, "STATUS", "STATUS5", "5")

	if err != nil {
		return nil, err
	}

	return UnmarshalStatusFive([]byte(response))
}

// StatusSix retrieves MQTT configuration information (STATUS 6) from the Sonoff device.
// This includes MQTT host, port, client ID, and other MQTT settings.
func (sonoffBasicR2 SonoffBasicR2) StatusSix(id string) (*StatusSix, error) {
	response, err := sonoffBasicR2.getCmndResponse(id, "STATUS", "STATUS6", "6")

	if err != nil {
		return nil, err
	}

	return UnmarshalStatusSix([]byte(response))
}

// StatusSeven retrieves time and date settings (STATUS 7) from the Sonoff device.
// This includes local time, daylight savings settings, and timezone information.
func (sonoffBasicR2 SonoffBasicR2) StatusSeven(id string) (*StatusSeven, error) {
	response, err := sonoffBasicR2.getCmndResponse(id, "STATUS", "STATUS7", "7")

	if err != nil {
		return nil, err
	}

	return UnmarshalStatusSeven([]byte(response))
}

// StatusEight retrieves sensor data (STATUS 8) from the Sonoff device.
// This includes the most recent readings from the device's sensors.
func (sonoffBasicR2 SonoffBasicR2) StatusEight(id string) (*StatusEight, error) {
	response, err := sonoffBasicR2.getCmndResponse(id, "STATUS", "STATUS8", "8")

	if err != nil {
		return nil, err
	}

	return UnmarshalStatusEight([]byte(response))
}

// StatusEleven retrieves runtime status information (STATUS 11) from the Sonoff device.
// This includes uptime, heap usage, WiFi information, and more.
func (sonoffBasicR2 SonoffBasicR2) StatusEleven(id string) (*StatusEleven, error) {
	response, err := sonoffBasicR2.getCmndResponse(id, "STATUS", "STATUS11", "11")

	if err != nil {
		return nil, err
	}

	return UnmarshalStatusEleven([]byte(response))
}

// StatusPhysicalButton retrieves the current configuration of the physical button on the Sonoff device.
// It checks the status of SetOption73, which controls whether the physical button is enabled (OFF) or disabled (ON).
func (sonoffBasicR2 SonoffBasicR2) StatusPhysicalButton(id string) (bool, error) {
	response, err := sonoffBasicR2.getCmndResponse(id, "SETOPTION73", "RESULT", "")

	if err != nil {
		return false, err
	}

	var data map[string]string

	if err := json.Unmarshal([]byte(response), &data); err != nil {
		return false, err
	}

	result, ok := data["SetOption73"]

	if !ok {
		return false, errors.New("SetOption73 not found")
	}

	// If SetOption73 is "OFF", the physical button is enabled; otherwise, it's disabled.
	return result == "OFF", nil
}

// PowerOn sends an MQTT command to turn on the device.
func (sonoffBasicR2 SonoffBasicR2) PowerOn(id string) {
	topicCmnd := sonoffBasicR2.getFullCmndTopic(id, "POWER")
	_ = sonoffBasicR2.server.Publish(topicCmnd, []byte("ON"), false, sonoffBasicR2.qos)
}

// PowerOff sends an MQTT command to turn off the device.
func (sonoffBasicR2 SonoffBasicR2) PowerOff(id string) {
	topicCmnd := sonoffBasicR2.getFullCmndTopic(id, "POWER")
	_ = sonoffBasicR2.server.Publish(topicCmnd, []byte("OFF"), false, sonoffBasicR2.qos)
}

// PowerToggle sends an MQTT command to toggle the power state of the Sonoff device.
// It switches the power between ON and OFF, depending on the current state.
func (sonoffBasicR2 SonoffBasicR2) PowerToggle(id string) {
	topicCmnd := sonoffBasicR2.getFullCmndTopic(id, "POWER")
	_ = sonoffBasicR2.server.Publish(topicCmnd, []byte("TOGGLE"), false, sonoffBasicR2.qos)
}

// PhysicalButtonOn sends an MQTT command to enable the physical button on the Sonoff device.
// This allows the device's physical button to control power toggling. It corresponds to the Tasmota command SetOption73.
func (sonoffBasicR2 SonoffBasicR2) PhysicalButtonOn(id string) {
	topicCmnd := sonoffBasicR2.getFullCmndTopic(id, "SETOPTION73")
	_ = sonoffBasicR2.server.Publish(topicCmnd, []byte("0"), false, sonoffBasicR2.qos)
}

// PhysicalButtonOff sends an MQTT command to disable the physical button on the Sonoff device.
// This prevents the device's physical button from toggling the power. It corresponds to the Tasmota command SetOption73.
func (sonoffBasicR2 SonoffBasicR2) PhysicalButtonOff(id string) {
	topicCmnd := sonoffBasicR2.getFullCmndTopic(id, "SETOPTION73")
	_ = sonoffBasicR2.server.Publish(topicCmnd, []byte("1"), false, sonoffBasicR2.qos)
}

// generateSubscriptionId generates a unique subscription ID for MQTT topics using a random number generator.
func (sonoffBasicR2 SonoffBasicR2) generateSubscriptionId() int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(math.MaxInt32)
}

// Helper methods for constructing the full MQTT topic paths for command (cmnd), telemetry (tele), and status (stat) topics.
func (sonoffBasicR2 SonoffBasicR2) getFullTopic(prefix string, id string, topic string) string {
	return fmt.Sprintf("%s/%s/%s", prefix, id, topic)
}

// getFullStatTopic constructs the full MQTT topic for device status ("stat") messages.
// It combines the topic prefix "stat", the device ID, and the specific status topic.
func (sonoffBasicR2 SonoffBasicR2) getFullStatTopic(id string, topic string) string {
	return sonoffBasicR2.getFullTopic("stat", id, topic)
}

// getFullCmndTopic constructs the full MQTT topic for command ("cmnd") messages.
// It combines the topic prefix "cmnd", the device ID, and the specific command topic.
func (sonoffBasicR2 SonoffBasicR2) getFullCmndTopic(id string, topic string) string {
	return sonoffBasicR2.getFullTopic("cmnd", id, topic)
}

// getFullTeleTopic constructs the full MQTT topic for telemetry ("tele") messages.
// It combines the topic prefix "tele", the device ID, and the specific telemetry topic.
func (sonoffBasicR2 SonoffBasicR2) getFullTeleTopic(id string, topic string) string {
	return sonoffBasicR2.getFullTopic("tele", id, topic)
}

// getCmndResponse sends a command to the Sonoff device and waits for a response.
// It publishes the command on the "cmnd" topic and subscribes to the corresponding "stat" topic to capture the response.
// If a response is not received within the defined timeout, it returns an error.
func (sonoffBasicR2 SonoffBasicR2) getCmndResponse(id string, topicCmnd string, topicStat string, value string) (string, error) {
	// Get the full topic for status and command
	fullTopicStat := sonoffBasicR2.getFullStatTopic(id, topicStat)
	fullTopicCmnd := sonoffBasicR2.getFullCmndTopic(id, topicCmnd)

	// Set a timeout for the response
	ctx, cancel := context.WithTimeout(
		sonoffBasicR2.mainContext,
		time.Duration(sonoffBasicR2.ctxCmndResponseTimeoutInSeconds)*time.Second,
	)

	defer cancel()

	// Channel to capture the response
	result := make(chan string, 1)

	// Generate a unique subscription ID
	subscriptionId := sonoffBasicR2.generateSubscriptionId()

	// Function to handle incoming status messages
	subscribeResponse := func(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet) {
		select {
		case result <- string(pk.Payload):
		case <-ctx.Done():
		}
	}

	// Subscribe to the status topic to receive the response
	err := sonoffBasicR2.server.Subscribe(fullTopicStat, subscriptionId, subscribeResponse)

	defer func(server *mqtt.Server, filter string, subscriptionId int) {
		_ = server.Unsubscribe(filter, subscriptionId)
	}(sonoffBasicR2.server, fullTopicStat, subscriptionId)

	if err != nil {
		return "", err
	}

	// Publish the command to the device
	_ = sonoffBasicR2.server.Publish(fullTopicCmnd, []byte(value), false, sonoffBasicR2.qos)

	// Wait for a response or timeout
	select {
	case <-ctx.Done():
		return "", fmt.Errorf(
			"operation not completed in %d seconds",
			sonoffBasicR2.ctxCmndResponseTimeoutInSeconds,
		)
	case data := <-result:
		return data, nil
	}
}
