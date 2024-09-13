package mqtt_sonoff_basic_r2_test

import (
	sonoff "github.com/fromsi/mqtt_sonoff_basic_r2"
	mqtt "github.com/mochi-mqtt/server/v2"
	mqttauth "github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Example_withServer() {
	// init

	sigSystem := make(chan os.Signal, 1)

	signal.Notify(sigSystem, syscall.SIGINT, syscall.SIGTERM)

	server := mqtt.New(
		&mqtt.Options{
			InlineClient: true,
		},
	)

	tcp := listeners.NewTCP(listeners.Config{ID: "t1", Address: ":1883"})
	err := server.AddListener(tcp)

	if err != nil {
		panic(err.Error())
	}

	err = server.AddHook(new(mqttauth.AllowHook), nil)

	if err != nil {
		panic(err.Error())
	}

	sonoffServer, err := sonoff.NewSonoffBasicR2WithServer(server, 2)

	if err != nil {
		panic(err.Error())
	}

	// run

	go func() {
		_ = sonoffServer.Serve()
		_ = server.Serve()
	}()

	// ... your code ...

	go func() {
		for {
			select {
			case id, ok := <-sonoffServer.TeleConnected():
				if !ok {
					return
				}

				log.Println("Connected", id)

				sonoffServer.PowerToggle(id)
			}
		}
	}()

	// stop

	select {
	case <-sigSystem:
		_ = sonoffServer.Close()
		_ = server.Close()
	}
}

func Example_withoutServer() {
	// init

	sigSystem := make(chan os.Signal, 1)

	signal.Notify(sigSystem, syscall.SIGINT, syscall.SIGTERM)

	server, err := sonoff.NewSonoffBasicR2("", 1883, 0)

	if err != nil {
		panic(err.Error())
	}

	// run

	go func() {
		_ = server.Serve()
	}()

	// ... your code ...

	go func() {
		for {
			select {
			case id, ok := <-server.TeleConnected():
				if !ok {
					return
				}

				log.Println("Connected", id)

				server.PowerToggle(id)
			}
		}
	}()

	// stop

	select {
	case <-sigSystem:
		_ = server.Close()
	}
}
