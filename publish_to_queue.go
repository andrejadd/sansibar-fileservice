package main

import (
    "github.com/nats-io/nats.go"
    "log"
)

func main() {
    // Connect to NATS
    nc, err := nats.Connect("nats://localhost:4222")
    if err != nil {
        log.Fatal(err)
    }
    defer nc.Close()

    // Publish a message to the queue
    err = nc.Publish("task_queue", []byte("This is the second one"))
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Message published to task_queue")
}

