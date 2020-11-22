package slack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"golang.org/x/net/websocket"
)

type Client struct {
	Name string

	users map[string]User
	chans map[string]Channel
	idmap map[string]string
	token string
	c     *websocket.Conn
	next  uint64
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
		chans["#"+c.Name] = c.ID
		chans[c.ID] = "#" + c.Name
	}

	ws, err := websocket.Dial(r.Url, "", "https://api.slack.com/")
	if err != nil {
		return Client{}, err
	}

	c := Client{
		c:     ws,
		token: token,
		idmap: chans,
		users: make(map[string]User),
		chans: make(map[string]Channel),

		Name: r.Self.ID,
	}
	c.fetchUsers()

	return c, nil
}

func (c Client) url(rel string, args ...interface{}) string {
	if !strings.HasPrefix(rel, "/") {
		rel = "/" + rel
	}
	return base + fmt.Sprintf(rel, args...)
}

func (c Client) Send(m Message) error {
	m.ID = atomic.AddUint64(&c.next, 1)
	if ch, ok := c.idmap[m.Channel]; ok {
		m.Channel = ch
	}

	c.intern(&m)
	return websocket.JSON.Send(c.c, m)
}

func (c Client) Receive() (Message, error) {
	var m Message
	err := websocket.JSON.Receive(c.c, &m)
	if err != nil {
		return Message{}, err
	}

	fmt.Printf("%s > ts: %s\n", m.Type, m.TS)
	if f, err := strconv.ParseFloat(m.TS, 64); err == nil {
		m.Received = time.Unix(int64(f), 0)
		fmt.Printf("setting received to %v\n", m.Received)
	}

	m.interned = true
	c.extern(&m)
	return m, err
}

func (c *Client) id2name(id string) string {
	if name, exists := c.idmap[id]; exists {
		return name
	}

	name := id
	if strings.HasPrefix(id, "U") {
		if user, found := c.FindUser(id); found {
			name = "@" + user.Name
		}
	}
	if strings.HasPrefix(id, "C") {
		if channel, found := c.FindChannel(id); found {
			name = "#" + channel.Name
		}
	}

	c.idmap[id] = name
	c.idmap[name] = id
	return name
}

func (c *Client) name2id(name string) string {
	if id, exists := c.idmap[name]; exists {
		return id
	}

	id := name
	if strings.HasPrefix(name, "@") {
		if user, found := c.FindUser(name); found {
			id = user.ID
		}
	}
	if strings.HasPrefix(id, "#") {
		if channel, found := c.FindChannel(name); found {
			id = channel.ID
		}
	}

	c.idmap[id] = name
	c.idmap[name] = id
	return id
}
