package slack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync/atomic"

	"golang.org/x/net/websocket"
)

type Client struct {
	Name     string
	Channels map[string]string

	c    *websocket.Conn
	next uint64
}

func Connect(token string) (Client, error) {
	res, err := http.Get(fmt.Sprintf("https://slack.com/api/rtm.start?token=%s", token))
	if err != nil {
		return Client{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return Client{}, fmt.Errorf("API request failed with code %d", res.StatusCode)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Client{}, err
	}

	var r struct {
		Ok    bool   `json:"ok"`
		Error string `json:"error"`
		Url   string `json:"url"`
		Self  struct {
			ID string `json:"id"`
		} `json:"self"`

		Channels []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"channels"`
	}
	if err := json.Unmarshal(b, &r); err != nil {
		return Client{}, err
	}

	if !r.Ok {
		return Client{}, fmt.Errorf("Slack error: %s", r.Error)
	}

	chans := make(map[string]string)
	for _, c := range r.Channels {
		chans[c.Name] = c.ID
	}

	ws, err := websocket.Dial(r.Url, "", "https://api.slack.com/")
	if err != nil {
		return Client{}, err
	}

	return Client{
		c:        ws,
		Name:     r.Self.ID,
		Channels: chans,
	}, nil
}

func (c Client) Send(m Message) error {
	m.ID = atomic.AddUint64(&c.next, 1)
	if ch, ok := c.Channels[m.Channel]; ok {
		m.Channel = ch
	}
	b, _ := json.Marshal(m)
	fmt.Printf("MESSAGE: %s\n", string(b))
	b, _ = json.Marshal(c.Channels)
	fmt.Printf("CHANNELS: %s\n", string(b))

	return websocket.JSON.Send(c.c, m)
}

func (c Client) Receive() (Message, error) {
	var m Message
	return m, websocket.JSON.Receive(c.c, &m)
}
