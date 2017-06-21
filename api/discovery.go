package api

import (
	"encoding/json"
	"log"

	"github.com/thesquare-avans/StreamService/context"
	"github.com/thesquare-avans/StreamService/distribution"
	"github.com/thesquare-avans/StreamService/transport"
	"github.com/zhouhui8915/go-socket.io-client"
)

type DiscoveryConnection struct {
	client *socketio_client.Client
	center *distribution.Center
}

func NewDiscoveryConnection(url string) (*DiscoveryConnection, error) {
	options := &socketio_client.Options{}
	client, err := socketio_client.NewClient(url, options)
	if err != nil {
		return nil, err
	}
	return &DiscoveryConnection{
		client: client,
	}, nil
}

func (d *DiscoveryConnection) Run() {
	d.client.On("start", handleStart)
}

func handleStart(msg string) {
	var payload transport.Payload
	err := json.Unmarshal([]byte(msg), &payload)
	if err != nil {
		log.Println("error:", err)
		return
	}
	payload.Verify(&context.GlobalContext().PrivateKey.PublicKey)

	var args StartArgs
	err = json.Unmarshal([]byte(payload.Payload), &args)
	if err != nil {
		log.Println("error:", err)
		return
	}
}
