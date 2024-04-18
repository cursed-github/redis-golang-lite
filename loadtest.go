package main

import (
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"
)

func loadtest() {
    var wg sync.WaitGroup
    clientCount := 50 // Limit of concurrent goroutines
    totalKeys := 10000 // Total keys to set and get
    sem := make(chan struct{}, clientCount) // Semaphore-like channel for limiting concurrency

    // Setting keys
    for i := 0; i < totalKeys; i++ {
        wg.Add(1)
        go func(key int) {
            defer wg.Done()
            sem <- struct{}{}        // Acquire semaphore
            setKey(key)
            <-sem                    // Release semaphore
        }(i + 1) // +1 to start keys from 1 to 10000
    }

    // Getting keys
    for i := 0; i < totalKeys; i++ {
        wg.Add(1)
        go func(key int) {
            defer wg.Done()
            sem <- struct{}{}        // Acquire semaphore
            getKey(key)
            <-sem                    // Release semaphore
        }(i + 1) // +1 to start keys from 1 to 10000
    }

    wg.Wait() // Wait for all goroutines to complete
}

func setKey(key int) {
    conn, err := net.Dial("tcp", "localhost:6379")
    if err != nil {
        fmt.Printf("Error connecting: %v\n", err)
        return
    }
    defer conn.Close()

    // Generate a random value for the key
    value := rand.Intn(10000) // Random values from 0 to 9999
    ttl := 3600 // TTL in seconds (1 hour)

    command := fmt.Sprintf("*5\r\n$3\r\nSET\r\n$%d\r\nkey%d\r\n$%d\r\n%d\r\n$2\r\nEX\r\n$4\r\n%d\r\n", 
                           len(fmt.Sprintf("key%d", key)), key, 
                           len(fmt.Sprintf("%d", value)), value, ttl)
    _, err = conn.Write([]byte(command))
    if err != nil {
        fmt.Printf("Error sending SET command: %v\n", err)
        return
    }

    // Wait to receive a response
    buffer := make([]byte, 1024)
    conn.SetReadDeadline(time.Now().Add(5 * time.Second)) // Set a timeout for reading
    n, err := conn.Read(buffer)
    if err != nil {
        fmt.Printf("Error reading response from SET: %v\n", err)
        return
    }

    fmt.Printf("Set key%d with TTL: %s\n", key, string(buffer[:n]))
}

func getKey(key int) {
    conn, err := net.Dial("tcp", "localhost:6379")
    if err != nil {
        fmt.Printf("Error connecting: %v\n", err)
        return
    }
    defer conn.Close()

    command := fmt.Sprintf("*2\r\n$3\r\nGET\r\n$%d\r\nkey%d\r\n", len(fmt.Sprintf("key%d", key)), key)
    _, err = conn.Write([]byte(command))
    if err != nil {
        fmt.Printf("Error sending GET command: %v\n", err)
        return
    }

    // Wait to receive a response
    buffer := make([]byte, 1024)
    conn.SetReadDeadline(time.Now().Add(5 * time.Second)) // Set a timeout for reading
    n, err := conn.Read(buffer)
    if err != nil {
        fmt.Printf("Error reading response from GET: %v\n", err)
        return
    }

    fmt.Printf("Get key%d: %s\n", key, string(buffer[:n]))
}
