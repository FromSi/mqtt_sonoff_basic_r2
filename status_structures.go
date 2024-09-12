package mqtt_sonoff_basic_r2

import (
	"encoding/json"
	"fmt"
	"time"
)

type TasmotaTime time.Time

func (tasmotaTime *TasmotaTime) UnmarshalJSON(value []byte) error {
	str := string(value[1 : len(value)-1])

	t, err := time.Parse("2006-01-02T15:04:05", str)

	if err != nil {
		return err
	}

	*tasmotaTime = TasmotaTime(t)

	return nil
}

func (tasmotaTime TasmotaTime) MarshalJSON() ([]byte, error) {
	t := time.Time(tasmotaTime)

	return []byte(fmt.Sprintf(`"%s"`, t.Format("2006-01-02T15:04:05"))), nil
}

func (tasmotaTime TasmotaTime) ToTime() time.Time {
	return time.Time(tasmotaTime)
}

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

func UnmarshalStatusZero(data []byte) (*StatusZero, error) {
	var result StatusZero

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

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

func UnmarshalStatusOne(data []byte) (*StatusOne, error) {
	var mapResult map[string]StatusOne

	if err := json.Unmarshal(data, &mapResult); err != nil {
		return nil, err
	}

	result := mapResult["StatusPRM"]

	return &result, nil
}

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

func UnmarshalStatusTwo(data []byte) (*StatusTwo, error) {
	var mapResult map[string]StatusTwo

	if err := json.Unmarshal(data, &mapResult); err != nil {
		return nil, err
	}

	result := mapResult["StatusFWR"]

	return &result, nil
}

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

func UnmarshalStatusThree(data []byte) (*StatusThree, error) {
	var mapResult map[string]StatusThree

	if err := json.Unmarshal(data, &mapResult); err != nil {
		return nil, err
	}

	result := mapResult["StatusLOG"]

	return &result, nil
}

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

func UnmarshalStatusFour(data []byte) (*StatusFour, error) {
	var mapResult map[string]StatusFour

	if err := json.Unmarshal(data, &mapResult); err != nil {
		return nil, err
	}

	result := mapResult["StatusMEM"]

	return &result, nil
}

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

func UnmarshalStatusFive(data []byte) (*StatusFive, error) {
	var mapResult map[string]StatusFive

	if err := json.Unmarshal(data, &mapResult); err != nil {
		return nil, err
	}

	result := mapResult["StatusNET"]

	return &result, nil
}

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

func UnmarshalStatusSix(data []byte) (*StatusSix, error) {
	var mapResult map[string]StatusSix

	if err := json.Unmarshal(data, &mapResult); err != nil {
		return nil, err
	}

	result := mapResult["StatusMQT"]

	return &result, nil
}

type StatusSeven struct {
	UTC      time.Time   `json:"UTC"`
	Local    TasmotaTime `json:"Local"`
	StartDST TasmotaTime `json:"StartDST"`
	EndDST   TasmotaTime `json:"EndDST"`
	Timezone string      `json:"Timezone"`
	Sunrise  string      `json:"Sunrise"`
	Sunset   string      `json:"Sunset"`
}

func UnmarshalStatusSeven(data []byte) (*StatusSeven, error) {
	var mapResult map[string]StatusSeven

	if err := json.Unmarshal(data, &mapResult); err != nil {
		return nil, err
	}

	result := mapResult["StatusTIM"]

	return &result, nil
}

type StatusEight struct {
	Time TasmotaTime `json:"Time"`
}

func UnmarshalStatusEight(data []byte) (*StatusEight, error) {
	var mapResult map[string]StatusEight

	if err := json.Unmarshal(data, &mapResult); err != nil {
		return nil, err
	}

	result := mapResult["StatusSNS"]

	return &result, nil
}

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

func UnmarshalStatusEleven(data []byte) (*StatusEleven, error) {
	var mapResult map[string]StatusEleven

	if err := json.Unmarshal(data, &mapResult); err != nil {
		return nil, err
	}

	result := mapResult["StatusSTS"]

	return &result, nil
}

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

func UnmarshalAutoGenerated(data []byte) (*AutoGenerated, error) {
	var result AutoGenerated

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
