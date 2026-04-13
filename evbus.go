package main

import (
	"fmt"
	"time"

	evbus "github.com/asaskevich/EventBus"
)

func test() {
	bus := evbus.New()
	bus.SubscribeAsync("topic:super", superHandler, true)
	bus.SubscribeAsync("topic:super", superHandler2, false)

	bus.Publish("topic:super", 123)

	bus.WaitAsync()
}

func superHandler(num int) {
	time.Sleep(time.Second)
	fmt.Println("handled num", num)
}

func superHandler2(num int) {
	time.Sleep(time.Second)
	fmt.Println("handled num from handler 2", num)
}
