package mqtt_sonoff_basic_r2

import (
	"encoding/json"
	"fmt"
	"time"
)

// TasmotaTime is a custom time type used to parse and format Tasmota's specific date-time format (YYYY-MM-DDTHH:MM:SS).
// Tasmota is a firmware for ESP8266 devices, like Sonoff, which communicates over MQTT and provides device status in JSON format.
// More about Tasmota: https://tasmota.github.io/docs/
type TasmotaTime time.Time

// UnmarshalJSON handles parsing the Tasmota-specific time format when unmarshaling JSON data.
func (tasmotaTime *TasmotaTime) UnmarshalJSON(value []byte) error {
	str := string(value[1 : len(value)-1]) // Removing the enclosing quotes

	t, err := time.Parse("2006-01-02T15:04:05", str) // Parse the string using Tasmota's format

	if err != nil {
		return err
	}

	*tasmotaTime = TasmotaTime(t)

	return nil
}

// MarshalJSON handles converting TasmotaTime into a string in the Tasmota-specific format when marshaling JSON data.
func (tasmotaTime TasmotaTime) MarshalJSON() ([]byte, error) {
	t := time.Time(tasmotaTime)

	return []byte(fmt.Sprintf(`"%s"`, t.Format("2006-01-02T15:04:05"))), nil
}

// ToTime converts TasmotaTime back to the standard Go time.Time type.
func (tasmotaTime TasmotaTime) ToTime() time.Time {
	return time.Time(tasmotaTime)
}

// StatusZero represents the JSON structure of the device status information returned by Tasmota.
// It contains general information about the device such as the module type, power state, and more.
// See: https://tasmota.github.io/docs/Commands/#management
type StatusZero struct {
	Module       int      `json:"Module"`
	DeviceName   string   `json:"DeviceName"`
	FriendlyName []string `json:"FriendlyName"`
	Topic        string   `json:"Topic"`
	ButtonTopic  string   `json:"ButtonTopic"`
	Power        string   `json:"Power"`
	PowerLock    string   `json:"PowerLock"`
	PowerOnState int      `json:"PowerOnState"`
	LedState     int      `json:"LedState"`
	LedMask      string   `json:"LedMask"`
	SaveData     int      `json:"SaveData"`
	SaveState    int      `json:"SaveState"`
	SwitchTopic  string   `json:"SwitchTopic"`
	SwitchMode   []int    `json:"SwitchMode"`
	ButtonRetain int      `json:"ButtonRetain"`
	SwitchRetain int      `json:"SwitchRetain"`
	SensorRetain int      `json:"SensorRetain"`
	PowerRetain  int      `json:"PowerRetain"`
	InfoRetain   int      `json:"InfoRetain"`
	StateRetain  int      `json:"StateRetain"`
	StatusRetain int      `json:"StatusRetain"`
}

// UnmarshalStatusZero unmarshals the device status (Status 0) from JSON data.
func UnmarshalStatusZero(data []byte) (*StatusZero, error) {
	var result StatusZero

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// StatusOne represents system and configuration details of a Tasmota device, such as baudrate and OTA URL.
// It corresponds to the Status 1 command in Tasmota.
// See: https://tasmota.github.io/docs/Commands/#management
type StatusOne struct {
	Baudrate      int         `json:"Baudrate"`
	SerialConfig  string      `json:"SerialConfig"`
	GroupTopic    string      `json:"GroupTopic"`
	OtaURL        string      `json:"OtaUrl"`
	RestartReason string      `json:"RestartReason"`
	Uptime        string      `json:"Uptime"`
	StartupUTC    TasmotaTime `json:"StartupUTC"`
	Sleep         int         `json:"Sleep"`
	CfgHolder     int         `json:"CfgHolder"`
	BootCount     int         `json:"BootCount"`
	BCResetTime   TasmotaTime `json:"BCResetTime"`
	SaveCount     int         `json:"SaveCount"`
	SaveAddress   string      `json:"SaveAddress"`
}

// UnmarshalStatusOne unmarshals the system status (Status 1) from JSON data.
func UnmarshalStatusOne(data []byte) (*StatusOne, error) {
	var mapResult map[string]StatusOne

	if err := json.Unmarshal(data, &mapResult); err != nil {
		return nil, err
	}

	result := mapResult["StatusPRM"]

	return &result, nil
}

// StatusTwo represents firmware details of the Tasmota device, such as version and build date.
// It corresponds to the Status 2 command in Tasmota.
// See: https://tasmota.github.io/docs/Commands/#management
type StatusTwo struct {
	Version       string      `json:"Version"`
	BuildDateTime TasmotaTime `json:"BuildDateTime"`
	Boot          int         `json:"Boot"`
	Core          string      `json:"Core"`
	SDK           string      `json:"SDK"`
	CPUFrequency  int         `json:"CpuFrequency"`
	Hardware      string      `json:"Hardware"`
	CR            string      `json:"CR"`
}

// UnmarshalStatusTwo unmarshals the firmware status (Status 2) from JSON data.
func UnmarshalStatusTwo(data []byte) (*StatusTwo, error) {
	var mapResult map[string]StatusTwo

	if err := json.Unmarshal(data, &mapResult); err != nil {
		return nil, err
	}

	result := mapResult["StatusFWR"]

	return &result, nil
}

// StatusThree represents the logging configuration of the device, including log levels and host settings.
// It corresponds to the Status 3 command in Tasmota.
// See: https://tasmota.github.io/docs/Commands/#management
type StatusThree struct {
	SerialLog  int      `json:"SerialLog"`
	WebLog     int      `json:"WebLog"`
	MqttLog    int      `json:"MqttLog"`
	SysLog     int      `json:"SysLog"`
	LogHost    string   `json:"LogHost"`
	LogPort    int      `json:"LogPort"`
	SSID       []string `json:"SSId"`
	TelePeriod int      `json:"TelePeriod"`
	Resolution string   `json:"Resolution"`
	SetOption  []string `json:"SetOption"`
}

// UnmarshalStatusThree unmarshals the logging status (Status 3) from JSON data.
func UnmarshalStatusThree(data []byte) (*StatusThree, error) {
	var mapResult map[string]StatusThree

	if err := json.Unmarshal(data, &mapResult); err != nil {
		return nil, err
	}

	result := mapResult["StatusLOG"]

	return &result, nil
}

// StatusFour represents the memory status of the device, including flash memory and heap size.
// It corresponds to the Status 4 command in Tasmota.
// See: https://tasmota.github.io/docs/Commands/#management
type StatusFour struct {
	ProgramSize      int      `json:"ProgramSize"`
	Free             int      `json:"Free"`
	Heap             int      `json:"Heap"`
	ProgramFlashSize int      `json:"ProgramFlashSize"`
	FlashSize        int      `json:"FlashSize"`
	FlashChipID      string   `json:"FlashChipId"`
	FlashFrequency   int      `json:"FlashFrequency"`
	FlashMode        string   `json:"FlashMode"`
	Features         []string `json:"Features"`
	Drivers          string   `json:"Drivers"`
	Sensors          string   `json:"Sensors"`
	I2CDriver        string   `json:"I2CDriver"`
}

// UnmarshalStatusFour unmarshals the memory status (Status 4) from JSON data.
func UnmarshalStatusFour(data []byte) (*StatusFour, error) {
	var mapResult map[string]StatusFour

	if err := json.Unmarshal(data, &mapResult); err != nil {
		return nil, err
	}

	result := mapResult["StatusMEM"]

	return &result, nil
}

// StatusFive represents the network configuration of the device, including IP addresses, DNS settings, and MAC address.
// It corresponds to the Status 5 command in Tasmota.
// See: https://tasmota.github.io/docs/Commands/#management
type StatusFive struct {
	Hostname   string  `json:"Hostname"`
	IPAddress  string  `json:"IPAddress"`
	Gateway    string  `json:"Gateway"`
	Subnetmask string  `json:"Subnetmask"`
	DNSServer1 string  `json:"DNSServer1"`
	DNSServer2 string  `json:"DNSServer2"`
	Mac        string  `json:"Mac"`
	Webserver  int     `json:"Webserver"`
	HTTPAPI    int     `json:"HTTP_API"`
	WifiConfig int     `json:"WifiConfig"`
	WifiPower  float64 `json:"WifiPower"`
}

// UnmarshalStatusFive unmarshals the network configuration (Status 5) from JSON data.
func UnmarshalStatusFive(data []byte) (*StatusFive, error) {
	var mapResult map[string]StatusFive

	if err := json.Unmarshal(data, &mapResult); err != nil {
		return nil, err
	}

	result := mapResult["StatusNET"]

	return &result, nil
}

// StatusSix represents MQTT configuration and settings of the Tasmota device, such as MQTT host, port, and client information.
// It corresponds to the Status 6 command in Tasmota.
// See: https://tasmota.github.io/docs/Commands/#management
type StatusSix struct {
	MqttHost       string `json:"MqttHost"`
	MqttPort       int    `json:"MqttPort"`
	MqttClientMask string `json:"MqttClientMask"`
	MqttClient     string `json:"MqttClient"`
	MqttUser       string `json:"MqttUser"`
	MqttCount      int    `json:"MqttCount"`
	MAXPACKETSIZE  int    `json:"MAX_PACKET_SIZE"`
	KEEPALIVE      int    `json:"KEEPALIVE"`
	SOCKETTIMEOUT  int    `json:"SOCKET_TIMEOUT"`
}

// UnmarshalStatusSix unmarshals the MQTT configuration (Status 6) from JSON data.
func UnmarshalStatusSix(data []byte) (*StatusSix, error) {
	var mapResult map[string]StatusSix

	if err := json.Unmarshal(data, &mapResult); err != nil {
		return nil, err
	}

	result := mapResult["StatusMQT"]

	return &result, nil
}

// StatusSeven represents the time settings of the device, including local time, daylight saving time (DST), and time zone.
// It corresponds to the Status 7 command in Tasmota.
// See: https://tasmota.github.io/docs/Commands/#management
type StatusSeven struct {
	UTC      time.Time   `json:"UTC"`
	Local    TasmotaTime `json:"Local"`
	StartDST TasmotaTime `json:"StartDST"`
	EndDST   TasmotaTime `json:"EndDST"`
	Timezone string      `json:"Timezone"`
	Sunrise  string      `json:"Sunrise"`
	Sunset   string      `json:"Sunset"`
}

// UnmarshalStatusSeven unmarshals the time settings (Status 7) from JSON data.
func UnmarshalStatusSeven(data []byte) (*StatusSeven, error) {
	var mapResult map[string]StatusSeven

	if err := json.Unmarshal(data, &mapResult); err != nil {
		return nil, err
	}

	result := mapResult["StatusTIM"]

	return &result, nil
}

// StatusEight represents sensor data collected from the device.
// It corresponds to the Status 8 command in Tasmota, also known as StatusSNS (Sensor Status).
// See: https://tasmota.github.io/docs/Commands/#management
type StatusEight struct {
	Time TasmotaTime `json:"Time"`
}

// UnmarshalStatusEight unmarshals the sensor data (Status 8) from JSON data.
func UnmarshalStatusEight(data []byte) (*StatusEight, error) {
	var mapResult map[string]StatusEight

	if err := json.Unmarshal(data, &mapResult); err != nil {
		return nil, err
	}

	result := mapResult["StatusSNS"]

	return &result, nil
}

// StatusEleven represents the runtime status of the device, including uptime, heap usage, and WiFi information.
// It corresponds to the Status 11 command in Tasmota, also known as StatusSTS (Status System).
// See: https://tasmota.github.io/docs/Commands/#management
type StatusEleven struct {
	Time      TasmotaTime `json:"Time"`
	Uptime    string      `json:"Uptime"`
	UptimeSec int         `json:"UptimeSec"`
	Heap      int         `json:"Heap"`
	SleepMode string      `json:"SleepMode"`
	Sleep     int         `json:"Sleep"`
	LoadAvg   int         `json:"LoadAvg"`
	MqttCount int         `json:"MqttCount"`
	POWER     string      `json:"POWER"`
	Wifi      struct {
		AP        int    `json:"AP"`
		SSID      string `json:"SSId"`
		BSSID     string `json:"BSSId"`
		Channel   int    `json:"Channel"`
		Mode      string `json:"Mode"`
		RSSI      int    `json:"RSSI"`
		Signal    int    `json:"Signal"`
		LinkCount int    `json:"LinkCount"`
		Downtime  string `json:"Downtime"`
	} `json:"Wifi"`
}

// UnmarshalStatusEleven unmarshals the runtime status (Status 11) from JSON data.
func UnmarshalStatusEleven(data []byte) (*StatusEleven, error) {
	var mapResult map[string]StatusEleven

	if err := json.Unmarshal(data, &mapResult); err != nil {
		return nil, err
	}

	result := mapResult["StatusSTS"]

	return &result, nil
}

// AutoGenerated combines various status commands into one structure.
// It represents the entire set of status information returned by a Tasmota device, including system, network, and sensor data.
// See the full Tasmota documentation: https://tasmota.github.io/docs/Commands/#management
type AutoGenerated struct {
	Status    StatusZero   `json:"Status"`
	StatusPRM StatusOne    `json:"StatusPRM"`
	StatusFWR StatusTwo    `json:"StatusFWR"`
	StatusLOG StatusThree  `json:"StatusLOG"`
	StatusMEM StatusFour   `json:"StatusMEM"`
	StatusNET StatusFive   `json:"StatusNET"`
	StatusMQT StatusSix    `json:"StatusMQT"`
	StatusTIM StatusSeven  `json:"StatusTIM"`
	StatusSNS StatusEight  `json:"StatusSNS"`
	StatusSTS StatusEleven `json:"StatusSTS"`
}

// UnmarshalAutoGenerated unmarshals the entire set of Tasmota status information from JSON data.
func UnmarshalAutoGenerated(data []byte) (*AutoGenerated, error) {
	var result AutoGenerated

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
