package slack

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type User struct {
	ID       string
	Name     string
	Avatar   string
	TZOffset int
}

func (c *Client) fetchUsers() {
	res, err := http.Get(c.url("/users.list?token=%s", c.token))
	if err != nil {
		return
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	var data struct {
		OK      bool `json:"ok"`
		Members []struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			TZOffset int    `json:"tz_offset"`

			Profile struct {
				Image512 string `json:"image_512"`
			} `json:"profile"`
		} `json:"members"`
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return
	}

	if !data.OK {
		return
	}

	for _, user := range data.Members {
		u := User{
			ID:       user.ID,
			Name:     user.Name,
			Avatar:   user.Profile.Image512,
			TZOffset: user.TZOffset,
		}
		c.users[user.ID] = u
		c.users["@"+user.Name] = u
	}
}

func (c *Client) FindUser(by string) (User, bool) {
	if user, ok := c.users[by]; ok {
		return user, true
	}

	c.fetchUsers()
	user, ok := c.users[by]
	return user, ok
}
