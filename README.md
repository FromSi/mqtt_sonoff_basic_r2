# MQTT Sonoff Basic R2
Communication with [Sonoff Basic R2](https://sonoff.tech/product/diy-smart-switches/basicr2) via `MQTT`. The `Sonoff Basic R2` must have [Tasmota](https://tasmota.github.io/docs) firmware. For `MQTT` the [Mochi-MQTT](https://github.com/mochi-mqtt/server) library is used.

## Install

```shell
go get github.com/formsi/mqtt_sonoff_basic_r2
```

**Note:** Current lib uses [Go Modules](https://go.dev/wiki/Modules) to manage dependencies.

**Note 2:** Шаблон для full topic: `%prefix%/%topic%/`

## Features
* Functions to start or stop the server
* Receive notification of connection or disconnection
* Changing Power ON/OFF/TOGGLE 
* Changing Physical Button ON/OFF 
* Getting Status (0-11 and physical_button) with timeout and structs

## Examples

### Functions to start or stop the server
More on the file `cmd/without_server/main.go`

```go
package main

import (
    sonoff "github.com/fromsi/mqtt_sonoff_basic_r2"
    "os"
    "os/signal"
    "syscall"
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
    // ...

    // stop

    select {
    case <-sigSystem:
        _ = server.Close()
    }
}

```

### Receive notification of connection or disconnection
```go
//...

func main() {
    // init
    // ...

    // run
    // ...

    // ... your code ...

    go func() {
        for {
            select {
            case id, ok := <-server.TeleConnected():
                if !ok {
                    return
                }

                log.Println("Connected", id)
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
    // ...
}
```

### Changing Power ON/OFF/TOGGLE
```go
//...

func main() {
    // init
    // ...

    // run
    // ...

    // ... your code ...

    // connected with id
    time.Sleep(1 * time.Second)
    
    server.PowerOn(id)
    
    time.Sleep(1 * time.Second)
    
    server.PowerOff(id)
    
    time.Sleep(1 * time.Second)
    
    server.PowerToggle(id)
    
    time.Sleep(1 * time.Second)
    
    server.PowerOff(id)
    // connected with id

    // stop
    // ...
}
```

### Changing Physical Button ON/OFF
```go
//...

func main() {
    // init
    // ...

    // run
    // ...

    // ... your code ...

    // connected with id
    value, _ := StatusPhysicalButton(id)
    log.Println(id, "PhysicalButton", value)

    time.Sleep(1 * time.Second)
    
    server.PhysicalButtonOn(id)
    value, _ = StatusPhysicalButton(id)
    log.Println(id, "PhysicalButton", value)
    
    time.Sleep(1 * time.Second)
    
    server.PhysicalButtonOff(id)
    value, _ = StatusPhysicalButton(id)
    log.Println(id, "PhysicalButton", value)
    // connected with id

    // stop
    // ...
}
```

### Getting Status (0-11) with timeout and structs
```go
//...

func main() {
    // init
    // ...

    // run
    // ...

    // ... your code ...

    // connected with id
    time.Sleep(1 * time.Second)
    
    server.PowerOn(id)
    
    response, _ = server.Status(id)
    log.Println(response.Status.Topic, "POWER", response.Status.Power, "TIME", response.StatusTIM.Local.ToTime().String())
    
    time.Sleep(1 * time.Second)
    
    server.PowerOff(id)
    
    response, _ = server.Status(id)
    log.Println(response.Status.Topic, "POWER", response.Status.Power, "TIME", response.StatusTIM.Local.ToTime().String())
    // connected with id

    // stop
    // ...
}
```

### Using the library as a wrapper for your server 
More on the (mochi-mqtt/server)[https://github.com/mochi-mqtt/server]

```go
package main

import (
    sonoff "github.com/fromsi/mqtt_sonoff_basic_r2"
    mqtt "github.com/mochi-mqtt/server/v2"
    mqttauth "github.com/mochi-mqtt/server/v2/hooks/auth"
    "github.com/mochi-mqtt/server/v2/listeners"
)

func main() {
    server := mqtt.New(
        &mqtt.Options{
            InlineClient: true,
        },
    )

    tcp := listeners.NewTCP(listeners.Config{ID: "t1", Address: ":1883"})
    _ = server.AddListener(tcp)
    _ = server.AddHook(new(mqttauth.AllowHook), nil)

    sonoffServer, _ := sonoff.NewSonoffBasicR2WithServer(server, 2)
	
    // run

    go func() {
        _ = sonoffServer.Serve()
        _ = server.Serve()
    }()

    // more on the file `cmd/without_server/main.go`
    // ... your code ...

    // stop

    sonoffServer.Close()
    server.Close()
}
```
