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

    // Subscribe to the task_queue as a queue group named "workers"
    _, err = nc.QueueSubscribe("task_queue", "workers", func(msg *nats.Msg) {
        log.Printf("Received message: %s", string(msg.Data))
        // Process the message here
    })
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Worker is waiting for messages...")
    
    // Keep the connection alive
    select {}
}

