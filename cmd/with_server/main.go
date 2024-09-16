package main

import (
	sonoff "github.com/fromsi/mqtt_sonoff_basic_r2"
	mqtt "github.com/mochi-mqtt/server/v2"
	mqttauth "github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
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

				response, _ := sonoffServer.Status(id)
				log.Println(response.Status.Topic, "POWER", response.Status.Power, "TIME", response.StatusTIM.Local.ToTime().String())

				time.Sleep(1 * time.Second)

				sonoffServer.PowerOn(id)

				response, _ = sonoffServer.Status(id)
				log.Println(response.Status.Topic, "POWER", response.Status.Power, "TIME", response.StatusTIM.Local.ToTime().String())

				time.Sleep(1 * time.Second)

				sonoffServer.PowerOff(id)

				response, _ = sonoffServer.Status(id)
				log.Println(response.Status.Topic, "POWER", response.Status.Power, "TIME", response.StatusTIM.Local.ToTime().String())

				time.Sleep(1 * time.Second)

				sonoffServer.PowerToggle(id)

				response, _ = sonoffServer.Status(id)
				log.Println(response.Status.Topic, "POWER", response.Status.Power, "TIME", response.StatusTIM.Local.ToTime().String())

				time.Sleep(1 * time.Second)

				sonoffServer.PowerOff(id)

				response, _ = sonoffServer.Status(id)
				log.Println(response.Status.Topic, "POWER", response.Status.Power, "TIME", response.StatusTIM.Local.ToTime().String())
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
