package main

import (
	"fmt"
	"time"

	"github.com/thesquare-avans/StreamService/stream"
)

func main() {
	server, err := stream.NewServer(":1312", "/home/sem/test")
	checkErr(err)
	fmt.Println("Listening and waiting...")

	err = server.WaitForClient()
	checkErr(err)
	defer server.Close()
	fmt.Println("Client connected")

	for {
		start := time.Now()
		len, err := server.ReceiveSingle()
		checkErr(err)
		fmt.Printf("Received one fragment, %d bytes, latency: %s\n", len, time.Now().Sub(start))
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
