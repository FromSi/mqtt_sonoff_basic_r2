package mqtt_sonoff_basic_r2

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const JsonData = `{
	  "Status": {
		"Module": 1,
		"DeviceName": "Tasmota",
		"FriendlyName": [
		  "Tasmota"
		],
		"Topic": "main",
		"ButtonTopic": "0",
		"Power": "0",
		"PowerLock": "0",
		"PowerOnState": 3,
		"LedState": 1,
		"LedMask": "FFFF",
		"SaveData": 1,
		"SaveState": 1,
		"SwitchTopic": "0",
		"SwitchMode": [
		  0,
		  0,
		  0,
		  0,
		  0,
		  0,
		  0,
		  0,
		  0,
		  0,
		  0,
		  0,
		  0,
		  0,
		  0,
		  0,
		  0,
		  0,
		  0,
		  0,
		  0,
		  0,
		  0,
		  0,
		  0,
		  0,
		  0,
		  0
		],
		"ButtonRetain": 0,
		"SwitchRetain": 0,
		"SensorRetain": 0,
		"PowerRetain": 0,
		"InfoRetain": 0,
		"StateRetain": 0,
		"StatusRetain": 0
	  },
	  "StatusPRM": {
		"Baudrate": 115200,
		"SerialConfig": "8N1",
		"GroupTopic": "tasmotas",
		"OtaUrl": "http://ota.tasmota.com/tasmota/release/tasmota.bin.gz",
		"RestartReason": "Software/System restart",
		"Uptime": "0T00:27:03",
		"StartupUTC": "2024-08-31T12:50:38",
		"Sleep": 50,
		"CfgHolder": 4617,
		"BootCount": 14,
		"BCResetTime": "2024-08-29T15:59:39",
		"SaveCount": 97,
		"SaveAddress": "FB000"
	  },
	  "StatusFWR": {
		"Version": "14.2.0(release-tasmota)",
		"BuildDateTime": "2024-08-14T12:36:35",
		"Boot": 31,
		"Core": "2_7_7",
		"SDK": "2.2.2-dev(38a443e)",
		"CpuFrequency": 80,
		"Hardware": "ESP8285N08",
		"CR": "373/699"
	  },
	  "StatusLOG": {
		"SerialLog": 2,
		"WebLog": 2,
		"MqttLog": 0,
		"SysLog": 0,
		"LogHost": "",
		"LogPort": 514,
		"SSId": [
		  "ALHN-ED72",
		  ""
		],
		"TelePeriod": 300,
		"Resolution": "558180C0",
		"SetOption": [
		  "00008009",
		  "2805C80001000600003C5A0A192800000000",
		  "00000080",
		  "00006000",
		  "00004000",
		  "00000000"
		]
	  },
	  "StatusMEM": {
		"ProgramSize": 648,
		"Free": 352,
		"Heap": 23,
		"ProgramFlashSize": 1024,
		"FlashSize": 1024,
		"FlashChipId": "144051",
		"FlashFrequency": 40,
		"FlashMode": "DOUT",
		"Features": [
		  "0809",
		  "8F9AC787",
		  "04368001",
		  "000000CF",
		  "010013C0",
		  "C000F981",
		  "00004004",
		  "00001000",
		  "54000020",
		  "00000080",
		  "00000000"
		],
		"Drivers": "1,2,!3,!4,!5,!6,7,!8,9,10,12,!16,!18,!19,!20,!21,!22,!24,26,!27,29,!30,!35,!37,!45,62,!68",
		"Sensors": "1,2,3,4,5,6",
		"I2CDriver": "7"
	  },
	  "StatusNET": {
		"Hostname": "main-0614",
		"IPAddress": "192.168.1.158",
		"Gateway": "192.168.1.1",
		"Subnetmask": "255.255.255.0",
		"DNSServer1": "192.168.1.1",
		"DNSServer2": "0.0.0.0",
		"Mac": "2C:F4:32:FA:22:66",
		"Webserver": 2,
		"HTTP_API": 1,
		"WifiConfig": 4,
		"WifiPower": 17
	  },
	  "StatusMQT": {
		"MqttHost": "192.168.1.129",
		"MqttPort": 1883,
		"MqttClientMask": "DVES_%06X",
		"MqttClient": "DVES_FA2266",
		"MqttUser": "DVES_USER",
		"MqttCount": 1,
		"MAX_PACKET_SIZE": 1200,
		"KEEPALIVE": 30,
		"SOCKET_TIMEOUT": 4
	  },
	  "StatusTIM": {
		"UTC": "2024-08-31T13:17:41Z",
		"Local": "2024-08-31T14:17:41",
		"StartDST": "2024-03-31T02:00:00",
		"EndDST": "2024-10-27T03:00:00",
		"Timezone": "+01:00",
		"Sunrise": "06:06",
		"Sunset": "19:33"
	  },
	  "StatusSNS": {
		"Time": "2024-08-31T14:17:41"
	  },
	  "StatusSTS": {
		"Time": "2024-08-31T14:17:41",
		"Uptime": "0T00:27:03",
		"UptimeSec": 1623,
		"Heap": 22,
		"SleepMode": "Dynamic",
		"Sleep": 50,
		"LoadAvg": 19,
		"MqttCount": 1,
		"POWER": "OFF",
		"Wifi": {
		  "AP": 1,
		  "SSId": "ALHN-ED72",
		  "BSSId": "EC:84:B4:0C:86:09",
		  "Channel": 3,
		  "Mode": "11n",
		  "RSSI": 68,
		  "Signal": -66,
		  "LinkCount": 1,
		  "Downtime": "0T00:00:04"
		}
	  }
	}`

func TestTasmotaTime_UnmarshalJSON(t *testing.T) {
	jsonData := `"2023-09-13T14:20:00"`

	var tasmotaTime TasmotaTime

	err := json.Unmarshal([]byte(jsonData), &tasmotaTime)

	if !assert.NoError(t, err) {
		return
	}

	expectedTime := time.Date(2023, 9, 13, 14, 20, 0, 0, time.UTC)

	assert.Equal(t, expectedTime, tasmotaTime.ToTime())
}

func TestTasmotaTime_MarshalJSON(t *testing.T) {
	tasmotaTime := TasmotaTime(time.Date(2023, 9, 13, 14, 20, 0, 0, time.UTC))

	jsonData, err := json.Marshal(tasmotaTime)

	if !assert.NoError(t, err) {
		return
	}

	expectedJson := `"2023-09-13T14:20:00"`

	assert.JSONEq(t, expectedJson, string(jsonData))
}

func TestTasmotaTime_ToTime(t *testing.T) {
	expectedTime := time.Date(2023, 9, 13, 14, 20, 0, 0, time.UTC)
	tasmotaTime := TasmotaTime(expectedTime)

	resultTime := tasmotaTime.ToTime()

	assert.Equal(t, expectedTime, resultTime)
}

func Test_UnmarshalStatusZero(t *testing.T) {
	statusZero, err := UnmarshalStatusZero([]byte(JsonData))

	if !assert.NoError(t, err) {
		return
	}

	expectedStatusZero := StatusZero{
		Module:       1,
		DeviceName:   "Tasmota",
		FriendlyName: []string{"Tasmota"},
		Topic:        "main",
		ButtonTopic:  "0",
		Power:        "0",
		PowerLock:    "0",
		PowerOnState: 3,
		LedState:     1,
		LedMask:      "FFFF",
		SaveData:     1,
		SaveState:    1,
		SwitchTopic:  "0",
		SwitchMode:   []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		ButtonRetain: 0,
		SwitchRetain: 0,
		SensorRetain: 0,
		PowerRetain:  0,
		InfoRetain:   0,
		StateRetain:  0,
		StatusRetain: 0,
	}

	assert.Equal(t, expectedStatusZero, *statusZero)
}

func Test_UnmarshalStatusOne(t *testing.T) {
	statusOne, err := UnmarshalStatusOne([]byte(JsonData))

	if !assert.NoError(t, err) {
		return
	}

	expectedStatusOne := StatusOne{
		Baudrate:      115200,
		SerialConfig:  "8N1",
		GroupTopic:    "tasmotas",
		OtaURL:        "http://ota.tasmota.com/tasmota/release/tasmota.bin.gz",
		RestartReason: "Software/System restart",
		Uptime:        "0T00:27:03",
		StartupUTC:    TasmotaTime(time.Date(2024, 8, 31, 12, 50, 38, 0, time.UTC)),
		Sleep:         50,
		CfgHolder:     4617,
		BootCount:     14,
		BCResetTime:   TasmotaTime(time.Date(2024, 8, 29, 15, 59, 39, 0, time.UTC)),
		SaveCount:     97,
		SaveAddress:   "FB000",
	}

	assert.Equal(t, expectedStatusOne, *statusOne)
}

func Test_UnmarshalStatusTwo(t *testing.T) {
	statusTwo, err := UnmarshalStatusTwo([]byte(JsonData))

	if !assert.NoError(t, err) {
		return
	}

	expectedStatusTwo := StatusTwo{
		Version:       "14.2.0(release-tasmota)",
		BuildDateTime: TasmotaTime(time.Date(2024, 8, 14, 12, 36, 35, 0, time.UTC)),
		Boot:          31,
		Core:          "2_7_7",
		SDK:           "2.2.2-dev(38a443e)",
		CPUFrequency:  80,
		Hardware:      "ESP8285N08",
		CR:            "373/699",
	}

	assert.Equal(t, expectedStatusTwo, *statusTwo)
}

func Test_UnmarshalStatusThree(t *testing.T) {
	statusThree, err := UnmarshalStatusThree([]byte(JsonData))

	if !assert.NoError(t, err) {
		return
	}

	expectedStatusThree := StatusThree{
		SerialLog:  2,
		WebLog:     2,
		MqttLog:    0,
		SysLog:     0,
		LogHost:    "",
		LogPort:    514,
		SSID:       []string{"ALHN-ED72", ""},
		TelePeriod: 300,
		Resolution: "558180C0",
		SetOption:  []string{"00008009", "2805C80001000600003C5A0A192800000000", "00000080", "00006000", "00004000", "00000000"},
	}

	assert.Equal(t, expectedStatusThree, *statusThree)
}

func Test_UnmarshalStatusFour(t *testing.T) {
	statusFour, err := UnmarshalStatusFour([]byte(JsonData))

	if !assert.NoError(t, err) {
		return
	}

	expectedStatusFour := StatusFour{
		ProgramSize:      648,
		Free:             352,
		Heap:             23,
		ProgramFlashSize: 1024,
		FlashSize:        1024,
		FlashChipID:      "144051",
		FlashFrequency:   40,
		FlashMode:        "DOUT",
		Features:         []string{"0809", "8F9AC787", "04368001", "000000CF", "010013C0", "C000F981", "00004004", "00001000", "54000020", "00000080", "00000000"},
		Drivers:          "1,2,!3,!4,!5,!6,7,!8,9,10,12,!16,!18,!19,!20,!21,!22,!24,26,!27,29,!30,!35,!37,!45,62,!68",
		Sensors:          "1,2,3,4,5,6",
		I2CDriver:        "7",
	}

	assert.Equal(t, expectedStatusFour, *statusFour)
}

func Test_UnmarshalStatusFive(t *testing.T) {
	statusFive, err := UnmarshalStatusFive([]byte(JsonData))

	if !assert.NoError(t, err) {
		return
	}

	expectedStatusFive := StatusFive{
		Hostname:   "main-0614",
		IPAddress:  "192.168.1.158",
		Gateway:    "192.168.1.1",
		Subnetmask: "255.255.255.0",
		DNSServer1: "192.168.1.1",
		DNSServer2: "0.0.0.0",
		Mac:        "2C:F4:32:FA:22:66",
		Webserver:  2,
		HTTPAPI:    1,
		WifiConfig: 4,
		WifiPower:  17,
	}

	assert.Equal(t, expectedStatusFive, *statusFive)
}

func Test_UnmarshalStatusSix(t *testing.T) {
	statusSix, err := UnmarshalStatusSix([]byte(JsonData))

	if !assert.NoError(t, err) {
		return
	}

	expectedStatusSix := StatusSix{
		MqttHost:       "192.168.1.129",
		MqttPort:       1883,
		MqttClientMask: "DVES_%06X",
		MqttClient:     "DVES_FA2266",
		MqttUser:       "DVES_USER",
		MqttCount:      1,
		MAXPACKETSIZE:  1200,
		KEEPALIVE:      30,
		SOCKETTIMEOUT:  4,
	}

	assert.Equal(t, expectedStatusSix, *statusSix)
}

func Test_UnmarshalStatusSeven(t *testing.T) {
	statusSeven, err := UnmarshalStatusSeven([]byte(JsonData))

	if !assert.NoError(t, err) {
		return
	}

	expectedStatusSeven := StatusSeven{
		UTC:      time.Date(2024, 8, 31, 13, 17, 41, 0, time.UTC),
		Local:    TasmotaTime(time.Date(2024, 8, 31, 14, 17, 41, 0, time.UTC)),
		StartDST: TasmotaTime(time.Date(2024, 3, 31, 2, 0, 0, 0, time.UTC)),
		EndDST:   TasmotaTime(time.Date(2024, 10, 27, 3, 0, 0, 0, time.UTC)),
		Timezone: "+01:00",
		Sunrise:  "06:06",
		Sunset:   "19:33",
	}

	assert.Equal(t, expectedStatusSeven, *statusSeven)
}

func Test_UnmarshalStatusEight(t *testing.T) {
	statusEight, err := UnmarshalStatusEight([]byte(JsonData))

	if !assert.NoError(t, err) {
		return
	}

	expectedStatusEight := StatusEight{
		Time: TasmotaTime(time.Date(2024, 8, 31, 14, 17, 41, 0, time.UTC)),
	}

	assert.Equal(t, expectedStatusEight, *statusEight)
}

func Test_UnmarshalStatusEleven(t *testing.T) {
	statusEleven, err := UnmarshalStatusEleven([]byte(JsonData))

	if !assert.NoError(t, err) {
		return
	}

	expectedStatusEleven := StatusEleven{
		Time:      TasmotaTime(time.Date(2024, 8, 31, 14, 17, 41, 0, time.UTC)),
		Uptime:    "0T00:27:03",
		UptimeSec: 1623,
		Heap:      22,
		SleepMode: "Dynamic",
		Sleep:     50,
		LoadAvg:   19,
		MqttCount: 1,
		POWER:     "OFF",
		Wifi: struct {
			AP        int    `json:"AP"`
			SSID      string `json:"SSId"`
			BSSID     string `json:"BSSId"`
			Channel   int    `json:"Channel"`
			Mode      string `json:"Mode"`
			RSSI      int    `json:"RSSI"`
			Signal    int    `json:"Signal"`
			LinkCount int    `json:"LinkCount"`
			Downtime  string `json:"Downtime"`
		}{
			AP:        1,
			SSID:      "ALHN-ED72",
			BSSID:     "EC:84:B4:0C:86:09",
			Channel:   3,
			Mode:      "11n",
			RSSI:      68,
			Signal:    -66,
			LinkCount: 1,
			Downtime:  "0T00:00:04",
		},
	}

	assert.Equal(t, expectedStatusEleven, *statusEleven)
}

func Test_UnmarshalStatus(t *testing.T) {
	status, err := UnmarshalStatus([]byte(JsonData))

	if !assert.NoError(t, err) {
		return
	}

	expectedStatus := Status{
		Status: StatusZero{
			Module:       1,
			DeviceName:   "Tasmota",
			FriendlyName: []string{"Tasmota"},
			Topic:        "main",
			ButtonTopic:  "0",
			Power:        "0",
			PowerLock:    "0",
			PowerOnState: 3,
			LedState:     1,
			LedMask:      "FFFF",
			SaveData:     1,
			SaveState:    1,
			SwitchTopic:  "0",
			SwitchMode:   []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			ButtonRetain: 0,
			SwitchRetain: 0,
			SensorRetain: 0,
			PowerRetain:  0,
			InfoRetain:   0,
			StateRetain:  0,
			StatusRetain: 0,
		},
		StatusPRM: StatusOne{
			Baudrate:      115200,
			SerialConfig:  "8N1",
			GroupTopic:    "tasmotas",
			OtaURL:        "http://ota.tasmota.com/tasmota/release/tasmota.bin.gz",
			RestartReason: "Software/System restart",
			Uptime:        "0T00:27:03",
			StartupUTC:    TasmotaTime(time.Date(2024, 8, 31, 12, 50, 38, 0, time.UTC)),
			Sleep:         50,
			CfgHolder:     4617,
			BootCount:     14,
			BCResetTime:   TasmotaTime(time.Date(2024, 8, 29, 15, 59, 39, 0, time.UTC)),
			SaveCount:     97,
			SaveAddress:   "FB000",
		},
		StatusFWR: StatusTwo{
			Version:       "14.2.0(release-tasmota)",
			BuildDateTime: TasmotaTime(time.Date(2024, 8, 14, 12, 36, 35, 0, time.UTC)),
			Boot:          31,
			Core:          "2_7_7",
			SDK:           "2.2.2-dev(38a443e)",
			CPUFrequency:  80,
			Hardware:      "ESP8285N08",
			CR:            "373/699",
		},
		StatusLOG: StatusThree{
			SerialLog:  2,
			WebLog:     2,
			MqttLog:    0,
			SysLog:     0,
			LogHost:    "",
			LogPort:    514,
			SSID:       []string{"ALHN-ED72", ""},
			TelePeriod: 300,
			Resolution: "558180C0",
			SetOption:  []string{"00008009", "2805C80001000600003C5A0A192800000000", "00000080", "00006000", "00004000", "00000000"},
		},
		StatusMEM: StatusFour{
			ProgramSize:      648,
			Free:             352,
			Heap:             23,
			ProgramFlashSize: 1024,
			FlashSize:        1024,
			FlashChipID:      "144051",
			FlashFrequency:   40,
			FlashMode:        "DOUT",
			Features:         []string{"0809", "8F9AC787", "04368001", "000000CF", "010013C0", "C000F981", "00004004", "00001000", "54000020", "00000080", "00000000"},
			Drivers:          "1,2,!3,!4,!5,!6,7,!8,9,10,12,!16,!18,!19,!20,!21,!22,!24,26,!27,29,!30,!35,!37,!45,62,!68",
			Sensors:          "1,2,3,4,5,6",
			I2CDriver:        "7",
		},
		StatusNET: StatusFive{
			Hostname:   "main-0614",
			IPAddress:  "192.168.1.158",
			Gateway:    "192.168.1.1",
			Subnetmask: "255.255.255.0",
			DNSServer1: "192.168.1.1",
			DNSServer2: "0.0.0.0",
			Mac:        "2C:F4:32:FA:22:66",
			Webserver:  2,
			HTTPAPI:    1,
			WifiConfig: 4,
			WifiPower:  17,
		},
		StatusMQT: StatusSix{
			MqttHost:       "192.168.1.129",
			MqttPort:       1883,
			MqttClientMask: "DVES_%06X",
			MqttClient:     "DVES_FA2266",
			MqttUser:       "DVES_USER",
			MqttCount:      1,
			MAXPACKETSIZE:  1200,
			KEEPALIVE:      30,
			SOCKETTIMEOUT:  4,
		},
		StatusTIM: StatusSeven{
			UTC:      time.Date(2024, 8, 31, 13, 17, 41, 0, time.UTC),
			Local:    TasmotaTime(time.Date(2024, 8, 31, 14, 17, 41, 0, time.UTC)),
			StartDST: TasmotaTime(time.Date(2024, 3, 31, 2, 0, 0, 0, time.UTC)),
			EndDST:   TasmotaTime(time.Date(2024, 10, 27, 3, 0, 0, 0, time.UTC)),
			Timezone: "+01:00",
			Sunrise:  "06:06",
			Sunset:   "19:33",
		},
		StatusSNS: StatusEight{
			Time: TasmotaTime(time.Date(2024, 8, 31, 14, 17, 41, 0, time.UTC)),
		},
		StatusSTS: StatusEleven{
			Time:      TasmotaTime(time.Date(2024, 8, 31, 14, 17, 41, 0, time.UTC)),
			Uptime:    "0T00:27:03",
			UptimeSec: 1623,
			Heap:      22,
			SleepMode: "Dynamic",
			Sleep:     50,
			LoadAvg:   19,
			MqttCount: 1,
			POWER:     "OFF",
			Wifi: struct {
				AP        int    `json:"AP"`
				SSID      string `json:"SSId"`
				BSSID     string `json:"BSSId"`
				Channel   int    `json:"Channel"`
				Mode      string `json:"Mode"`
				RSSI      int    `json:"RSSI"`
				Signal    int    `json:"Signal"`
				LinkCount int    `json:"LinkCount"`
				Downtime  string `json:"Downtime"`
			}{
				AP:        1,
				SSID:      "ALHN-ED72",
				BSSID:     "EC:84:B4:0C:86:09",
				Channel:   3,
				Mode:      "11n",
				RSSI:      68,
				Signal:    -66,
				LinkCount: 1,
				Downtime:  "0T00:00:04",
			},
		},
	}

	assert.Equal(t, expectedStatus, *status)
}
