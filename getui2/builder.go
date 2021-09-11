package getui2

import (
	"encoding/json"
	"fmt"
	"github.com/lithammer/shortuuid/v3"
	"strings"
)

// Payload represents a notification which holds the content that will be
// marshalled as JSON.
type Payload struct {
	content map[string]interface{}
	_custom map[string]interface{}
}

type audience struct {
	Cid   []string `json:"cid,omitempty"`
	Alias []string `json:"alias,omitempty"`
}

type settings struct {
	Ttl      int            `json:"ttl,omitempty"`
	Strategy map[string]int `json:"strategy,omitempty"`
}

type message struct {
	Notification interface{} `json:"notification,omitempty"`
}

type notification struct {
	Title        string `json:"title,omitempty"`
	Body         string `json:"body,omitempty"`
	BigText      string `json:"big_text,omitempty"`
	BigImage     string `json:"big_image,omitempty"`
	Logo         string `json:"logo,omitempty"`
	LogoUrl      string `json:"logo_url,omitempty"`
	ChannelId    string `json:"channel_id,omitempty"`
	ChannelName  string `json:"channel_name,omitempty"`
	ChannelLevel int    `json:"channel_level,omitempty"`
	ClickType    string `json:"click_type,omitempty"`
	Intent       string `json:"intent,omitempty"`
	Url          string `json:"url,omitempty"`
	Payload      string `json:"payload,omitempty"`
	NotifyId     int    `json:"notify_id,omitempty"`
	RingName     string `json:"ring_name,omitempty"`
	BadgeAddNum  int    `json:"badge_add_num,omitempty"`
}

// NewPayload returns a new Payload struct
func NewPayload() *Payload {
	return &Payload{
		map[string]interface{}{
			"request_id": strings.ToLower(shortuuid.New()),
		},
		map[string]interface{}{},
	}
}

func (p *Payload) AlertTitle(title string) *Payload {
	p.message().notification().Title = title
	return p
}

func (p *Payload) AlertBody(body string) *Payload {
	p.message().notification().Body = body
	return p
}

func (p *Payload) AlertPayload(payload string) *Payload {
	p.message().notification().Payload = payload
	return p
}

func (p *Payload) ClickType(clickType string) *Payload {
	p.message().notification().ClickType = clickType
	return p
}

func (p *Payload) Intent(intent string) *Payload {
	p.message().notification().Intent = intent
	return p
}

func (p *Payload) BadgeAddNum(badgeAddNum int) *Payload {
	p.message().notification().BadgeAddNum = badgeAddNum
	return p
}

func (p *Payload) Ttl(ttl int) *Payload {
	p.settings().Ttl = ttl
	return p
}

func (p *Payload) Strategy(name string, _type int) *Payload {
	strategy := p.settings().Strategy
	if strategy == nil {
		strategy = make(map[string]int)
		strategy["default"] = 1
	}
	strategy[name] = _type
	p.settings().Strategy = strategy
	return p
}

func (p *Payload) Cid(cid string) *Payload {
	p.audience().Cid = []string{cid}
	return p
}

func (p *Payload) Custom(key string, val interface{}) *Payload {
	p.custom()[key] = val
	return p
}

func (p *Payload) SetIntent() *Payload {
	p.ClickType("intent")
	template := "intent:#Intent;launchFlags=0x04000000;component=com.fansy.bark/io.dcloud.PandoraEntry;S.UP-OL-SU=true;S.title=%v;S.content=%v;S.payload=%v;end"
	intent := fmt.Sprintf(template, p.message().notification().Title, p.message().notification().Body, p.message().notification().Payload)
	p.Intent(intent)
	return p
}

func (p *Payload) SetUps() *Payload {
	ups := map[string]interface{}{
		"notification": p.message().notification(),
	}
	p.channel()["android"] = map[string]interface{}{
		"ups": ups,
	}
	return p
}

// MarshalJSON returns the JSON encoded version of the Payload
func (p *Payload) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(p._custom)
	if err == nil {
		p.AlertPayload(string(b))
	}
	return json.Marshal(p.content)
}

func (p *Payload) audience() *audience {
	if _, ok := p.content["audience"].(*audience); !ok {
		p.content["audience"] = &audience{}
	}
	return p.content["audience"].(*audience)
}

func (p *Payload) settings() *settings {
	if _, ok := p.content["settings"].(*settings); !ok {
		p.content["settings"] = &settings{}
	}
	return p.content["settings"].(*settings)
}

func (p *Payload) message() *message {
	if _, ok := p.content["push_message"].(*message); !ok {
		p.content["push_message"] = &message{}
	}
	return p.content["push_message"].(*message)
}

func (p *Payload) channel() map[string]interface{} {
	if _, ok := p.content["push_channel"].(map[string]interface{}); !ok {
		p.content["push_channel"] = make(map[string]interface{})
	}
	return p.content["push_channel"].(map[string]interface{})
}

func (m *message) notification() *notification {
	if _, ok := m.Notification.(*notification); !ok {
		m.Notification = &notification{}
	}
	return m.Notification.(*notification)
}

func (p *Payload) custom() map[string]interface{} {
	return p._custom
}
