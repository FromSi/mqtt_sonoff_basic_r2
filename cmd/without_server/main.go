package main

import (
	sonoff "github.com/fromsi/mqtt_sonoff_basic_r2"
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

				r, _ := server.StatusPhysicalButton(id)
				log.Println("StatusPhysicalButton", r)

				response, _ := server.Status(id)
				log.Println(response.Status.Topic, "POWER", response.Status.Power, "TIME", response.StatusTIM.Local.ToTime().String())

				time.Sleep(1 * time.Second)

				server.PowerOn(id)

				response, _ = server.Status(id)
				log.Println(response.Status.Topic, "POWER", response.Status.Power, "TIME", response.StatusTIM.Local.ToTime().String())

				time.Sleep(1 * time.Second)

				server.PowerOff(id)

				response, _ = server.Status(id)
				log.Println(response.Status.Topic, "POWER", response.Status.Power, "TIME", response.StatusTIM.Local.ToTime().String())

				time.Sleep(1 * time.Second)

				server.PowerToggle(id)

				response, _ = server.Status(id)
				log.Println(response.Status.Topic, "POWER", response.Status.Power, "TIME", response.StatusTIM.Local.ToTime().String())

				time.Sleep(1 * time.Second)

				server.PowerOff(id)

				response, _ = server.Status(id)
				log.Println(response.Status.Topic, "POWER", response.Status.Power, "TIME", response.StatusTIM.Local.ToTime().String())
			}
		}
	}()

	go func() {
		for {
			select {
			case id, ok := <-server.TeleDisconnected():
				if !ok {
					return
				}

				log.Println("Disconnected", id)
			}
		}
	}()

	// stop

	select {
	case <-sigSystem:
		_ = server.Close()
	}
}
