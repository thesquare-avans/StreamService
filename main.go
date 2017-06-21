package main

import (
	"log"
	"net/http"

	"github.com/thesquare-avans/StreamService/fsd"
	"github.com/thesquare-avans/StreamService/stream"
)

func main() {
	ch := make(chan *fsd.Fragment, 16)

	server, err := stream.NewServer("", 1234, ch)
	log.Println("Listening")
	checkErr(err)

	go func() {
		log.Println("Waiting for client")
		err := server.Run()
		checkErr(err)
		err = server.Close()
		log.Println("Stop listening, return:", err)
		checkErr(err)
	}()

	streamServer := stream.NewStreamServer()

	go func() {
		log.Println("Listening HTTP")
		checkErr(http.ListenAndServe(":8080", streamServer))
	}()

	for {
		fragment := <-ch
		duration, err := fragment.GetDuration()
		checkErr(err)
		log.Printf("Received fragment, duration: %s, length: %d\n", duration, fragment.Length)
		streamServer.PushFragment(fragment)
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
