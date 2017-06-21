package api

type Error struct {
	Code string `json:"code"`
}

func (e *Error) Error() string {
	return e.Code
}

type User struct {
	Id        string `json:"id"`
	PublicKey string `json:"publicKey"`
}

type DefaultPayload struct {
	Success bool   `json:"success"`
	Error   *Error `json:"error"`
}

type RegisterArgs struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}

type StartArgs struct {
	StreamId string `json:"streamId"`
	Streamer *User  `json:"streamer"`
}

type StreamInfo struct {
	Id   string `json:"id"`
	Host string `json:"host"`
	Port int    `json:"port"`
}

type StartPayload struct {
	Success bool        `json:"success"`
	Error   *Error      `json:"error"`
	Data    *StreamInfo `json:"data"`
}
