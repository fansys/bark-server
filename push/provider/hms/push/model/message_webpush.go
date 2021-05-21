package model

import (
	"bark-server/push/provider/hms/push/constant"
	"time"
)

type WebPushConfig struct {
	Data         string               `json:"data,omitempty"`
	Headers      *WebPushHeaders      `json:"headers,omitempty"`
	HmsOptions   *HmsWebPushOption    `json:"hms_options,omitempty"`
	Notification *WebPushNotification `json:"notification,omitempty"`
}

type WebPushHeaders struct {
	TTL     string `json:"ttl,omitempty"`
	Topic   string `json:"topics,omitempty"`
	Urgency string `json:"urgency,omitempty"`
}

type HmsWebPushOption struct {
	Link string `json:"link,omitempty"`
}

type WebPushNotification struct {
	Title              string           `json:"title,omitempty"`
	Body               string           `json:"body,omitempty"`
	Actions            []*WebPushAction `json:"actions,omitempty"`
	Badge              string           `json:"badge,omitempty"`
	Dir                string           `json:"dir,omitempty"`
	Icon               string           `json:"icon,omitempty"`
	Image              string           `json:"image,omitempty"`
	Lang               string           `json:"lang,omitempty"`
	Renotify           bool             `json:"renotify,omitempty"`
	RequireInteraction bool             `json:"require_interaction,omitempty"`
	Silent             bool             `json:"silent,omitempty"`
	Tag                string           `json:"tag,omitempty"`
	Timestamp          int64            `json:"timestamp,omitempty"`
	Vibrate            []int            `json:"vibrate,omitempty"`
}

type WebPushAction struct {
	Action string `json:"action,omitempty"`
	Icon   string `json:"icon,omitempty"`
	Title  string `json:"title,omitempty"`
}

func GetDefaultWebPushConfig() *WebPushConfig {
	return &WebPushConfig{
		Data:       "web push data",
		Headers:    getDefaultHeaders(),
		HmsOptions: getDefaultHmsOptions(),
	}
}
func getDefaultHeaders() *WebPushHeaders {
	return &WebPushHeaders{
		TTL:     "990",
		Topic:   "topic",
		Urgency: constant.UrgencyVeryLow,
	}
}

func getDefaultHmsOptions() *HmsWebPushOption {
	return &HmsWebPushOption{
		Link: "https://www.huawei.com",
	}
}

func GetDefaultWebNotification() *WebPushNotification {
	return &WebPushNotification{
		Dir:       constant.DirAuto,
		Silent:    true,
		Timestamp: time.Now().Unix(),
	}
}
