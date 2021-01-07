package slack

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/net/websocket"
	"io/ioutil"
	"net/http"
	"strings"
	"sync/atomic"
)

// Client defines a client for Slack API.
type Client struct {
	Id      string
	ws      *websocket.Conn
	counter uint64
	token   string
}

type SlackMembersResult struct {
	Ok      bool        `json:"ok"`
	Members []SlackUser `json:"members"`
}

type SlackChannelsResult struct {
	Ok       bool           `json:"ok"`
	Channels []SlackChannel `json:"channels"`
}

type SlackUser struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	IsAdmin     bool   `json:"is_admin"`
	IsBot       bool   `json:"is_bot"`
	ProfileDesc string `json:"profileDesc"`
	Class       string `json:"class"`
}

type SlackChannel struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// These two structures represent the response of the Slack API rtm.start.
// Only some fields are included. The rest are ignored by json.Unmarshal.
type responseRtmStart struct {
	Ok    bool         `json:"ok"`
	Error string       `json:"error"`
	Url   string       `json:"url"`
	Self  responseSelf `json:"self"`
}

type responseSelf struct {
	Id string `json:"Id"`
}

// These are the messages read off and written into the websocket. Since this
// struct serves as both read and write, we include the "Id" field which is
// required only for writing.
type Message struct {
	Id      uint64 `json:"Id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
	User    string `json:"user"`
}

// slackStart does a rtm.start, and returns a websocket URL and user ID. The
// websocket URL can be used to initiate an RTM session.
func Start(token string) (wsurl, id string, err error) {
	url := fmt.Sprintf("https://slack.com/api/rtm.start?token=%s", token)
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = fmt.Errorf("API request failed with code %d", resp.StatusCode)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var respObj responseRtmStart
	if err = json.Unmarshal(body, &respObj); err != nil {
		return
	}

	if !respObj.Ok {
		err = fmt.Errorf("Slack error: %s", respObj.Error)
		return
	}

	wsurl = respObj.Url
	id = respObj.Self.Id
	return
}

func (c *Client) GetMessage() (m Message, err error) {
	err = websocket.JSON.Receive(c.ws, &m)
	return
}

func (c *Client) PostMessage(m Message) error {
	m.Id = atomic.AddUint64(&c.counter, 1)
	return websocket.JSON.Send(c.ws, m)
}

// Starts a websocket-based Real Time API session and return the websocket
// and the ID of the (bot-)user whom the token belongs to.
func (c *Client) Connect(token string) error {
	wsurl, id, err := Start(token)
	if err != nil {
		return err
	}

	ws, err := websocket.Dial(wsurl, "", SlackOriginUrl)
	if err != nil {
		return err
	}

	c.ws = ws
	c.Id = id
	c.token = token

	return nil
}

func (c *Client) GetTeamUsers() ([]SlackUser, error) {
	resp, err := http.DefaultClient.Get("https://slack.com/api/users.list?token=" + c.token)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("unexpected http return code received: %d", resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var res SlackMembersResult
	if err := json.Unmarshal(respBytes, &res); err != nil {
		return nil, err
	}
	if !res.Ok {
		return nil, errors.New("members response not ok")
	}
	return res.Members, nil
}

func (c *Client) GetTeamChannels() ([]SlackChannel, error) {
	resp, err := http.DefaultClient.Get("https://slack.com/api/channels.list?token=" + c.token)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("unexpected http return code received: %d", resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var res SlackChannelsResult
	if err := json.Unmarshal(respBytes, &res); err != nil {
		return nil, err
	}
	if !res.Ok {
		return nil, errors.New("channels response not ok")
	}
	return res.Channels, nil
}

func IsSlackUserName(input string) bool {
	return strings.HasPrefix(input, "<@")
}

func IsSlackChannelName(input string) bool {
	return strings.HasPrefix(input, "<#")
}

func GetSlackUserName(input string) string {
	return strings.TrimSuffix(strings.TrimPrefix(input, "<@"), ">")
}

func GetSlackChannelName(input string) string {
	return strings.TrimSuffix(strings.TrimPrefix(input, "<#"), ">")
}
