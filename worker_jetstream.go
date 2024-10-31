package main

import (
    "log"
    "time"
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

    // Subscribe to "task_queue" with a queue group "workers"
    _, err = js.QueueSubscribe("task_queue", "work", func(msg *nats.Msg) {
        log.Printf("Received task: %s", string(msg.Data))

        // Simulate task processing (5 seconds)
        time.Sleep(5 * time.Second)

        // Acknowledge the message after processing
        if err := msg.Ack(); err != nil {
            log.Printf("Failed to acknowledge message: %v", err)
        } else {
            log.Printf("Task acknowledged")
        }
    }, nats.ManualAck(), nats.Durable("worker_group"))

    if err != nil {
        log.Fatal(err)
    }

    log.Println("Worker subscribed")

    // Keep the worker running
    select {}
}


