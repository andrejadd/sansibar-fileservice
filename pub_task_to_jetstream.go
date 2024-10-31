package main

import (
    "log"
    "github.com/nats-io/nats.go"
)

func main() {
    nc, err := nats.Connect(nats.DefaultURL)
    if err != nil {
        log.Fatal(err)
    }
    defer nc.Close()

    js, err := nc.JetStream()
    if err != nil {
        log.Fatal(err)
    }

    // Publish a message to the "task_queue" subject
    ack, err := js.QueuePublish("task_queue", "work", []byte("This is a task to be processed"))
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Message published with sequence: %d", ack.Sequence)
}
