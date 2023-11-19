package hub

import "github.com/fakovacic/editor/internal/app"

type Client struct {
	Username string    `json:"username"`
	Color    string    `json:"color"`
	Ready    bool      `json:"ready"`
	Position *Position `json:"position,omitempty"`
}

type Position struct {
	Index  int `json:"index"`
	Length int `json:"length"`
}

func ToClient(client *app.Client) Client {
	vc := Client{
		Username: client.Username,
		Color:    client.Color,
		Ready:    client.Ready,
	}

	if client.Position != nil {
		vc.Position = &Position{
			Index:  client.Position.Index,
			Length: client.Position.Length,
		}
	}

	return vc
}
