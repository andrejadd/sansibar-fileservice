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

    // Create a stream that will hold the tasks (messages)
    _, err = js.AddStream(&nats.StreamConfig{
        Name:     "TASKS",
        Subjects: []string{"task_queue"},
    })
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Stream created for tasks")
}

