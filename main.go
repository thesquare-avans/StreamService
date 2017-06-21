package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/thesquare-avans/StreamService/distribution"
	"github.com/thesquare-avans/StreamService/fsd"
	"github.com/thesquare-avans/StreamService/stream"
)

const (
	ConcurrentStreams = 4
	StreamStartPort   = 1234
	StreamServerPort  = ":8080"
)

func main() {
	center := distribution.NewCenter()

	for i := 0; i < ConcurrentStreams; i++ {
		ch := make(chan *fsd.Fragment, 16)
		port := StreamStartPort + i

		log.Printf("Starting StreamServer %d on port %d", i, port)
		server, err := stream.NewServer("", port, ch)
		checkErr(err)

		go func() {
			for {
				err := server.Run()
				logErr(err)
			}
		}()

		streamId := strconv.Itoa(i)
		logErr(center.NewStream(streamId))

		go func() {
			for {
				fragment := <-ch
				duration, err := fragment.GetDuration()
				logErr(err)
				log.Printf("Received fragment, duration: %s, length: %d\n", duration, fragment.Length)
				logErr(center.PushToStream(streamId, fragment))
			}
		}()
	}

	streamServer := stream.NewStreamServer(center)
	checkErr(http.ListenAndServe(StreamServerPort, streamServer))
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func logErr(err error) {
	if err != nil {
		log.Println("error:", err)
	}
}
