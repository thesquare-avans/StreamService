package main

import (
	"log"
	"net/http"

	"github.com/thesquare-avans/StreamService/distribution"
	"github.com/thesquare-avans/StreamService/fsd"
	"github.com/thesquare-avans/StreamService/hls"
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

	center := distribution.NewCenter()
	checkErr(center.NewStream("0"))

	go func() {
		handler := hls.NewHandler(center, "", 5)
		log.Println("Listening HTTP")
		checkErr(http.ListenAndServe(":8080", handler))
	}()

	for {
		fragment := <-ch
		duration, err := fragment.GetVideoLength()
		checkErr(err)
		log.Printf("Received fragment, duration: %.3fms, length: %d\n", duration, fragment.Length)
		center.PushToStream("0", fragment)
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
