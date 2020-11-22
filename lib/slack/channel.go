package slack

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Channel struct {
	ID    string
	Name  string
	Topic string
}

func (c *Client) fetchChannels() {
	res, err := http.Get(c.url("/conversations.list?token=%s", c.token))
	if err != nil {
		return
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	var data struct {
		OK       bool `json:"ok"`
		Channels []struct {
			ID   string `json:"id"`
			Name string `json:"name"`

			Topic struct {
				Value string `json:"value"`
			} `json:"topic"`
		} `json:"members"`
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return
	}

	if !data.OK {
		return
	}

	for _, channel := range data.Channels {
		ch := Channel{
			ID:    channel.ID,
			Name:  channel.Name,
			Topic: channel.Topic.Value,
		}
		c.chans[channel.ID] = ch
		c.chans["#"+channel.Name] = ch
	}
}

func (c *Client) FindChannel(by string) (Channel, bool) {
	if channel, ok := c.chans[by]; ok {
		return channel, true
	}

	c.fetchChannels()
	channel, ok := c.chans[by]
	return channel, ok
}
