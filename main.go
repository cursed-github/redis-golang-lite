package main

import (
	"fmt"
	"kvstore/resp"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func handleRequest(conn net.Conn){
	defer conn.Close()
	respParser := resp.Resp{}
	payload, err:= respParser.DeserializeResp(conn)
	if err != nil {
		fmt.Println("Error parsing:", err)
		return 
	}

	// Process the payload
	fmt.Printf("Received payload: %+v\n", payload)

	_, err = conn.Write([]byte("Received\n"))
	if err != nil {
		fmt.Println("Error writing:", err)
		return 
	}
}

func main() {
	listner, err:= net.Listen("tcp", ":6379")
	if err!=nil {
		fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }
	defer listner.Close()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-stopChan
		fmt.Println("Shutting down server...")
		listner.Close()
		// Add any cleanup tasks here
		os.Exit(0)
	}()
	for {
		con, err := listner.Accept()
		if err!=nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}

		go handleRequest(con)
	}
}